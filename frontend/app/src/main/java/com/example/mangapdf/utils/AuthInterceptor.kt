package com.example.mangapdf.utils

import com.example.mangapdf.data.local.SecurePreferencesManager
import okhttp3.Interceptor
import okhttp3.Response

class AuthInterceptor : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val token = SecurePreferencesManager.geAccessToken()
        val request = chain.request().newBuilder()
            .addHeader("Authorization", "Bearer $token")
            .build()
        return chain.proceed(request)
    }
}
