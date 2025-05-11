data class MangaDto(
    val id: Int,
    val rus_name: String?,
    val eng_name: String,
    val rating: RatingDto,
    val cover: CoverDto,
    val ageRestriction: AgeRestrictionDto,
    val name: String,
    val releaseDateString: String,
    val slug: String,
    val slug_url: String,
    val is_licensed: Boolean,
    val model: String,
    val site: Int,
    val status: StatusDto,
    val type: TypeDto
)

data class RatingDto(
    val average: String,
    val averageFormated: String,
    val user: Int,
    val votes: Int,
    val votesFormated: String
)

data class CoverDto(
    val default: String,
    val filename: String,
    val md: String,
    val thumbnail: String
)

data class AgeRestrictionDto(
    val id: Int,
    val label: String
)

data class StatusDto(
    val id: Int,
    val label: String
)

data class TypeDto(
    val id: Int,
    val label: String
)


