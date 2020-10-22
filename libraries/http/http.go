package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"gin-frame/libraries/log"
	"gin-frame/libraries/util"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracingLog "github.com/opentracing/opentracing-go/log"
)

func HttpSend(c *gin.Context, method, url, logId string, data map[string]interface{}) map[string]interface{} {
	ctx := c.Request.Context()
	var (
		statement = url
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

	logFormat.LogId = logId

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
		logFormat.XHop = logFormat.XHop.Next()

		logFormat.Module = "databus/http"

		if err != nil {
			log.Errorf(logFormat, "http[%s]:[%s], error: %s", method, statement, err)
			return
		}
		log.Infof(logFormat, "http[%s]:%s, success", method, statement)
	}()

	client := &http.Client{}

	//请求数据
	byteDates, err := json.Marshal(data)
	util.Must(err)
	reader := bytes.NewReader(byteDates)

	//url
	url = url + "?logid=" + logId

	//构建req
	req, err = http.NewRequest(method, url, reader)
	util.Must(err)

	//设置请求header
	req.Header.Add("content-type", "application/json")

	tracer, _ := c.Get("Tracer")
	parentSpanContext, _ := c.Get("ParentSpanContext")

	span := opentracing.StartSpan(
		"httpDo",
		opentracing.ChildOf(parentSpanContext.(opentracing.SpanContext)),
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
		ext.SpanKindRPCClient,
	)
	defer span.Finish()

	injectErr := tracer.(opentracing.Tracer).Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

	if injectErr != nil {
		span.LogFields(opentracingLog.String("inject-error", err.Error()))
	}

	//发送请求
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
	}

	if b != nil {
		res, err := simplejson.NewJson(b)
		util.Must(err)

		ret["data"] = res
	}

	span.SetTag("code", resp.StatusCode)
	span.SetTag("msg", ret["msg"])
	span.SetTag("data", ret["data"])

	return ret
}
