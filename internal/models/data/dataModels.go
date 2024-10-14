package data

import (
	"time"

	"github.com/plyovchev/sumup-assignment-notifications/internal/models/external"
)

// Notification service

type Notification struct {
	external.NotificationInput
	CreatedAt time.Time
}

/**
 * &Notification{Name: "payment_cancelled", Type: "", Message: "Your payment has been cancelled", DeliveryChanges: []DeliveryChannel{Email, Slack}}
 */
