package com.example.mangapdf.data.repositry

import com.example.mangapdf.models.Manga

class MangaRepository {

    private val apiService = RetrofitInstance.api

    suspend fun getMangaList(): Result<List<Manga>> {
        return try {
            val response = apiService.getMangaList()
            val mangaList = response.data.map {
                Manga(
                    id = it.id,
                    title = it.rus_name ?: it.eng_name,
                    rating = it.rating.averageFormated.toDoubleOrNull() ?: 0.0,
                    thumbnailUrl = it.cover.thumbnail
                )
            }
            Result.success(mangaList)
        } catch (exception: Exception) {
            Result.failure(exception)
        }
    }
}


