package com.example.mangapdf.ui.home

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.DiffUtil
import androidx.recyclerview.widget.ListAdapter
import androidx.recyclerview.widget.RecyclerView
import com.example.mangapdf.models.Manga
import com.example.mangapdf.databinding.MangaItemBinding
import com.bumptech.glide.Glide


class MangaAdapter : ListAdapter<Manga, MangaAdapter.MangaViewHolder>(DIFF) {

    companion object {
        val DIFF = object : DiffUtil.ItemCallback<Manga>() {
            override fun areItemsTheSame(oldItem: Manga, newItem: Manga): Boolean =
                oldItem.id == newItem.id

            override fun areContentsTheSame(oldItem: Manga, newItem: Manga): Boolean =
                oldItem == newItem
        }
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): MangaViewHolder {
            val binding = MangaItemBinding.inflate(LayoutInflater.from(parent.context), parent, false)
        return MangaViewHolder(binding)
    }

    override fun onBindViewHolder(holder: MangaViewHolder, position: Int) {
        holder.bind(getItem(position))
    }

    inner class MangaViewHolder(private val binding: MangaItemBinding) :
        RecyclerView.ViewHolder(binding.root) {

        fun bind(manga: Manga) {
            binding.title.text = manga.title
            binding.rating.text = "★ ${manga.rating}"
            Glide.with(binding.thumbnail.context)
                .load(manga.thumbnailUrl)
                .into(binding.thumbnail)

            val marginInPx = binding.root.context.resources.displayMetrics.density * 16
            val screenWidth = binding.root.context.resources.displayMetrics.widthPixels
            val cardWidth = ((screenWidth - marginInPx * 3) / 2).toInt()

            val layoutParams = binding.root.layoutParams
            layoutParams.width = cardWidth
            binding.root.layoutParams = layoutParams
        }

    }
}
