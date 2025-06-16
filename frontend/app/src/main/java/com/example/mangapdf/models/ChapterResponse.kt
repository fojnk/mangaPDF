package com.example.mangapdf.models

data class ChapterResponse(
    val is_mtr:         Boolean,
    val payload:        List<Chapter>,
    val status:         String,
    val translators:    List<Translator>,
    val user_hash:      String
)

data class Chapter(
    val path:   String,
    val title:  String,
    var isSelected: Boolean = false
)

data class Translator (
    val id:     String,
    val name:   String
)