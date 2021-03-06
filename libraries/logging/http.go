package logging

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/why444216978/go-util/conversion"
)

type contextKey uint64

const (
	contextLogID contextKey = iota
	contextHTTPRequestBodyFields
	contextHTTPResponseBodyFields
	contextHTTPLogFields
	contextTraceID
)

// ExtractLogID init log id
func ExtractLogID(req *http.Request) string {
	logID := req.Header.Get(LOG_FIELD)

	if logID == "" {
		logID = newObjectId().Hex()
	}

	req.Header.Add(LOG_FIELD, logID)

	return logID
}

// WithLogID inject log id to context
func WithLogID(ctx context.Context, val interface{}) context.Context {
	return context.WithValue(ctx, contextLogID, val)
}

// ValueLogID extract log id from context
func ValueLogID(ctx context.Context) string {
	val := ctx.Value(contextLogID)
	logID, ok := val.(string)
	if !ok {
		return ""
	}
	return logID
}

// WithTraceID inject trace_id id to context
func WithTraceID(ctx context.Context, val interface{}) context.Context {
	return context.WithValue(ctx, contextTraceID, val)
}

// ValueTraceID extract trace id from context
func ValueTraceID(ctx context.Context) string {
	val := ctx.Value(contextTraceID)
	logID, ok := val.(string)
	if !ok {
		return ""
	}
	return logID
}

// WithHTTPFields inject common http log fields to context
func WithHTTPFields(ctx context.Context, fields Fields) context.Context {
	return context.WithValue(ctx, contextHTTPLogFields, fields)
}

// ValueHTTPFields extrect common http log fields from context
func ValueHTTPFields(ctx context.Context) Fields {
	val := ctx.Value(contextHTTPLogFields)
	fields, ok := val.(Fields)
	if !ok {
		return Fields{}
	}
	return fields
}

// WithHTTPRequestBody inject common http request body to context
func WithHTTPRequestBody(ctx context.Context, body interface{}) context.Context {
	return context.WithValue(ctx, contextHTTPRequestBodyFields, body)
}

// ValueHTTPRequestBody extrect common http request body from context
func ValueHTTPRequestBody(ctx context.Context) interface{} {
	return ctx.Value(contextHTTPRequestBodyFields)
}

// WithHTTPResponseBody inject common http response body to context
func WithHTTPResponseBody(ctx context.Context, body interface{}) context.Context {
	return context.WithValue(ctx, contextHTTPResponseBodyFields, body)
}

// ValueHTTPResponseBody extrect common http request body from context
func ValueHTTPResponseBody(ctx context.Context) interface{} {
	return ctx.Value(contextHTTPResponseBodyFields)
}

// AddTraceID add trace id to global fields
func AddTraceID(ctx context.Context, traceID string) context.Context {
	fields := ValueHTTPFields(ctx)
	fields.TraceID = traceID
	return WithHTTPFields(ctx, fields)
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

// MergeHTTPFields merge extend log fields and  http common fields
func MergeHTTPFields(ctx context.Context, extend map[string]interface{}) map[string]interface{} {
	fields := ValueHTTPFields(ctx)
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
