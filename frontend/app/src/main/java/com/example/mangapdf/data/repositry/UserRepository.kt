import com.example.mangapdf.models.AuthResponse
import com.example.mangapdf.models.LoginRequest
import com.example.mangapdf.models.RegisterRequest
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response
import java.net.NetworkInterface

class UserRepository {
    private val apiService = RetrofitInstance.api

    /**
     * Полчение ip адреса
     */
    public fun getLocalIpAddress(): String? {
        try {
            val networkInterfaces = NetworkInterface.getNetworkInterfaces()
            while (networkInterfaces.hasMoreElements()) {
                val networkInterface = networkInterfaces.nextElement()
                val inetAddresses = networkInterface.inetAddresses
                while (inetAddresses.hasMoreElements()) {
                    val inetAddress = inetAddresses.nextElement()
                    if (!inetAddress.isLoopbackAddress) {
                        return inetAddress.hostAddress
                    }
                }
            }
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return null
    }


    fun login(username: String, password: String, callback: (result: Result<AuthResponse?>) -> Unit) {
        val loginRequest = LoginRequest(username, password)
        apiService.login(getLocalIpAddress().toString(), loginRequest).enqueue(object : Callback<AuthResponse> {
            override fun onResponse(call: Call<AuthResponse>, response: Response<AuthResponse>) {
                if (response.isSuccessful) {
                    callback(Result.success(response.body()))
                } else {
                    val errorMessage = response.errorBody()?.string() ?: "Unknown error"
                    callback(Result.failure(Exception(errorMessage)))
                }
            }
            override fun onFailure(call: Call<AuthResponse>, t: Throwable) {
                callback(Result.failure(t))
            }
        })
    }


    fun register(email: String, username: String, password: String, callback: (result: Result<AuthResponse?>) -> Unit) {
        val registerRequest = RegisterRequest(username, email, password)
        apiService.register(getLocalIpAddress().toString(), registerRequest).enqueue(object : Callback<AuthResponse> {
            override fun onResponse(call: Call<AuthResponse>, response: Response<AuthResponse>) {
                if (response.isSuccessful) {
                    callback(Result.success(response.body()))
                } else {
                    val errorMessage = response.errorBody()?.string() ?: "Unknown error"
                    callback(Result.failure(Exception(errorMessage)))
                }
            }
            override fun onFailure(call: Call<AuthResponse>, t: Throwable) {
                callback(Result.failure(t))
            }
        })
    }
}