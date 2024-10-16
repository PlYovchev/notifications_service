package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plyovchev/sumup-assignment-notifications/internal/config"
	"github.com/plyovchev/sumup-assignment-notifications/internal/db"
	"github.com/plyovchev/sumup-assignment-notifications/internal/errors"
	"github.com/plyovchev/sumup-assignment-notifications/internal/logger"
	"github.com/plyovchev/sumup-assignment-notifications/internal/models/data"
	"github.com/plyovchev/sumup-assignment-notifications/internal/models/external"
	"github.com/plyovchev/sumup-assignment-notifications/internal/repositories"
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

	dbClient := db.NewDbClient(db.SCHEMA, lgr, handler.config)
	notificationRepository := repositories.NewNotificationRepository(dbClient)

	notifications := createNotificationFromInput(notificationInput)
	for _, notification := range notifications {
		if _, err := notificationRepository.Create(&notification); err != nil {
			dbApiErr := &external.APIError{
				HTTPStatusCode: http.StatusInternalServerError,
				ErrorCode:      errors.FailedToInsertInDb,
				Message:        "Failed to insert a record in the database.",
				DebugID:        requestId,
			}

			lgr.Error().
				Err(err).
				Int("HttpStatusCode", dbApiErr.HTTPStatusCode).
				Str("ErrorCode", dbApiErr.ErrorCode).
				Msg(dbApiErr.Message)

			ginContext.AbortWithStatusJSON(dbApiErr.HTTPStatusCode, dbApiErr)
			return
		}
	}

	// notificationService := services.NewNotificationService(handler.config, lgr)
	// notificationService.SendNotification(notification)
}

func createNotificationFromInput(notificationInput external.NotificationInput) []data.Notification {
	if len(notificationInput.DeliveryChannels) == 0 {
		return nil
	}

	var notifications []data.Notification = make([]data.Notification, len(notificationInput.DeliveryChannels))
	for i, deliveryChannel := range notificationInput.DeliveryChannels {
		notifications[i] = data.Notification{
			Key:             notificationInput.Key,
			Message:         notificationInput.Message,
			DeliveryChannel: deliveryChannel,
			Status:          data.Pending,
		}
	}

	return notifications
}
