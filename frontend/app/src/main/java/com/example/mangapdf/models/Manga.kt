package com.example.mangapdf.models

import android.os.Parcelable
import kotlinx.parcelize.Parcelize

@Parcelize
data class Manga(
    val id: String,
    val title: String,
    val rating: Double,
    val thumbnailUrl: String,
    val description: String
) : Parcelable
