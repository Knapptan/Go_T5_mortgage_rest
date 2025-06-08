// Package http содержит настройку маршрутов и определения интерфейсов обработчиков.
package http

import (
	"github.com/gin-gonic/gin"
)

// HandlerInterface определяет интерфейс для HTTP-обработчиков маршрутов.
type HandlerInterface interface {
	Execute(c *gin.Context)
	GetCache(c *gin.Context)
}

// SetupRoutes регистрирует маршруты Gin и middleware для указанного обработчика.
func SetupRoutes(router *gin.Engine, handler HandlerInterface) {
	router.Use(LoggingMiddleware())
	router.POST("/execute", handler.Execute)
	router.GET("/cache", handler.GetCache)
}
