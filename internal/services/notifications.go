package services

import (
	"github.com/plyovchev/sumup-assignment-notifications/internal/config"
	"github.com/plyovchev/sumup-assignment-notifications/internal/logger"
	"github.com/plyovchev/sumup-assignment-notifications/internal/models/data"
	"github.com/plyovchev/sumup-assignment-notifications/internal/models/external"
	"github.com/plyovchev/sumup-assignment-notifications/internal/services/notifiers"
)

type Notifier interface {
	SendNotification(notification *data.Notification) error
}

type NotificationsService interface {
	createNotifierForNotificationType(deliveryChannel external.DeliveryChannel) Notifier
	SendNotification(notification *data.Notification) error
}

type notificationService struct {
	config *config.Config
	logger *logger.AppLogger
}

func NewNotificationService(config *config.Config, logger *logger.AppLogger) NotificationsService {
	return &notificationService{
		config: config,
		logger: logger,
	}
}

func (notificationService *notificationService) createNotifierForNotificationType(deliveryChannel external.DeliveryChannel) Notifier {
	if deliveryChannel == external.Email {
		emailSenderConfig := notifiers.EmailSenderConfig{
			From:       notificationService.config.Email.From,
			Password:   notificationService.config.Email.Password,
			Recipients: notificationService.config.Email.Recipients,
			SmtpHost:   notificationService.config.Email.SmtpHost,
			SmtpPort:   notificationService.config.Email.SmtpPort,
		}

		return notifiers.NewEmailNotifier(emailSenderConfig, notificationService.logger)
	} else if deliveryChannel == external.Slack {
		return notifiers.NewSlackNotifier(notificationService.config.Slack.WebhookUrl, notificationService.logger)
	}

	return nil
}

func (notificationService *notificationService) SendNotification(notification *data.Notification) error {
	for _, deliveryChannel := range notification.DeliveryChannels {
		notifier := notificationService.createNotifierForNotificationType(deliveryChannel)

		if err := notifier.SendNotification(notification); err != nil {
			notificationService.logger.Error().Err(err).Msgf("The notification with key '%s' could not be sent!", notification.Key)
			return err
		}
	}

	return nil
}
