package notifiers

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/plyovchev/sumup-assignment-notifications/internal/logger"
	"github.com/plyovchev/sumup-assignment-notifications/internal/models/data"
)

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

	req, _ := http.NewRequest("POST", notifier.webhookUrl, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return err
	}

	defer resp.Body.Close()

	return nil
}
