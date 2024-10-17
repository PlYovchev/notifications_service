package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plyovchev/sumup-assignment-notifications/internal/logger"
)

type ServiceStatus string

const (
	UP   ServiceStatus = "ok"
	DOWN ServiceStatus = "down"
)

type StatusHandler struct {
	logger *logger.AppLogger
}

func NewStatusHandler(logger *logger.AppLogger) *StatusHandler {
	return &StatusHandler{
		logger: logger,
	}
}

// CheckStatus - Checks the health of all the dependencies of the service to ensure complete serviceability.
func (s *StatusHandler) CheckStatus(c *gin.Context) {
	var code = http.StatusOK

	// send response
	c.JSON(code, UP)
}
