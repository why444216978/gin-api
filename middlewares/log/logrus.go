package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-api/app_const"
	"gin-api/libraries/config"
	"gin-api/libraries/logging"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/why444216978/go-util/conversion"
	"github.com/why444216978/go-util/sys"
	util_time "github.com/why444216978/go-util/time"
)

func Logger() gin.HandlerFunc {

	logFilePath := "./logs"
	logFileName := "api.log"

	// 日志文件
	fileName := path.Join(logFilePath, logFileName)

	// 写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	// 实例化
	logger := logrus.New()

	// 设置输出
	logger.Out = src

	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		fileName+"."+util_time.Date("YmdHi", time.Now()),
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(1*time.Minute),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 新增 Hook
	logger.AddHook(lfHook)

	logCfg := config.GetConfigToJson("log", "log")
	queryLogField := logCfg["query_field"].(string)
	headerLogField := logCfg["header_field"].(string)
	return func(c *gin.Context) {
		var logID string
		switch {
		case c.Query(queryLogField) != "":
			logID = c.Query(queryLogField)
		case c.Request.Header.Get(headerLogField) != "":
			logID = c.Request.Header.Get(headerLogField)
		default:
			logID = logging.NewObjectId().Hex()
		}
		c.Header(headerLogField, logID)

		reqBody := []byte{}
		if c.Request.Body != nil { // Read
			reqBody, _ = ioutil.ReadAll(c.Request.Body)
		}
		reqBodyMap, _ := conversion.JsonToMap(string(reqBody))
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

		hostIP, _ := sys.ExternalIP()

		start := time.Now()

		responseWriter := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		c.Next() // 处理请求

		responseBody := responseWriter.body.String()
		responseBodyMap, _ := conversion.JsonToMap(responseBody)

		common := &logging.Common{
			LogID: logID,
		}
		logging.WriteLogCommon(c, common)

		fields := logging.Fields{
			Header:   c.Request.Header,
			Method:   c.Request.Method,
			Request:  reqBodyMap,
			Response: responseBodyMap,
			Code:     c.Writer.Status(),
			CallerIP: c.ClientIP(),
			HostIP:   hostIP,
			Port:     app_const.SERVICE_PORT,
			API:      c.Request.RequestURI,
			Module:   "HTTP",
			Cost:     int64(time.Now().Sub(start)),
		}
		fields.Common = *common

		data, _ := conversion.StructToJson(common)

		var logFields logrus.Fields
		json.Unmarshal([]byte(data), &logFields)

		logger.WithFields(logFields).Info()
	}
}
