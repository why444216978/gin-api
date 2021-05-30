package logging

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	CONTEXT_LOG_KEY = "log"
)

func GetLogCommon(c *gin.Context) (comm *Common) {
	h := c.Request.Context().Value(CONTEXT_LOG_KEY)
	comm, ok := h.(*Common)
	if !ok {
		comm = &Common{}
	}
	return
}

func WriteLogCommon(c *gin.Context, comm *Common) {
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), CONTEXT_LOG_KEY, comm))
}
