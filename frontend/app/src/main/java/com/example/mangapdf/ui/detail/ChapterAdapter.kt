package com.example.mangapdf.ui.detail


import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.example.mangapdf.databinding.ItemChapterBinding
import com.example.mangapdf.models.Chapter

class ChapterAdapter(
    private var chapters: List<Chapter>
) : RecyclerView.Adapter<ChapterAdapter.ChapterViewHolder>() {

    inner class ChapterViewHolder(val binding: ItemChapterBinding) : RecyclerView.ViewHolder(binding.root)

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ChapterViewHolder {
        val binding = ItemChapterBinding.inflate(LayoutInflater.from(parent.context), parent, false)
        return ChapterViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ChapterViewHolder, position: Int) {
        val chapter = chapters[position]

        holder.binding.tvChapterNumber.text = chapter.title
        holder.binding.cbSelected.isChecked = chapter.isSelected

        holder.binding.cbSelected.setOnCheckedChangeListener(null)
        holder.binding.cbSelected.setOnCheckedChangeListener { _, isChecked ->
            chapter.isSelected = isChecked
        }
    }

    override fun getItemCount() = chapters.size

    fun updateData(newChapters: List<Chapter>) {
        chapters = newChapters
        notifyDataSetChanged()
    }

    fun getSelectedChapters(): List<Chapter> {
        return chapters.filter { it.isSelected }
    }
}

