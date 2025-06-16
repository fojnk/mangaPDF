package com.example.mangapdf.ui.login

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import android.widget.EditText
import android.widget.TextView
import android.widget.Toast
import androidx.fragment.app.Fragment
import androidx.lifecycle.ViewModelProvider
import androidx.navigation.fragment.findNavController
import com.example.mangapdf.R
import com.example.mangapdf.databinding.ActivityLoginBinding

class LoginFragment : Fragment() {

    private var _binding: ActivityLoginBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        val loginViewModel =
            ViewModelProvider(this)[LoginViewModel::class.java]

        _binding = ActivityLoginBinding.inflate(inflater, container, false)
        val root: View = binding.root

        val etName: EditText = binding.etName
        val etPassword: EditText = binding.etPassword
        val btnLogin: Button = binding.btnLogin
        val tvRegister: TextView = binding.tvRegister

        btnLogin.setOnClickListener {
            val name = etName.text.toString()
            val password = etPassword.text.toString()

            loginViewModel.login(name, password)
        }

        tvRegister.setOnClickListener {
            findNavController().navigate(R.id.navigation_register)
        }

        loginViewModel.loginStatus.observe(viewLifecycleOwner) { status ->
            if (status.isSuccess) {
                findNavController().navigate(R.id.navigation_home)
            } else {
                Toast.makeText(requireContext(), status.message, Toast.LENGTH_LONG).show()
            }

        }

        return root
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}


