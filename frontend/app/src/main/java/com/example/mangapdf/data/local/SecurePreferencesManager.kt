package com.example.mangapdf.data.local

import android.content.Context
import android.content.SharedPreferences
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKeys

object SecurePreferencesManager {

    private const val NAME = "secure_prefs"
    private const val ACCESS_TOKEN_KEY = "access_token"
    private const val REFRESH_TOKEN_KEY = "refresh_token"

    private lateinit var encryptedSharedPreferences: SharedPreferences

    /**
     * Инициализация EncryptedSharedPreferences.
     */
    fun initialize(context: Context) {
        if (!::encryptedSharedPreferences.isInitialized) {
            val masterKey = MasterKeys.getOrCreate(MasterKeys.AES256_GCM_SPEC)

            encryptedSharedPreferences = EncryptedSharedPreferences.create(
                NAME,
                masterKey,
                context,
                EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
                EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
            )
        }
    }

    /**
     * Сохранение токенов.
     */
    fun saveTokens(accessToken: String, refreshToken: String) {
        val editor = encryptedSharedPreferences.edit()
        editor.putString(ACCESS_TOKEN_KEY, accessToken)
        editor.putString(REFRESH_TOKEN_KEY, refreshToken)
        editor.apply()
    }

    /**
     * Получение токенов.
     */
    fun getTokens(): Pair<String?, String?> {
        val accessToken = encryptedSharedPreferences.getString(ACCESS_TOKEN_KEY, null)
        val refreshToken = encryptedSharedPreferences.getString(REFRESH_TOKEN_KEY, null)
        return Pair(accessToken, refreshToken)
    }

    /**
     * Получение access токена.
     */
    fun geAccessToken(): String? {
        return encryptedSharedPreferences.getString(ACCESS_TOKEN_KEY, null)
    }

    /**
     *  Удаление токенов.
     */
    fun clearTokens() {
        val editor = encryptedSharedPreferences.edit()
        editor.remove(ACCESS_TOKEN_KEY)
        editor.remove(REFRESH_TOKEN_KEY)
        editor.apply()
    }
}
