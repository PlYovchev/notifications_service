package data

import (
	"time"

	"github.com/plyovchev/sumup-assignment-notifications/internal/db"
)

type NotificationType string

const (
	Info    NotificationType = "Info"
	Warning NotificationType = "Warning"
	Error   NotificationType = "Error"
)

type DeliveryChannel string

const (
	Email DeliveryChannel = "Email"
	Slack DeliveryChannel = "Slack"
)

type NotificationStatus string

const (
	Pending   NotificationStatus = "pending"
	Completed NotificationStatus = "completed"
	Failed    NotificationStatus = "failed"
)

type Notification struct {
	Id              int                `gorm:"primary_key" json:"id"`
	Key             string             `json:"key"`
	Message         string             `json:"message"`
	Status          NotificationStatus `json:"status"`
	DeliveryChannel DeliveryChannel    `json:"delivery_channel"`
	CreatedAt       time.Time          `json:"created_at"`
}

// TableName returns the table name of account struct and it is used by gorm.
func (Notification) TableName() string {
	return db.SCHEMA + "." + db.NOTIFICATION_TABLE
}

// NewAccount is constructor.
func NewNotification(key string, message string, status NotificationStatus, deliveryChannel DeliveryChannel) *Notification {
	return &Notification{Key: key, Message: message, Status: status, DeliveryChannel: deliveryChannel}
}

func (notification *Notification) ToString() string {
	return notification.Key + " " + notification.Message + " " + string(notification.DeliveryChannel)
}
