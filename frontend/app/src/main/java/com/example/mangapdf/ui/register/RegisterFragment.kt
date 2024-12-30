package com.example.mangapdf.ui.register

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
import com.example.mangapdf.databinding.ActivityRegisterBinding

class RegisterFragment : Fragment() {

    private var _binding: ActivityRegisterBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        val registerViewModel =
            ViewModelProvider(this)[RegisterViewModel::class.java]

        _binding = ActivityRegisterBinding.inflate(inflater, container, false)
        val root: View = binding.root


        val etName: EditText        = binding.etName
        val etEmail: EditText       = binding.etEmail
        val etPassword: EditText    = binding.etPassword
        val btnRegister: Button     = binding.btnRegister
        val tvLogin: TextView       = binding.tvLogin;

        btnRegister.setOnClickListener {
            val name = etName.text.toString()
            val email = etEmail.text.toString()
            val password = etPassword.text.toString()

            // Регистрация на сервере.
            registerViewModel.register(name, email, password)
        }

        tvLogin.setOnClickListener {
            findNavController().navigate(R.id.navigation_login)
        }

        registerViewModel.registerStatus.observe(viewLifecycleOwner) { status ->
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


