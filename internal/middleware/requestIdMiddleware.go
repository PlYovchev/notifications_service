package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/plyovchev/notifications-service/internal/config"
)

func ReqIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.Request.Header.Get(config.RequestIdentifier)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx := context.WithValue(c.Request.Context(), config.ContextKey(config.RequestIdentifier), reqID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(config.RequestIdentifier, reqID)
		c.Next()
	}
}
