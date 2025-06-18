package com.example.mangapdf.ui.detail

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.fragment.app.Fragment
import androidx.lifecycle.ViewModelProvider
import androidx.navigation.fragment.navArgs
import androidx.recyclerview.widget.LinearLayoutManager
import com.bumptech.glide.Glide
import com.example.mangapdf.databinding.FragmentDetailBinding
import com.yandex.mobile.ads.rewarded.*
import com.yandex.mobile.ads.common.AdRequestError
import com.yandex.mobile.ads.common.AdError
import com.yandex.mobile.ads.common.AdRequestConfiguration
import com.yandex.mobile.ads.common.ImpressionData
import androidx.core.view.ViewCompat

class DetailFragment : Fragment() {

    private var _binding: FragmentDetailBinding? = null
    private val binding get() = _binding!!

    private val args: DetailFragmentArgs by navArgs()

    private lateinit var viewModel: DetailViewModel
    private lateinit var chapterAdapter: ChapterAdapter

    private var rewardedAdLoader: RewardedAdLoader? = null
    private var rewardedAd: RewardedAd? = null

    private val adUnitId = "demo-rewarded-yandex"

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentDetailBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        val manga = args.manga

        binding.textViewTitle.text = manga.title
        binding.textViewDescription.text = manga.description
        binding.textViewRating.text = "Рейтинг: ${manga.rating}"

        Glide.with(this)
            .load(manga.thumbnailUrl)
            .placeholder(android.R.color.darker_gray)
            .into(binding.imageViewThumbnail)

        viewModel = ViewModelProvider(
            this,
            ViewModelProvider.AndroidViewModelFactory.getInstance(requireActivity().application)
        )[DetailViewModel::class.java]

        chapterAdapter = ChapterAdapter(emptyList())
        binding.rvChapters.layoutManager = LinearLayoutManager(requireContext())
        binding.rvChapters.adapter = chapterAdapter

        viewModel.chapters.observe(viewLifecycleOwner) { chapters ->
            chapterAdapter.updateData(chapters)
        }

        viewModel.error.observe(viewLifecycleOwner) { errorMsg ->
            Toast.makeText(requireContext(), errorMsg ?: "Ошибка загрузки глав", Toast.LENGTH_SHORT).show()
        }

        viewModel.status.observe(viewLifecycleOwner) { status ->
            if (status == "ready") {
                viewModel.task.value?.let {
                    viewModel.downloadPdf(it, manga.title)
                    Toast.makeText(requireContext(), "Загрузка завершена", Toast.LENGTH_SHORT).show()
                }

            }
        }

        viewModel.loadChapters(manga)

        loadRewardedAd()

        binding.btnDownloadPdf.setOnClickListener {
            val selectedChapters = chapterAdapter.getSelectedChapters()

            if (selectedChapters.isEmpty()) {
                Toast.makeText(requireContext(), "Выберите хотя бы одну главу", Toast.LENGTH_SHORT).show()
                return@setOnClickListener
            }

            val selectedPaths = selectedChapters.map { it.path }

            showRewardedAd {
                viewModel.downloadManga(manga, selectedPaths)
            }
        }
    }

    private fun loadRewardedAd() {
        rewardedAdLoader = RewardedAdLoader(requireContext()).apply {
            setAdLoadListener(object : RewardedAdLoadListener {
                override fun onAdLoaded(ad: RewardedAd) {
                    rewardedAd = ad
                }

                override fun onAdFailedToLoad(adRequestError: AdRequestError) {
                    Toast.makeText(requireContext(), "Ошибка загрузки рекламы: ${adRequestError.description}", Toast.LENGTH_SHORT).show()
                }
            })
        }

        val adRequestConfiguration = AdRequestConfiguration.Builder(adUnitId).build()
        rewardedAdLoader?.loadAd(adRequestConfiguration)
    }


    private fun showRewardedAd(onReward: () -> Unit) {
        rewardedAd?.apply {
            setAdEventListener(object : RewardedAdEventListener {
                override fun onAdShown() {}

                override fun onAdDismissed() {
                    rewardedAd?.setAdEventListener(null)
                    rewardedAd = null
                    loadRewardedAd()

                    binding.root.requestLayout()
                    binding.root.invalidate()
                    ViewCompat.requestApplyInsets(binding.root)
                }

                override fun onAdFailedToShow(adError: AdError) {
                    rewardedAd?.setAdEventListener(null)
                    rewardedAd = null
                    Toast.makeText(requireContext(), "Не удалось показать рекламу", Toast.LENGTH_SHORT).show()
                    loadRewardedAd()
                }

                override fun onAdClicked() {}
                override fun onAdImpression(impressionData: ImpressionData?) {}

                override fun onRewarded(reward: Reward) {
                    onReward()
                }
            })
            show(requireActivity())
        } ?: run {
            Toast.makeText(requireContext(), "Реклама не загружена. Скачивание начнётся сразу.", Toast.LENGTH_SHORT).show()
            onReward()
            loadRewardedAd()
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        rewardedAd?.setAdEventListener(null)
        rewardedAd = null
        rewardedAdLoader?.setAdLoadListener(null)
        rewardedAdLoader = null
        _binding = null
    }
}
