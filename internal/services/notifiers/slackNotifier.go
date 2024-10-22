package notifiers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/plyovchev/notifications-service/internal/logger"
	"github.com/plyovchev/notifications-service/internal/models/data"
)

const slackWebhookTimeout = 30 * time.Second

type SlackNotifier struct {
	logger     *logger.AppLogger
	webhookUrl string
}

func NewSlackNotifier(webhookUrl string, logger *logger.AppLogger) *SlackNotifier {
	return &SlackNotifier{
		logger:     logger,
		webhookUrl: webhookUrl,
	}
}

func (notifier *SlackNotifier) SendNotification(notification *data.Notification) error {
	notifier.logger.Debug().Msg("Sending slack message")

	var jsonStr = fmt.Sprintf("{\"text\":\"%s\"}", notification.Message)
	var jsonBytes = []byte(jsonStr)

	ctx, cancel := context.WithTimeout(context.Background(), slackWebhookTimeout)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, notifier.webhookUrl, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}

	defer resp.Body.Close()

	return nil
}
