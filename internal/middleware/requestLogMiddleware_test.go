package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/plyovchev/sumup-assignment-notifications/internal/config"
	"github.com/plyovchev/sumup-assignment-notifications/internal/logger"
	"github.com/plyovchev/sumup-assignment-notifications/internal/middleware"
)

func TestRequestLogMiddleware(_ *testing.T) {
	type requestLogMiddlewareTestCase struct {
		Description  string
		InputReqPath string
	}

	var testCases = []requestLogMiddlewareTestCase{
		{
			Description:  "improve assertions-1",
			InputReqPath: "/test/1",
		},
		{
			Description:  "improve assertions-2",
			InputReqPath: "/test/2",
		},
	}

	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(resp)
	r.Use(middleware.RequestLogMiddleware(logger.Setup(&config.Config{})))

	for _, tc := range testCases {
		r.GET(tc.InputReqPath, func(ctx *gin.Context) {
			ctx.String(200, "OK")
		})

		c.Request, _ = http.NewRequest(http.MethodGet, tc.InputReqPath, nil)
		r.ServeHTTP(resp, c.Request)
	}
}
