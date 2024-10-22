package notifiers

import (
	"net/smtp"

	"github.com/plyovchev/notifications-service/internal/logger"
	"github.com/plyovchev/notifications-service/internal/models/data"
)

type EmailSenderConfig struct {
	From       string
	Password   string
	Recipients []string
	SmtpHost   string
	SmtpPort   string
}

type EmailNotifier struct {
	EmailSenderConfig
	logger *logger.AppLogger
}

func NewEmailNotifier(emailSenderConfig EmailSenderConfig, logger *logger.AppLogger) *EmailNotifier {
	return &EmailNotifier{
		logger:            logger,
		EmailSenderConfig: emailSenderConfig,
	}
}

func (notifier *EmailNotifier) SendNotification(notification *data.Notification) error {
	notifier.logger.Debug().Msg("Sending email.")

	// Authentication.
	auth := smtp.PlainAuth("", notifier.From, notifier.Password, notifier.SmtpHost)

	err := smtp.SendMail(notifier.SmtpHost+":"+notifier.SmtpPort, auth, notifier.From, notifier.Recipients, []byte(notification.Message))

	if err != nil {
		return err
	}

	notifier.logger.Debug().Msg("Email has been sent.")

	return nil
}
