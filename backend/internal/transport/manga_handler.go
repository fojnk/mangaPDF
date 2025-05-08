package transport

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Manga List
// @Security ApiKeyAuth
// @Tags Manga
// @Description Get Manga List
// @ID magna
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

	var js map[string]interface{}
	if err := json.Unmarshal(resp.Bytes(), &js); err != nil {
		NewTransportErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, js)
}
