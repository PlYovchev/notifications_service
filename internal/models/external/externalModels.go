package external

// APIError represents the structure of an API error response.
type APIError struct {
	HTTPStatusCode int    `json:"httpStatusCode"`
	Message        string `json:"message"`
	DebugID        string `json:"debugId"`
	ErrorCode      string `json:"errorCode"`
}

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

type NotificationInput struct {
	Key              string            `json:"Key"`
	Type             NotificationType  `json:"type" binding:"required"`
	Message          string            `json:"message" binding:"required"`
	DeliveryChannels []DeliveryChannel `json:"deliveryChannels"`
}
