package server

import (
	"io"
	"sync"

	"github.com/gin-contrib/gzip"
	"github.com/plyovchev/sumup-assignment-notifications/internal/config"
	"github.com/plyovchev/sumup-assignment-notifications/internal/logger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/plyovchev/sumup-assignment-notifications/internal/handlers"
	"github.com/plyovchev/sumup-assignment-notifications/internal/middleware"
	"github.com/plyovchev/sumup-assignment-notifications/internal/util"
)

var startOnce sync.Once

func StartService(serviceEnv config.ServiceEnv, cfg *config.Config, lgr *logger.AppLogger) {
	startOnce.Do(func() {
		r := WebRouter(serviceEnv, cfg, lgr)
		err := r.Run(":" + serviceEnv.Port)
		if err != nil {
			panic(err)
		}
	})
}

func WebRouter(serviceEnv config.ServiceEnv, cfg *config.Config, lgr *logger.AppLogger) *gin.Engine {
	ginMode := gin.ReleaseMode
	if util.IsDevMode(serviceEnv.Name) {
		ginMode = gin.DebugMode
		gin.ForceConsoleColor()
	}
	gin.SetMode(ginMode)
	gin.EnableJsonDecoderDisallowUnknownFields()

	// Middleware
	gin.DefaultWriter = io.Discard
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.ReqIDMiddleware())
	router.Use(middleware.ResponseHeadersMiddleware())
	router.Use(middleware.RequestLogMiddleware(lgr))
	router.Use(gin.Recovery())

	internalAPIGrp := router.Group("/internal")
	internalAPIGrp.Use(middleware.AuthMiddleware())
	pprof.RouteRegister(internalAPIGrp, "pprof")
	status := handlers.NewStatusHandler(lgr)
	router.GET("/status", status.CheckStatus) // /status

	// Routes - notifications
	externalAPIGrp := router.Group("/public-api/v1")
	externalAPIGrp.Use(middleware.AuthMiddleware())
	externalAPIGrp.Use(middleware.QueryParamsCheckMiddleware(lgr))
	{
		notificationsGroup := externalAPIGrp.Group("notifications")
		{
			notifications := handlers.NewNotificationsHandler(cfg, lgr)
			notificationsGroup.POST("/push-notification", notifications.PushNotification)
		}
	}

	lgr.Info().Msg("Registered routes")
	for _, item := range router.Routes() {
		lgr.Info().
			Str("method", item.Method).
			Str("path", item.Path).
			Send()
	}
	return router
}
