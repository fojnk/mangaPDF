package com.example.mangapdf.models

data class RegisterResponse(
    val accessToken: String,
    val id: String,
    val refreshToken: String
)