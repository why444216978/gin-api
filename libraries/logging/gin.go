package logging

import (
	"bytes"
	"context"
	"gin-api/app_const"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/conversion"
	"github.com/why444216978/go-util/sys"
)

type contextKey uint64

const (
	contextLogID contextKey = iota
	contextHTTPLogFields
)

// WithLogID inject log id to context
func WithLogID(c *gin.Context, val interface{}) {
	ctx := context.WithValue(c.Request.Context(), contextLogID, val)
	c.Request = c.Request.WithContext(ctx)
}

// ValueLogID extrect log id from context
func ValueLogID(c *gin.Context) string {
	val := c.Request.Context().Value(contextLogID)
	logID, ok := val.(string)
	if !ok {
		return ""
	}
	return logID
}

// WithHTTPFields inject common http log fields to context
func WithHTTPFields(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), contextLogID, InitHTTPFields(c))
	c.Request = c.Request.WithContext(ctx)
}

// ValueHTTPFields extrect common http log fields from context
func ValueHTTPFields(c *gin.Context) Fields {
	val := c.Request.Context().Value(contextLogID)
	fields, ok := val.(Fields)
	if !ok {
		return Fields{}
	}
	return fields
}

// GetRequestBody get http request body
func GetRequestBody(c *gin.Context) map[string]interface{} {
	reqBody := []byte{}
	if c.Request.Body != nil { // Read
		reqBody, _ = ioutil.ReadAll(c.Request.Body)
	}
	reqBodyMap, _ := conversion.JsonToMap(string(reqBody))
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

	return reqBodyMap
}

// InitHTTPFields init http fields
func InitHTTPFields(c *gin.Context) Fields {
	hostIP, _ := sys.ExternalIP()
	return Fields{
		LogID:    GetLogID(c),
		Header:   c.Request.Header,
		Method:   c.Request.Method,
		Request:  GetRequestBody(c),
		CallerIP: c.ClientIP(),
		HostIP:   hostIP,
		Port:     app_const.SERVICE_PORT,
		API:      c.Request.RequestURI,
		Module:   MODULE_HTTP,
	}
}

// MergeHTTPFields merge extend log fields and  http common fields
func MergeHTTPFields(c *gin.Context, extend map[string]interface{}) map[string]interface{} {
	fields := InitHTTPFields(c)
	common, _ := conversion.StructToMap(fields)

	ret := make(map[string]interface{})
	for k, v := range common {
		ret[k] = v
	}
	for k, v := range extend {
		ret[k] = v
	}

	return ret
}
