package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		logger.Infof("| %3d | %13v | %15s | %-7s | %s",
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
		)
	}
}
