package com.example.mangapdf.ui.detail

import android.app.Application
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.viewModelScope
import com.example.mangapdf.data.repositry.MangaRepository
import com.example.mangapdf.models.Chapter
import com.example.mangapdf.models.Manga
import kotlinx.coroutines.launch

class DetailViewModel(application: Application) : AndroidViewModel(application) {

    private val mangaRepository = MangaRepository()

    private val _chapters = MutableLiveData<List<Chapter>>()
    val chapters: LiveData<List<Chapter>> get() = _chapters

    private val _error = MutableLiveData<String>()
    val error: LiveData<String> get() = _error

    fun loadChapters(manga: Manga) {
        viewModelScope.launch {
            val result = mangaRepository.getChapters(manga)
            result.onSuccess { chaptersList ->
                _chapters.value = chaptersList
            }
            result.onFailure { exception ->
                _error.value = exception.message
            }
        }
    }
}
