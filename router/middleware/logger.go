package middleware

import (
	"schedrestd/common"
	"schedrestd/common/logger"
	"github.com/gin-gonic/gin"
	"time"
)

func LoggerToFile(logger logger.AipLogger) gin.HandlerFunc {
	return func(c *gin.Context) {

		// For current context, set the logger
		c.Set(common.LoggerName, logger)

		// Request time
		startTime := time.Now()

		// Handle next request
		c.Next()

		// End time
		endTime := time.Now()

		// Spent time
		latencyTime := endTime.Sub(startTime)

		// HTTP Method
		reqMethod := c.Request.Method

		// URI
		reqUri := c.Request.RequestURI

		// Status Code
		statusCode := c.Writer.Status()

		// Client IP
		clientIP := c.ClientIP()

		// Log format for current request
		logger.Infof("[SCHEDEST] %v | %3d | %13v | %15s | %s | %s |",
			endTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
