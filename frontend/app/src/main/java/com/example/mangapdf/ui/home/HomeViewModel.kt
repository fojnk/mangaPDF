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

    private val _mangaList = MutableLiveData<List<Manga>>()
    val mangaList: LiveData<List<Manga>> get() = _mangaList

    private val _error = MutableLiveData<String>()
    val error: LiveData<String> get() = _error

    fun loadManga() {
        viewModelScope.launch {
            val result = mangaRepository.getMangaList()
            result.onSuccess { mangaList ->
                _mangaList.value = mangaList
            }
            result.onFailure { exception ->
                _error.value = exception.message
            }
        }
    }
}


