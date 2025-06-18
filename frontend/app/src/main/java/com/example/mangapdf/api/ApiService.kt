package com.example.mangapdf.api

import MangaDto
import com.example.mangapdf.models.AuthResponse
import com.example.mangapdf.models.ChapterResponse
import com.example.mangapdf.models.DownloadOps
import com.example.mangapdf.models.LoginRequest
import com.example.mangapdf.models.RegisterRequest
import com.example.mangapdf.models.StatusResponse
import okhttp3.ResponseBody
import retrofit2.Call
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.POST
import retrofit2.http.Query

interface ApiService {

    @POST("auth/login")
    fun login(
        @Header("Ip") ip: String,
        @Body loginRequest: LoginRequest,
    ): Call<AuthResponse>

    @POST("auth/register")
    fun register(
        @Header("Ip") ip: String,
        @Body registerRequest: RegisterRequest
    ): Call<AuthResponse>



    @GET("/api/v1/manga/list")
    suspend fun getMangaList(
        @Query("offset") offset: Int
    ): List<MangaDto>


    @GET("/api/v1/manga/chapters")
    suspend fun getChapters(
        @Query("manga_id") mangaId: String
    ): ChapterResponse



    @POST("/api/v1/manga/download")
    suspend fun downloadManga(
        @Body downloadOps: DownloadOps,
    ): String



    @GET("/api/v1/manga/status")
    suspend fun getStatus(
        @Query("task_id") mangaId: String
    ): StatusResponse

    @GET("/api/v1/manga/result")
    suspend fun getResult(
        @Query("task_id") mangaId: String
    ): Response<ResponseBody>

}

