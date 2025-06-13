package transport

import (
	"encoding/json"
	"net/http"

	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/gin-gonic/gin"
)

// @Summary Get Manga List
// @Security ApiKeyAuth
// @Tags Manga
// @Description Get Manga List
// @ID manga
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /api/v1/manga/list [get]
func (h *Handler) getManga(c *gin.Context) {
	resp, err := h.services.MangaService.GetMangaList()

	if err != nil {
		NewTransportErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Get Manga Chapters
// @Security ApiKeyAuth
// @Tags Manga
// @Description Get Manga Chapters
// @ID manga-chapters
// @Param manga_id query string true "manga_id"
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /api/v1/manga/chapters [get]
func (h *Handler) getMangaChapters(c *gin.Context) {
	mangaId := c.Query("manga_id")

	resp, err := h.services.MangaService.GetChaptersList(mangaId)

	if err != nil {
		NewTransportErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var js map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &js); err != nil {
		NewTransportErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, js)
}

// @Summary Download Manga
// @Security ApiKeyAuth
// @Tags Manga
// @Description Download Manga Chapters
// @ID manga-download
// @Param manga_id query string true "manga_id"
// @Param downloadOps body models.DownloadOpts true "download_opt"
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /api/v1/manga/download [get]
func (h *Handler) downloadMangaChapters(c *gin.Context) {
	mangaId := c.Query("manga_id")

	var input models.DownloadOpts

	if err := c.BindJSON(&input); err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	out := h.services.MangaService.DownloadManga(input, mangaId)

	c.JSON(http.StatusOK, out)
}
