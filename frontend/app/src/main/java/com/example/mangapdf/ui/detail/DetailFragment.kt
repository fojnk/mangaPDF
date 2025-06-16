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

class DetailFragment : Fragment() {

    private var _binding: FragmentDetailBinding? = null
    private val binding get() = _binding!!

    private val args: DetailFragmentArgs by navArgs()

    private lateinit var viewModel: DetailViewModel
    private lateinit var chapterAdapter: ChapterAdapter

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

        viewModel.loadChapters(manga)

        binding.btnDownloadPdf.setOnClickListener {
            Toast.makeText(requireContext(), "Нажали скачать PDF для ${manga.title}", Toast.LENGTH_SHORT).show()
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
