// Package http часть с настройкой маршрутов
package http

import (
	"github.com/gin-gonic/gin"
)

// HandlerInterface определяет контракт для обработчиков
type HandlerInterface interface {
	Execute(c *gin.Context)
	GetCache(c *gin.Context)
}

func SetupRoutes(router *gin.Engine, handler HandlerInterface) {
	router.Use(LoggingMiddleware())
	router.POST("/execute", handler.Execute)
	router.GET("/cache", handler.GetCache)
}
