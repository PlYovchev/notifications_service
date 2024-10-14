package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/plyovchev/sumup-assignment-notifications/internal/config"
	"github.com/plyovchev/sumup-assignment-notifications/internal/errors"
	"github.com/plyovchev/sumup-assignment-notifications/internal/logger"
	"github.com/plyovchev/sumup-assignment-notifications/internal/models/data"
	"github.com/plyovchev/sumup-assignment-notifications/internal/models/external"
	"github.com/plyovchev/sumup-assignment-notifications/internal/services"
)

type NotificationsHandler struct {
	config *config.Config
	logger *logger.AppLogger
}

func NewNotificationsHandler(cfg *config.Config, logger *logger.AppLogger) *NotificationsHandler {
	return &NotificationsHandler{
		config: cfg,
		logger: logger,
	}
}

func (handler *NotificationsHandler) PushNotification(ginContext *gin.Context) {
	lgr, requestId := handler.logger.WithReqID(ginContext)

	var notificationInput external.NotificationInput
	if err := ginContext.ShouldBindJSON(&notificationInput); err != nil {
		apiErr := &external.APIError{
			HTTPStatusCode: http.StatusBadRequest,
			ErrorCode:      errors.PushNotificationInvalidParams,
			Message:        "Invalid push notification request body",
			DebugID:        requestId,
		}

		lgr.Error().
			Err(err).
			Int("HttpStatusCode", apiErr.HTTPStatusCode).
			Str("ErrorCode", apiErr.ErrorCode).
			Msg(apiErr.Message)

		ginContext.AbortWithStatusJSON(apiErr.HTTPStatusCode, apiErr)
		return
	}

	notification := createNotificationFromInput(notificationInput)

	notificationService := services.NewNotificationService(handler.config, lgr)
	notificationService.SendNotification(notification)
}

func createNotificationFromInput(notificationInput external.NotificationInput) *data.Notification {
	return &data.Notification{
		NotificationInput: notificationInput,
		CreatedAt:         time.Now(),
	}
}
