package com.example.mangapdf.data.repositry

import com.example.mangapdf.models.Chapter
import com.example.mangapdf.models.DownloadOps
import com.example.mangapdf.models.Manga

class MangaRepository {

    private val apiService = RetrofitInstance.api

    suspend fun getMangaList(): Result<List<Manga>> {
        return try {
            val response = apiService.getMangaList()
            val mangaList = response.map {
                Manga(
                    id = it.id,
                    title = it.title,
                    rating = it.rating.toDoubleOrNull() ?: 0.0,
                    thumbnailUrl = it.image,
                    description = it.description
                )
            }
            Result.success(mangaList)
        } catch (exception: Exception) {
            Result.failure(exception)
        }
    }

    suspend fun getChapters(manga: Manga): Result<List<Chapter>> {
        return try {
            val response = apiService.getChapters(manga.id)
            Result.success(response.payload)
        } catch (exception: Exception) {
            Result.failure(exception)
        }
    }

    suspend fun downloadManga(manga: Manga, chapters: List<String>): Result<String> {
        return try {
            val chaptersString = "[${chapters.joinToString(", ") { "\"$it\"" }}]"
            val response = apiService.downloadManga(
                DownloadOps(
                    chapters = chaptersString,
                    manga_id = manga.id,
                    type = "chapters",
                )
            )
            Result.success(response)
        } catch (exception: Exception) {
            Result.failure(exception)
        }
    }
}


