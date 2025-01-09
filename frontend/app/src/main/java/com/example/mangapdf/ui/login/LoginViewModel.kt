package com.example.mangapdf.ui.login

import UserRepository
import android.app.Application
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import com.example.mangapdf.data.local.SecurePreferencesManager

class LoginViewModel(application: Application) : AndroidViewModel(application) {

    private val _loginStatus = MutableLiveData<LoginStatus>()
    val loginStatus: LiveData<LoginStatus> get() = _loginStatus

    private val userRepository = UserRepository()

    fun login(username: String, password: String) {
        userRepository.login(username, password) { result ->
            result.onSuccess {response ->
                if (response != null) {
                    SecurePreferencesManager.saveTokens(response.accessToken, response.refreshToken)
                    _loginStatus.value = LoginStatus(true, response.accessToken, response.refreshToken, null)
                }
            }
            result.onFailure {exception ->
                _loginStatus.value = LoginStatus(false, null, null, exception.message)
            }
        }
    }
}

data class LoginStatus(
    val isSuccess: Boolean,
    val accessToken: String?,
    val refreshToken: String?,
    val message: String?
)
