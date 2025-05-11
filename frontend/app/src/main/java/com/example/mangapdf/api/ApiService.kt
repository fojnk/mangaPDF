package com.example.mangapdf.api

import com.example.mangapdf.models.AuthResponse
import com.example.mangapdf.models.LoginRequest
import com.example.mangapdf.models.MangaResponse
import com.example.mangapdf.models.RegisterRequest
import retrofit2.Call
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.POST

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
    suspend fun getMangaList(): MangaResponse

}

