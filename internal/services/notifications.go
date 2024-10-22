package services

import (
	"sync"
	"time"

	"github.com/plyovchev/notifications-service/internal/config"
	"github.com/plyovchev/notifications-service/internal/logger"
	"github.com/plyovchev/notifications-service/internal/models/data"
	"github.com/plyovchev/notifications-service/internal/repositories"
	"github.com/plyovchev/notifications-service/internal/services/notifiers"
	"github.com/plyovchev/notifications-service/internal/util"
)

const (
	notificationServicePollingTime = 30 * time.Second
	retryAttempts                  = 3
	channelBufferSize              = 10
)

type NotificationsService interface {
	SendNotification(notification *data.Notification) error
	OnNotificationsReceived(notificationIds []int)
	StartNotificationService()
}

type notificationService struct {
	config                       *config.Config
	logger                       *logger.AppLogger
	notificationRepository       repositories.NotificationRepository
	receivedNotificationsChannel chan []int
	isNotificationChannelOpen    bool
	lock                         sync.Mutex
}

func NewNotificationService(
	repository repositories.NotificationRepository,
	config *config.Config,
	logger *logger.AppLogger,
) NotificationsService {
	return &notificationService{
		notificationRepository:    repository,
		config:                    config,
		logger:                    logger,
		isNotificationChannelOpen: false,
	}
}

// A hook which to wake the service's polling thread and notify it that new notifications arrived
// and should be processed.
func (service *notificationService) OnNotificationsReceived(notificationIds []int) {
	service.lock.Lock()
	if service.isNotificationChannelOpen {
		service.receivedNotificationsChannel <- notificationIds
	}

	service.lock.Unlock()
}

// Start the notification service observer functionality.
// The observer functionality waits for notificationIds to arrive over a channel,
// or executes after a specified period/timeout to process and send all pending notifications.
func (service *notificationService) StartNotificationService() {
	service.logger.Info().Msg("Notification service observer started")

	service.lock.Lock()
	{
		service.receivedNotificationsChannel = make(chan []int, channelBufferSize)
		service.isNotificationChannelOpen = true
	}
	service.lock.Unlock()

	go func(receivedNotificationChannel chan []int) {
		for {
			select {
			case receivedNotificationIds := <-receivedNotificationChannel:
				service.processPendingNotifications(receivedNotificationIds)
			case <-time.After(notificationServicePollingTime):
				service.processPendingNotifications(nil)
			}
		}
	}(service.receivedNotificationsChannel)
}

// Stops the notification service observer functionality.
func (service *notificationService) StopNotificationService() {
	service.lock.Lock()
	{
		service.isNotificationChannelOpen = false
		close(service.receivedNotificationsChannel)
	}
	service.lock.Unlock()
}

// Process any pending notifications.
//
// The notificationIds are ids of the notifications that should be processed if they are pending.
// The notificationIds could be nil in which case all stored pending notifications are processed.
func (service *notificationService) processPendingNotifications(notificationIds []int) {
	service.logger.Debug().Msg("Processing pending notifications started")

	var notifications *[]data.Notification
	var err error

	if len(notificationIds) > 0 {
		notifications, err = service.notificationRepository.FindAllByIds(notificationIds)
	} else {
		notifications, err = service.notificationRepository.FindAllByStatus(data.Pending)
	}

	if err != nil {
		var idsArr = ""
		if notificationIds != nil {
			idsArr = util.ArrayToString(notificationIds, ", ")
		}

		service.logger.Error().
			Err(err).
			Str("notificationIds", idsArr).
			Msg("Could not retrieve notifications data")
		return
	}

	if len(*notifications) == 0 {
		service.logger.Debug().Msg("Processing pending notification finished due to 0 pending notification")
		return
	}

	retryNotificationIds := make(map[int]bool)
	for i := 0; i < retryAttempts; i++ {
		for _, notification := range *notifications {
			shouldRetry, present := retryNotificationIds[notification.Id]

			// If there is not record in the retry map about this notification then send it;
			// If there is a record and it indicates that the notification should be sent again, then retry.
			if (present && !shouldRetry) || notification.Status != data.Pending {
				continue
			}

			err = service.SendNotification(&notification)

			if err != nil {
				service.logger.Error().
					Err(err).
					Int("notificationId", notification.Id).
					Msg("Sending notification failed.")

				retryNotificationIds[notification.Id] = true
			} else {
				retryNotificationIds[notification.Id] = false
			}
		}
	}

	service.updateNotificationStatuses(*notifications, retryNotificationIds)

	service.logger.Debug().Msg("Processing pending notification finished")
}

// Update in the repository the new status of the processed notifications.
func (service *notificationService) updateNotificationStatuses(
	notifications []data.Notification,
	failedNotificationsMap map[int]bool,
) {
	for _, notification := range notifications {
		updatedNotification := &data.Notification{
			Id:              notification.Id,
			Key:             notification.Key,
			Message:         notification.Message,
			DeliveryChannel: notification.DeliveryChannel,
			CreatedAt:       notification.CreatedAt,
		}

		shouldRetry, present := failedNotificationsMap[notification.Id]
		if present && !shouldRetry {
			updatedNotification.Status = data.Completed
		} else {
			updatedNotification.Status = data.Failed
		}

		_, err := service.notificationRepository.Save(updatedNotification)
		service.logger.Error().
			Err(err).
			Int("notificationId", notification.Id).
			Msg("Failed to update notification status.")
	}
}

func (service *notificationService) SendNotification(notification *data.Notification) error {
	notifier := notifiers.CreateNotifierForChannel(notification.DeliveryChannel, service.config, service.logger)

	if err := notifier.SendNotification(notification); err != nil {
		service.logger.Error().Err(err).Msgf("The notification with key '%s' could not be sent!", notification.Key)
		return err
	}

	return nil
}
