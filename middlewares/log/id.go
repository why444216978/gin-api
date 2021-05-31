package log

import (
	"gin-api/libraries/logging"

	"github.com/gin-gonic/gin"
)

func WithLogID() gin.HandlerFunc {
	return func(c *gin.Context) {
		common := &logging.Common{
			LogID: logging.GetLogID(c),
		}
		logging.WriteLogCommon(c, common)
		c.Next()
	}
}
