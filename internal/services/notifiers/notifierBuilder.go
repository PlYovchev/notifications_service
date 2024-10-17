package notifiers

import (
	"github.com/plyovchev/sumup-assignment-notifications/internal/config"
	"github.com/plyovchev/sumup-assignment-notifications/internal/logger"
	"github.com/plyovchev/sumup-assignment-notifications/internal/models/data"
)

// An interface for sending a notification to a 3rd party service.
type Notifier interface {
	SendNotification(notification *data.Notification) error
}

// Builder function for creation of a specific notifier
// which could perform the delivery over the specified channel.
func CreateNotifierForChannel(
	deliveryChannel data.DeliveryChannel,
	config *config.Config,
	logger *logger.AppLogger,
) Notifier {
	if deliveryChannel == data.Email {
		emailSenderConfig := EmailSenderConfig{
			From:       config.Email.From,
			Password:   config.Email.Password,
			Recipients: config.Email.Recipients,
			SmtpHost:   config.Email.SmtpHost,
			SmtpPort:   config.Email.SmtpPort,
		}

		return NewEmailNotifier(emailSenderConfig, logger)
	} else if deliveryChannel == data.Slack {
		return NewSlackNotifier(config.Slack.WebhookUrl, logger)
	}

	return nil
}
