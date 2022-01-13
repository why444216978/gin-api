package logger

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/why444216978/go-util/conversion"
	"github.com/why444216978/go-util/snowflake"

	"github.com/why444216978/gin-api/library/logger"
)

// ExtractLogID init log id
func ExtractLogID(req *http.Request) string {
	logID := req.Header.Get(logger.LogHeader)

	if logID == "" {
		logID = snowflake.Generate().String()
	}

	req.Header.Add(logger.LogHeader, logID)

	return logID
}

// GetRequestBody get http request body
func GetRequestBody(req *http.Request) map[string]interface{} {
	reqBody := []byte{}
	if req.Body != nil { // Read
		reqBody, _ = ioutil.ReadAll(req.Body)
	}
	reqBodyMap, _ := conversion.JsonToMap(string(reqBody))
	req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

	return reqBodyMap
}
