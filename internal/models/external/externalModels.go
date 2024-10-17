package external

import "github.com/plyovchev/sumup-assignment-notifications/internal/models/data"

// APIError represents the structure of an API error response.
type APIError struct {
	HTTPStatusCode int    `json:"httpStatusCode"`
	Message        string `json:"message"`
	DebugID        string `json:"debugId"`
	ErrorCode      string `json:"errorCode"`
}

// The input properties of a notification request.
type NotificationInput struct {
	Key              string                 `json:"Key"`
	Message          string                 `json:"message" binding:"required"`
	DeliveryChannels []data.DeliveryChannel `json:"deliveryChannels"`
}
