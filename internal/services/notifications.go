package services

import (
	"time"

	"github.com/plyovchev/sumup-assignment-notifications/internal/config"
	"github.com/plyovchev/sumup-assignment-notifications/internal/logger"
	"github.com/plyovchev/sumup-assignment-notifications/internal/models/data"
	"github.com/plyovchev/sumup-assignment-notifications/internal/repositories"
	"github.com/plyovchev/sumup-assignment-notifications/internal/services/notifiers"
	"github.com/plyovchev/sumup-assignment-notifications/internal/util"
)

const (
	notificationServicePollingTimeInSeconds = 30 * time.Second
	retryAttempts                           = 3
	channelBufferSize                       = 10
)

type Notifier interface {
	SendNotification(notification *data.Notification) error
}

type NotificationsService interface {
	createNotifierForNotificationType(deliveryChannel data.DeliveryChannel) Notifier
	SendNotification(notification *data.Notification) error
	OnNotificationsReceived() error
	StartNotificationService()
}

type notificationService struct {
	config                       *config.Config
	logger                       *logger.AppLogger
	notificationRepository       repositories.NotificationRepository
	receivedNotificationsChannel chan []int
}

func NewNotificationService(repository repositories.NotificationRepository, config *config.Config, logger *logger.AppLogger) NotificationsService {
	receivedNotificationsChannel := make(chan []int, channelBufferSize)

	return &notificationService{
		notificationRepository:       repository,
		config:                       config,
		logger:                       logger,
		receivedNotificationsChannel: receivedNotificationsChannel,
	}
}

func (service *notificationService) OnNotificationsReceived() error {

}

func (service *notificationService) StartNotificationService() {
	go func(receivedNotificationChannel chan []int) {
		for true {
			select {
			case receivedNotificationIds := <-receivedNotificationChannel:
				service.processPendingNotifications(receivedNotificationIds)
			case <-time.After(notificationServicePollingTimeInSeconds):
				service.processPendingNotifications(nil)
			}

		}
	}(service.receivedNotificationsChannel)
}

func (service *notificationService) processPendingNotifications(notificationIds []int) {
	var notifications *[]data.Notification
	var err error

	if len(notificationIds) > 0 {
		notifications, err = service.notificationRepository.FindAllByIds(notificationIds)
	} else {
		notifications, err = service.notificationRepository.FindAllByStatus(data.Pending)
	}

	if err != nil {
		var idsArr string = ""
		if notificationIds != nil {
			idsArr = util.ArrayToString(notificationIds, ", ")
		}

		service.logger.Error().
			Err(err).
			Str("notificationIds", idsArr).
			Msg("Could not retrieve notifications data")
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

	for _, notification := range *notifications {
		updatedNotification := &data.Notification{
			Id:              notification.Id,
			Key:             notification.Key,
			Message:         notification.Message,
			DeliveryChannel: notification.DeliveryChannel,
			CreatedAt:       notification.CreatedAt,
		}

		shouldRetry, present := retryNotificationIds[notification.Id]
		if present && !shouldRetry {
			updatedNotification.Status = data.Completed
		} else {
			updatedNotification.Status = data.Failed
		}

		service.notificationRepository.Save(updatedNotification)
	}
}

func (service *notificationService) createNotifierForNotificationType(deliveryChannel data.DeliveryChannel) Notifier {
	if deliveryChannel == data.Email {
		emailSenderConfig := notifiers.EmailSenderConfig{
			From:       service.config.Email.From,
			Password:   service.config.Email.Password,
			Recipients: service.config.Email.Recipients,
			SmtpHost:   service.config.Email.SmtpHost,
			SmtpPort:   service.config.Email.SmtpPort,
		}

		return notifiers.NewEmailNotifier(emailSenderConfig, service.logger)
	} else if deliveryChannel == data.Slack {
		return notifiers.NewSlackNotifier(service.config.Slack.WebhookUrl, service.logger)
	}

	return nil
}

func (service *notificationService) SendNotification(notification *data.Notification) error {
	notifier := service.createNotifierForNotificationType(notification.DeliveryChannel)

	if err := notifier.SendNotification(notification); err != nil {
		service.logger.Error().Err(err).Msgf("The notification with key '%s' could not be sent!", notification.Key)
		return err
	}

	return nil
}
