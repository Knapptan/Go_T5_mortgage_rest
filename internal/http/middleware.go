// Package http часть с middleware для логов
package http

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("status_code: %d, duration: %d ns\n", c.Writer.Status(), duration.Nanoseconds())
	}
}
