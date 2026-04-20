package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		log.Printf("[%s] %s %s | %d | %v | %s",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			c.Writer.Status(),
			time.Since(start),
			c.ClientIP(),
		)
	}
}
