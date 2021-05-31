package logging

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/conversion"
)

const (
	LOG_FIELD = "Log-Id"
)

type Common struct {
	LogID string `json:"log_id"`
}

type Fields struct {
	Common
	Header   http.Header `json:"header"`
	Method   string      `json:"method"`
	Request  interface{} `json:"request"`
	Response interface{} `json:"response"`
	Code     int         `json:"code"`
	CallerIP string      `json:"caller_ip"`
	HostIP   string      `json:"host_ip"`
	Port     int         `json:"port"`
	API      string      `json:"api"`
	TraceID  string      `json:"trace_id"`
	SpanID   string      `json:"span_id"`
	Cost     int64       `json:"cost"`
	Module   string      `json:"module"`
	Trace    interface{} `json:"trace"`
}

func GetRequestBody(c *gin.Context) map[string]interface{} {
	reqBody := []byte{}
	if c.Request.Body != nil { // Read
		reqBody, _ = ioutil.ReadAll(c.Request.Body)
	}
	reqBodyMap, _ := conversion.JsonToMap(string(reqBody))
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

	return reqBodyMap
}
