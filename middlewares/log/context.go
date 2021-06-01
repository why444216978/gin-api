package log

import (
	"gin-api/libraries/logging"

	"github.com/gin-gonic/gin"
)

func InitContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		logging.WithLogID(c, logging.ExtractLogID(c))
		logging.WithHTTPFields(c)

		c.Next()
	}
}
