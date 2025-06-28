package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// Logger returns a Gin middleware handler that logs requests using Logrus.
func Logger() gin.HandlerFunc {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		if !c.IsAborted() {
			logrus.Infof("[%s] \"%s %s %s\" %d %s",
				start.Format(time.RFC3339),
				c.Request.Method,
				path,
				raw,
				c.Writer.Status(),
				time.Since(start),
			)
		}
	}
}
