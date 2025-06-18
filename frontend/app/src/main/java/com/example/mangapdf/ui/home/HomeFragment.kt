package com.example.mangapdf.ui.home

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.fragment.app.Fragment
import androidx.lifecycle.ViewModelProvider
import androidx.navigation.fragment.findNavController
import androidx.recyclerview.widget.GridLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.example.mangapdf.databinding.FragmentHomeBinding

class HomeFragment : Fragment() {

    private var _binding: FragmentHomeBinding? = null
    private val binding get() = _binding!!

    private lateinit var homeViewModel: HomeViewModel
    private lateinit var mangaAdapter: MangaAdapter

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        _binding = FragmentHomeBinding.inflate(inflater, container, false)
        val root: View = binding.root

        homeViewModel = ViewModelProvider(this)[HomeViewModel::class.java]

        mangaAdapter = MangaAdapter { manga ->
            val action = HomeFragmentDirections.actionHomeToDetail(manga)
            findNavController().navigate(action)
        }

        val layoutManager = GridLayoutManager(context, 2)
        binding.recyclerView.layoutManager = layoutManager
        binding.recyclerView.adapter = mangaAdapter

        binding.recyclerView.addOnScrollListener(object : RecyclerView.OnScrollListener() {
            override fun onScrolled(recyclerView: RecyclerView, dx: Int, dy: Int) {
                super.onScrolled(recyclerView, dx, dy)

                val totalItemCount = layoutManager.itemCount
                val lastVisibleItem = layoutManager.findLastVisibleItemPosition()

                if (totalItemCount <= lastVisibleItem + 5) {
                    homeViewModel.loadMoreManga()
                }
            }
        })

        homeViewModel.mangaList.observe(viewLifecycleOwner) { list ->
            mangaAdapter.submitList(list)
        }

        homeViewModel.error.observe(viewLifecycleOwner) { error ->
            error?.let {
                Toast.makeText(requireContext(), "Ошибка: $it", Toast.LENGTH_LONG).show()
            }
        }

        homeViewModel.loadMoreManga() // Initial load

        return root
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}
