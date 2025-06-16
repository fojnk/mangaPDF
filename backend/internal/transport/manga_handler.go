package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

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
// @Param offset query int true "offset"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /api/v1/manga/list [get]
func (h *Handler) getManga(c *gin.Context) {
	val := c.Query("offset")
	offset, err := strconv.Atoi(val)
	if err != nil {
		offset = 0
	}

	resp, err := h.services.MangaService.GetMangaList(offset)

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

type DownloadInput struct {
	ChaptersList []string `json:"chapters"`
	MangaName    string   `json:"manga_id"`
	Type         string   `json:"type"`
}

// @Summary Download Manga
// @Security ApiKeyAuth
// @Tags Manga
// @Description Download Manga Chapters
// @ID manga-download
// @Param downloadOps body DownloadInput true "download_opt"
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /api/v1/manga/download [post]
func (h *Handler) downloadMangaChapters(c *gin.Context) {
	var input DownloadInput
	var downloadOpt models.DownloadOpts

	if err := c.BindJSON(&input); err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	downloadOpt = models.DownloadOpts{
		Chapters:  input.ChaptersList,
		MangaURL:  input.MangaName,
		PDFall:    "1",
		PDFch:     "0",
		PDFvol:    "0",
		Del:       "1",
		Type:      input.Type,
		UserHash:  "0",
		CBZ:       "0",
		SavePath:  "",
		PrefTrans: "0",
	}

	out := h.services.TaskService.CreateTask(downloadOpt, input.MangaName)

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
