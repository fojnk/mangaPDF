package transport

import (
	"encoding/json"
	"net/http"

	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
// @Router /api/v1/manga/download [post]
func (h *Handler) downloadMangaChapters(c *gin.Context) {
	mangaId := c.Query("manga_id")

	var input models.DownloadOpts

	if err := c.BindJSON(&input); err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	out := h.services.TaskService.CreateTask(input, mangaId)

	c.JSON(http.StatusOK, out)
}

// @Summary Download Manga Status
// @Security ApiKeyAuth
// @Tags Manga
// @Description Download Manga Status
// @ID manga-download-status
// @Param task_id query string true "task_id"
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /api/v1/manga/status [get]
func (h *Handler) downloadStatus(c *gin.Context) {
	taskId := c.Query("task_id")
	status, ok := h.services.TaskService.GetStatus(taskId)

	logrus.Info(status)

	if ok {
		c.JSON(http.StatusOK, status)
	} else {
		c.JSON(http.StatusNotFound, "task not found")
	}
}

// @Summary Download Manga Result
// @Security ApiKeyAuth
// @Tags Manga
// @Description Download Manga Result
// @ID manga-download-result
// @Param task_id query string true "task_id"
// @Produce application/pdf
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /api/v1/manga/result [get]
func (h *Handler) downloadResult(c *gin.Context) {
	taskId := c.Query("task_id")

	status, ok := h.services.TaskService.GetStatus(taskId)

	logrus.Info(status)

	if ok {
		if status.Status != "ready" {
			c.JSON(http.StatusAccepted, gin.H{"error": "not ready"})
			return
		}

		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=\"result.pdf\"")
		c.File(status.File)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}
