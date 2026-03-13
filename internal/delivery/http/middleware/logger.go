package middleware

import (
	"time"

	"github.com/genss333/go-clean-architecture/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ZapLogger() gin.HandlerFunc {
	skipPaths := map[string]bool{"/metrics": true}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if skipPaths[path] {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		logger.Log.Info("HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
