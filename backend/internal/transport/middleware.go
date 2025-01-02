package transport

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationKey = "Authorization"
	UserId           = "user_id"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationKey)

	splitParts := strings.Split(header, " ")
	if len(splitParts) != 2 {
		NewTransportErrorResponse(c, http.StatusUnauthorized, "bad token (user doesn't authorized)")
		return
	}

	id, _, err := h.services.Authorization.ParseToken(splitParts[1])
	if err != nil {
		NewTransportErrorResponse(c, http.StatusUnauthorized, "parse token failed")
		return
	}

	c.Set(UserId, id)
}
