package log

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"gin-api/global"
// 	"gin-api/libraries/logging"
// 	"os"
// 	"path"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
// 	"github.com/rifflock/lfshook"
// 	"github.com/sirupsen/logrus"
// 	"github.com/why444216978/go-util/conversion"
// 	"github.com/why444216978/go-util/sys"
// 	util_time "github.com/why444216978/go-util/time"
// )

// func Logger() gin.HandlerFunc {

// 	logFilePath := "./logs"
// 	logFileName := "gin-api.log"

// 	// 日志文件
// 	fileName := path.Join(logFilePath, logFileName)

// 	// 写入文件
// 	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
// 	if err != nil {
// 		fmt.Println("err", err)
// 	}

// 	// 实例化
// 	logger := logrus.New()

// 	// 设置输出
// 	logger.Out = src

// 	// 设置日志级别
// 	logger.SetLevel(logrus.DebugLevel)

// 	// 设置 rotatelogs
// 	logWriter, err := rotatelogs.New(
// 		fileName+"."+util_time.Date("YmdHi", time.Now()),
// 		rotatelogs.WithLinkName(fileName),
// 		rotatelogs.WithMaxAge(7*24*time.Hour),
// 		rotatelogs.WithRotationTime(1*time.Minute),
// 	)

// 	writeMap := lfshook.WriterMap{
// 		logrus.InfoLevel:  logWriter,
// 		logrus.FatalLevel: logWriter,
// 		logrus.DebugLevel: logWriter,
// 		logrus.WarnLevel:  logWriter,
// 		logrus.ErrorLevel: logWriter,
// 		logrus.PanicLevel: logWriter,
// 	}

// 	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
// 		TimestampFormat: "2006-01-02 15:04:05",
// 	})

// 	// 新增 Hook
// 	logger.AddHook(lfHook)

// 	return func(c *gin.Context) {
// 		start := time.Now()

// 		responseWriter := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
// 		c.Writer = responseWriter

// 		c.Next()

// 		resp := responseWriter.body.String()
// 		respMap, _ := conversion.JsonToMap(resp)

// 		hostIP, _ := sys.ExternalIP()

// 		fields := logging.Fields{
// 			LogID:    logging.ValueLogID(c),
// 			Header:   c.Request.Header,
// 			Method:   c.Request.Method,
// 			Request:  logging.GetRequestBody(c),
// 			Response: respMap,
// 			Code:     c.Writer.Status(),
// 			CallerIP: c.ClientIP(),
// 			HostIP:   hostIP,
// 			Port:     global.Global.AppPort,
// 			API:      c.Request.RequestURI,
// 			Module:   logging.ModuleHTTP,
// 			Cost:     int64(time.Now().Sub(start)),
// 		}

// 		data, _ := conversion.StructToJson(fields)

// 		var logFields logrus.Fields
// 		json.Unmarshal([]byte(data), &logFields)

// 		logger.WithFields(logFields).Info()
// 	}
// }
