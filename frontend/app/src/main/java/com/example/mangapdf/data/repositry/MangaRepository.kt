package com.example.mangapdf.data.repositry

import RetrofitInstance
import android.os.Environment
import com.example.mangapdf.models.Chapter
import com.example.mangapdf.models.DownloadOps
import com.example.mangapdf.models.Manga
import java.io.File
import java.io.FileOutputStream

class MangaRepository {

    private val apiService = RetrofitInstance.api

    suspend fun getMangaList(offset: Int): Result<List<Manga>> {
        return try {
            val response = apiService.getMangaList(offset)
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
            val response = apiService.downloadManga(
                DownloadOps(
                    chapters = chapters,
                    manga_id = manga.id,
                    type = "chapters",
                )
            )
            Result.success(response)
        } catch (exception: Exception) {
            Result.failure(exception)
        }
    }

    suspend fun getStatus(task: String): Result<String> {
        return try {
            val response = apiService.getStatus(task)
            Result.success(response.Status)
        } catch (exception: Exception) {
            Result.failure(exception)
        }
    }


    suspend fun getPDF(task: String, mangaTitle: String): File? {
        return try {
            val response = apiService.getResult(task)

            if (response.isSuccessful && response.body() != null) {
                val body = response.body()!!

                val fileName = "${mangaTitle}_${task}.pdf"
                val downloadsFolder = Environment.getExternalStoragePublicDirectory(Environment.DIRECTORY_DOWNLOADS)
                val outputFile = File(downloadsFolder, fileName)

                body.byteStream().use { input ->
                    FileOutputStream(outputFile).use { output ->
                        input.copyTo(output)
                    }
                }
                outputFile
            } else {
                null
            }
        } catch (e: Exception) {
            e.printStackTrace()
            null
        }
    }



}


