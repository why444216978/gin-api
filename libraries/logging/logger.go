package logging

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/sys"
)

const (
	SHORTDATEFORMAT    = "20060102"
	DATEFORMAT         = "2006-01-02"
	DEFAULT_LOG_LEVEL  = DEBUG
	DEFAULT_LOG_PATH   = "./"
	DEFAULT_LOG_FILE   = "log.log"
	CONTEXT_LOG_HEADER = "log_header"
)

var (
	defaultLogger    *Logger
	defaultLogHeader *LogHeader
	initOnce         sync.Once
)

func Init(logConfig *LogConfig) {
	initOnce.Do(func() {
		//设置默认值
		if logConfig.Path == "" {
			logConfig.Path = DEFAULT_LOG_PATH
		}
		if logConfig.File == "" {
			logConfig.File = DEFAULT_LOG_FILE
		}

		if logConfig.RotatingFileHandler == "" {
			logConfig.RotatingFileHandler = TIMED_ROTATING_FILE_HANDLER
		}

		defaultLogger = NewLogger(logConfig)
		defaultLogHeader = &LogHeader{LogId: NewObjectId().Hex()}
		defaultLogHeader.HostIp, _ = sys.GetInternalIP()
	})
}

type Logger struct {
	c *LogConfig

	hdlrs []LogHandler

	logLevel LogLevel
	shutdown bool //when true,no longer receive new msg

	sync.RWMutex
}

/*
* 命名参考Python TimedRotatingFileHandler
* 目前的实现本质上就是个TimedRotatingFileHandler
* https://docs.python.org/2/library/logging.handlers.html#timedrotatingfilehandler
 */
func NewLogger(c *LogConfig) *Logger {
	logger := &Logger{
		c:        c,
		logLevel: LogLevel(c.Level),
		hdlrs:    []LogHandler{FileHandlerAdapterWithConfig(c)},
	}

	if c.Debug {
		//NewStdColorfulHandler is a singleton call
		logger.hdlrs = append(logger.hdlrs, NewStdColorfulHandler())
	}

	for _, hdlr := range logger.hdlrs {
		go hdlr.Run()
	}

	return logger
}

func (l *Logger) Debugf(header *LogHeader, format string, v ...interface{}) {
	l.logf(DEBUG, header, format, v...)
}
func (l *Logger) Debug(header *LogHeader, v ...interface{}) {
	l.log(DEBUG, header, v...)
}
func (l *Logger) Infof(header *LogHeader, format string, v ...interface{}) {
	l.logf(INFO, header, format, v...)
}
func (l *Logger) Info(header *LogHeader, v ...interface{}) {
	l.log(INFO, header, v...)
}
func (l *Logger) Warnf(header *LogHeader, format string, v ...interface{}) {
	l.logf(WARN, header, format, v...)
}
func (l *Logger) Warn(header *LogHeader, v ...interface{}) {
	l.log(WARN, header, v...)
}
func (l *Logger) Errorf(header *LogHeader, format string, v ...interface{}) {
	l.logf(ERROR, header, format, v...)
}
func (l *Logger) Error(header *LogHeader, v ...interface{}) {
	l.log(ERROR, header, v...)
}

//do nothing if the logger is not initialized
func (l *Logger) logf(lvl LogLevel, header *LogHeader, format string, v ...interface{}) {
	if l != nil {
		if l.LogLevel() <= lvl {
			l.output(4, lvl, header, fmt.Sprintf(format, v...))
		}
	}
}

//do nothing if the logger is not initialized
func (l *Logger) log(lvl LogLevel, header *LogHeader, v ...interface{}) {
	if l != nil {
		if l.LogLevel() <= lvl {
			if len(v) == 1 {
				l.output(4, lvl, header, v[0])
			} else {
				l.output(4, lvl, header, fmt.Sprint(v...))
			}
		}
	}
}

func (l *Logger) output(calldepth int, lvl LogLevel, header *LogHeader, s interface{}) {
	pc, file, line, _ := runtime.Caller(calldepth)

	//TODO
	//此处使用了指针,可能有race的问题,concurrent read and write

	record := &Record{
		Timestamp: ts(time.Now()),
		Level:     lvl,
		Msg:       s,
		File:      filepath.Base(file),
		Line:      line,
		Func:      runtime.FuncForPC(pc).Name(),
		LogHeader: *header,
	}
	record.MilliSecond = millts(record.Timestamp)
	record.HumanTime = hts(record.Timestamp)

	for _, h := range l.hdlrs {
		if l.c.AsyncFormatter {
			//异步formatter
			h.Receive(record)
		} else {
			//同步formatter
			h.SyncFormatterReceive(record)
		}

	}
}

func (l *Logger) LogLevel() (lvl LogLevel) {
	l.RLock()
	lvl = l.logLevel
	l.RUnlock()
	return
}

func (l *Logger) Notify(sig os.Signal) {
	switch sig {
	case syscall.SIGUSR1:
		for _, h := range l.hdlrs {
			h.Notify(&LogSignal{
				Action: SignalReopen,
			})
		}
	}
}

func (l *Logger) Shutdown() <-chan bool {
	l.Lock()
	if l.shutdown {
		l.Unlock()
		c := make(chan bool, 1)
		c <- true
		return c
	} else {
		l.shutdown = true
		l.Unlock()
		c := make(chan bool, 1)

		var wg sync.WaitGroup
		for _, h := range l.hdlrs {
			wg.Add(1)
			//Method params are evaluated when the method is invoke, not at the moment literal statement.
			//Here you must pass h as a param to go routine,for h is an interface.
			//Or you'll get the unexpected result: in every goroutine, h pointing to the last item in l.hdlrs
			// after the for iteration and before the method invoke.
			go func(_h LogHandler) {
				<-_h.Shutdown()
				defer wg.Done()
			}(h)
		}

		wg.Wait()
		c <- true
		return c
	}
}

func Debugf(header *LogHeader, format string, v ...interface{}) {
	defaultLogger.Debugf(header, format, v...)
}
func Debug(header *LogHeader, v ...interface{}) {
	defaultLogger.Debug(header, v...)
}
func Infof(header *LogHeader, format string, v ...interface{}) {
	defaultLogger.Infof(header, format, v...)
}
func Info(header *LogHeader, v ...interface{}) {
	defaultLogger.Info(header, v...)
}
func Errorf(header *LogHeader, format string, v ...interface{}) {
	defaultLogger.Errorf(header, format, v...)
}
func Error(header *LogHeader, v ...interface{}) {
	defaultLogger.Error(header, v...)
}
func ErrorCtx(c *gin.Context, v ...interface{}) {
	defaultLogger.Error(GetLogHeader(c), v...)
}
func Warnf(header *LogHeader, format string, v ...interface{}) {
	defaultLogger.Warnf(header, format, v...)
}
func Warn(header *LogHeader, v ...interface{}) {
	defaultLogger.Warn(header, v...)
}
func Shutdown() <-chan bool {
	if defaultLogger != nil {
		return defaultLogger.Shutdown()
	} else {
		c := make(chan bool, 1)
		c <- true
		return c
	}
}

func Notify(sig os.Signal) {
	if defaultLogger != nil {
		defaultLogger.Notify(sig)
	}
}

func GetLogHeader(c *gin.Context) *LogHeader {
	return c.Request.Context().Value(CONTEXT_LOG_HEADER).(*LogHeader)
}

func WriteLogHeader(c *gin.Context, header *LogHeader) {
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), CONTEXT_LOG_HEADER, header))
}
