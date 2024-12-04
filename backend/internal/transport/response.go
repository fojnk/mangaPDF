package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type transort_error struct {
	Message string `json:"message"`
}

func NewTransportErrorResponse(c *gin.Context, status_code int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(status_code, transort_error{message})
}
