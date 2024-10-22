package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/plyovchev/notifications-service/internal/config"
	"github.com/plyovchev/notifications-service/internal/errors"
	"github.com/plyovchev/notifications-service/internal/logger"
	"github.com/plyovchev/notifications-service/internal/models/data"
	"github.com/plyovchev/notifications-service/internal/models/external"
	"github.com/plyovchev/notifications-service/internal/repositories"
	"github.com/plyovchev/notifications-service/internal/services"
	"github.com/plyovchev/notifications-service/internal/util"
)

type NotificationsHandler struct {
	config                 *config.Config
	notificationService    services.NotificationsService
	notificationRepository repositories.NotificationRepository
	logger                 *logger.AppLogger
}

func NewNotificationsHandler(
	cfg *config.Config,
	notificationService services.NotificationsService,
	notificationRepository repositories.NotificationRepository,
	logger *logger.AppLogger,
) *NotificationsHandler {
	return &NotificationsHandler{
		config:                 cfg,
		notificationService:    notificationService,
		notificationRepository: notificationRepository,
		logger:                 logger,
	}
}

// Handles a push notification request. Expects a HTTP POST request.
// The body of the request should contain an input in the form of NotificationInput.
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

	// Persist the newly created notifications from the input.
	notifications := createNotificationsFromInput(notificationInput)
	for _, notification := range notifications {
		if _, err := handler.notificationRepository.Create(notification); err != nil {
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

	// Notify the notification service that new notifications have been received.
	// The notification ids are also sent so the new notifications could be prioritized.
	notificationIds := util.Map(notifications, func(notification *data.Notification) int { return notification.Id })
	handler.notificationService.OnNotificationsReceived(notificationIds)

	ginContext.JSON(http.StatusOK, notificationIds)
}

func createNotificationsFromInput(notificationInput external.NotificationInput) []*data.Notification {
	if len(notificationInput.DeliveryChannels) == 0 {
		return nil
	}

	var notifications = make([]*data.Notification, len(notificationInput.DeliveryChannels))
	for i, deliveryChannel := range notificationInput.DeliveryChannels {
		notifications[i] = &data.Notification{
			Key:             notificationInput.Key,
			Message:         notificationInput.Message,
			DeliveryChannel: deliveryChannel,
			Status:          data.Pending,
		}
	}

	return notifications
}
