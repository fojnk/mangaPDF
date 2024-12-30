package com.example.mangapdf.ui.register

import android.app.Application
import androidx.lifecycle.AndroidViewModel
import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.example.mangapdf.repositry.UserRepository
import com.example.mangapdf.ui.login.LoginStatus

class RegisterViewModel(application: Application) : AndroidViewModel(application) {

    private val _registerStatus = MutableLiveData<RegisterStatus>()
    val registerStatus: LiveData<RegisterStatus> get() = _registerStatus

    private val userRepository = UserRepository()


    fun register(name: String, email: String, password: String) {
        userRepository.register(email, name, password) { result ->
            result.onSuccess {response ->
                if (response != null) {
                    _registerStatus.value = RegisterStatus(true, response.accessToken, response.refreshToken, null)
                }
            }
            result.onFailure {exception ->
                _registerStatus.value = RegisterStatus(false, null, null, exception.message)
            }
        }
    }
}

data class RegisterStatus(
    val isSuccess: Boolean,
    val accessToken: String?,
    val refreshToken: String?,
    val message: String?
)
