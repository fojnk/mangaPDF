package com.example.mangapdf.models

data class DownloadOps (
    val chapters: List<String>,
    val manga_id: String,
    val type: String = "chapters"
)
