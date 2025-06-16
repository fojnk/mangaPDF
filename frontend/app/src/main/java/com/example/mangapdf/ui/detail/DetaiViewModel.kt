package com.example.mangapdf.ui.detail

import android.app.Application
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.viewModelScope
import com.example.mangapdf.data.repositry.MangaRepository
import com.example.mangapdf.models.Chapter
import com.example.mangapdf.models.Manga
import kotlinx.coroutines.delay
import kotlinx.coroutines.launch

class DetailViewModel(application: Application) : AndroidViewModel(application) {

    private val mangaRepository = MangaRepository()

    private val _chapters = MutableLiveData<List<Chapter>>()
    val chapters: LiveData<List<Chapter>> get() = _chapters

    private val _task = MutableLiveData<String>()
    val task: LiveData<String> get() = _task

    private val _status = MutableLiveData<String>()
    val status: LiveData<String> get() = _status

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

    fun downloadManga(manga: Manga, selectedPaths: List<String>) {
        viewModelScope.launch {
            val result = mangaRepository.downloadManga(manga, selectedPaths)
            result.onSuccess { response ->
                _task.value = response

                pollUntilReady(response)
            }
            result.onFailure { exception ->
                _error.value = exception.message
            }
        }
    }

    private fun pollUntilReady(taskId: String, maxAttempts: Int = 10, delayMillis: Long = 10000L) {
        viewModelScope.launch {
            repeat(maxAttempts) {
                val result = mangaRepository.getStatus(taskId)

                result
                    .onSuccess { status ->
                        if (status == "ready") {
                            _status.value = status
                            return@launch
                        } else {
                            _status.value = status
                        }
                    }
                    .onFailure { exception ->
                        _error.value = "Ошибка при проверке статуса: ${exception.message}"
                        return@launch
                    }

                delay(delayMillis)
            }

            _error.value = "Превышено количество попыток ожидания готовности файла"
        }
    }


}
