// Package http часть с настройкой маршрутов
package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handler *Handler) {
	// router.Use(LoggingMiddleware())
	router.POST("/execute", handler.Execute)
	router.GET("/cache", handler.GetCache)
}
