package com.example.mangapdf.ui.home

import android.app.Application
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.viewModelScope
import com.example.mangapdf.data.repositry.MangaRepository
import com.example.mangapdf.models.Manga
import kotlinx.coroutines.launch

class HomeViewModel(application: Application) : AndroidViewModel(application) {

    private val mangaRepository = MangaRepository()

    private val _mangaList = MutableLiveData<List<Manga>>(emptyList())
    val mangaList: LiveData<List<Manga>> get() = _mangaList

    private val _error = MutableLiveData<String>()
    val error: LiveData<String> get() = _error

    private var isLoading = false
    private var offset = 0
    private val limit = 50
    private var allDataLoaded = false

    fun loadMoreManga() {
        if (isLoading || allDataLoaded) return

        isLoading = true
        viewModelScope.launch {
            val result = mangaRepository.getMangaList(offset)
            result.onSuccess { newList ->
                val currentList = _mangaList.value ?: emptyList()
                _mangaList.value = currentList + newList
                offset += newList.size
                if (newList.size < limit) {
                    allDataLoaded = true
                }
            }.onFailure {
                _error.value = it.message
            }
            isLoading = false
        }
    }
}
