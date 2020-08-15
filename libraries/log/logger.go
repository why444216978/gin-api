package log

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
)

const (
	DATETIME          = "2006-01-02 15:04:05"
	SHORTDATEFORMAT   = "20060102"
	DATEFORMAT        = "2006-01-02"
	DEFAULT_LOG_LEVEL = DEBUG
	DEFAULT_LOG_PATH  = "./"
	DEFAULT_LOG_FILE  = "default.log"
)

var (
	runLogger   *Logger
	logInitOnce sync.Once

	errorLogger   *Logger
	errorInitOnce sync.Once
)

func InitError(logConfig *LogConfig, path, file string) *Logger {
	errorInitOnce.Do(func() {
		//设置默认值
		if logConfig.Path == "" {
			logConfig.Path = DEFAULT_LOG_PATH
		}
		if logConfig.File == "" {
			logConfig.File = DEFAULT_LOG_FILE
		}
		if path != "" {
			logConfig.Path = path
		}
		if file != "" {
			logConfig.File = file
		}

		if logConfig.RotatingFileHandler == "" {
			logConfig.RotatingFileHandler = TIMED_ROTATING_FILE_HANDLER
		}

		errorLogger = NewLogger(logConfig)
	})

	if path != logConfig.Path {
		logConfig.Path = path
	}

	if file != logConfig.File {
		logConfig.File = file
	}

	return errorLogger
}

func InitRun(logConfig *LogConfig, path, file string) *Logger {
	logInitOnce.Do(func() {
		//设置默认值
		if logConfig.Path == "" {
			logConfig.Path = DEFAULT_LOG_PATH
		}
		if logConfig.File == "" {
			logConfig.File = DEFAULT_LOG_FILE
		}
		if path != "" {
			logConfig.Path = path
		}
		if file != "" {
			logConfig.File = file
		}

		if logConfig.RotatingFileHandler == "" {
			logConfig.RotatingFileHandler = TIMED_ROTATING_FILE_HANDLER
		}

		runLogger = NewLogger(logConfig)
	})

	if path != logConfig.Path {
		logConfig.Path = path
	}

	if file != logConfig.File {
		logConfig.File = file
	}

	return runLogger
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

func (l *Logger) Debugf(header *LogFormat, format string, v ...interface{}) {
	l.logf(DEBUG, header, format, v...)
}
func (l *Logger) Debug(header *LogFormat, v ...interface{}) {
	l.log(DEBUG, header, v...)
}
func (l *Logger) Infof(header *LogFormat, format string, v ...interface{}) {
	l.logf(INFO, header, format, v...)
}
func (l *Logger) Info(header *LogFormat, v ...interface{}) {
	l.log(INFO, header, v...)
}
func (l *Logger) Warnf(header *LogFormat, format string, v ...interface{}) {
	l.logf(WARN, header, format, v...)
}
func (l *Logger) Warn(header *LogFormat, v ...interface{}) {
	l.log(WARN, header, v...)
}
func (l *Logger) Errorf(header *LogFormat, format string, v ...interface{}) {
	l.logf(ERROR, header, format, v...)
}
func (l *Logger) Error(header *LogFormat, v ...interface{}) {
	l.log(ERROR, header, v...)
}

//do nothing if the logger is not initialized
func (l *Logger) logf(lvl LogLevel, header *LogFormat, format string, v ...interface{}) {
	if l != nil {
		if l.LogLevel() <= lvl {
			l.output(4, lvl, header, fmt.Sprintf(format, v...))
		}
	}
}

//do nothing if the logger is not initialized
func (l *Logger) log(lvl LogLevel, header *LogFormat, v ...interface{}) {
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

func (l *Logger) output(calldepth int, lvl LogLevel, logContent *LogFormat, s interface{}) {
	//pc, file, line, _ := runtime.Caller(calldepth)

	//TODO
	//此处使用了指针,可能有race的问题,concurrent read and write

	logContent.Level = lvl
	logContent.HumanTime = logContent.StartTime.Format("2006-01-02 15:04:05")
	logContent.Msg = s
	logContent.EndTime = time.Now()
	//logContent.MilliSecond = millts(logContent.Timestamp)
	latencyTime := logContent.EndTime.Sub(logContent.StartTime).Microseconds() // 执行时间
	logContent.LatencyTime = latencyTime
	logContent.TimeUnit = "Microseconds"

	for _, h := range l.hdlrs {
		if l.c.AsyncFormatter {
			//异步formatter
			h.Receive(logContent)
		} else {
			//同步formatter
			h.SyncFormatterReceive(logContent)
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

func Debugf(header *LogFormat, format string, v ...interface{}) {
	runLogger.Debugf(header, format, v...)
}
func Debug(header *LogFormat, v ...interface{}) {
	runLogger.Debug(header, v...)
}
func Infof(header *LogFormat, format string, v ...interface{}) {
	runLogger.Infof(header, format, v...)
}
func Info(header *LogFormat, v ...interface{}) {
	runLogger.Info(header, v...)
}
func Errorf(header *LogFormat, format string, v ...interface{}) {
	errorLogger.Errorf(header, format, v...)
}
func Error(header *LogFormat, v ...interface{}) {
	errorLogger.Error(header, v...)
}
func Warnf(header *LogFormat, format string, v ...interface{}) {
	runLogger.Warnf(header, format, v...)
}
func Warn(header *LogFormat, v ...interface{}) {
	runLogger.Warn(header, v...)
}
func Shutdown() <-chan bool {
	if runLogger != nil {
		return runLogger.Shutdown()
	} else {
		c := make(chan bool, 1)
		c <- true
		return c
	}
}

func Notify(sig os.Signal) {
	if runLogger != nil {
		runLogger.Notify(sig)
	}
}
