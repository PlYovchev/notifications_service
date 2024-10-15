package server_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/plyovchev/sumup-assignment-notifications/internal/config"
	"github.com/plyovchev/sumup-assignment-notifications/internal/logger"
	"github.com/plyovchev/sumup-assignment-notifications/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestListOfRoutes(t *testing.T) {
	serviceEnv := config.ServiceEnv{Name: "test"}
	config := &config.Config{}
	lgr := logger.Setup(serviceEnv)
	router := server.WebRouter(serviceEnv, config, lgr)
	list := router.Routes()
	mode := gin.Mode()

	assert.Equal(t, gin.ReleaseMode, mode)

	assertRoutePresent(t, list, gin.RouteInfo{
		Method: http.MethodGet,
		Path:   "/status",
	})
}

func assertRoutePresent(t *testing.T, gotRoutes gin.RoutesInfo, wantRoute gin.RouteInfo) {
	for _, gotRoute := range gotRoutes {
		if gotRoute.Path == wantRoute.Path && gotRoute.Method == wantRoute.Method {
			return
		}
	}
	t.Errorf("route not found: %v", wantRoute)
}
