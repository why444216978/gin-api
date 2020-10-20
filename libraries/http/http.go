package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"gin-frame/libraries/config"
	"gin-frame/libraries/log"
	"gin-frame/libraries/util"
	"gin-frame/libraries/xhop"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func HttpSend(ctx *gin.Context, method, url string, data map[string]interface{}) map[string]interface{} {
	var (
		statement = url
		parent    = opentracing.SpanFromContext(ctx)
		span      opentracing.Span
		startAt   = time.Now()
		endAt     time.Time
		logFormat = log.LogHeaderFromContext(ctx)
		err       error
		ret       = make(map[string]interface{})
		req       *http.Request
	)

	if logFormat == nil {
		logFormat = log.NewLog()
	}

	logFormat.LogId = ctx.Writer.Header().Get(config.GetHeaderLogIdField())

	lastModule := logFormat.Module
	lastStartTime := logFormat.StartTime
	lastEndTime := logFormat.EndTime
	defer func() {
		logFormat.Module = lastModule
		lastStartTime = lastStartTime
		lastEndTime = lastEndTime
	}()
	defer func() {
		endAt = time.Now()
		logFormat.StartTime = startAt
		logFormat.EndTime = endAt
		latencyTime := logFormat.EndTime.Sub(logFormat.StartTime).Microseconds() // 执行时间
		logFormat.LatencyTime = latencyTime
		logFormat.XHop = xhop.NextXhop(ctx, config.GetXhopField())

		logFormat.Module = "databus/http"

		if err != nil {
			log.Errorf(logFormat, "http[%s]:[%s], error: %s", method, statement, err)
			return
		}
		log.Infof(logFormat, "http[%s]:%s, success", method, statement)
	}()

	if parent == nil {
		span = opentracing.StartSpan("httpDo")
	} else {
		span = opentracing.StartSpan("httpDo", opentracing.ChildOf(parent.Context()))
	}
	defer span.Finish()

	span.SetTag("http.type", method)
	span.SetTag("http.statement", url)
	span.SetTag("error", err != nil)

	client := &http.Client{}

	byteDates, err := json.Marshal(data)
	util.Must(err)
	reader := bytes.NewReader(byteDates)

	req, err = http.NewRequest(method, url, reader)

	req.Header.Add("content-type", "application/json")

	//trace.InjectTrace(ctx, logFormat, req)

	resp, err := client.Do(req)
	util.Must(err)
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	util.Must(err)

	ret["code"] = resp.StatusCode
	ret["msg"] = "success"
	ret["data"] = make(map[string]interface{})

	if resp.StatusCode != http.StatusOK {
		ret["msg"] = "http code:" + strconv.Itoa(resp.StatusCode)
		return ret
	}

	if b != nil {
		res, err := simplejson.NewJson(b)
		util.Must(err)

		ret["data"] = res
	}

	return ret
}
