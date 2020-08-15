package log

// StreamHandler
// FileHandler
// RotatingFileHandler

import (
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	//for TIMEOUT there is no more data, we do somthing
	HANDLE_TIMEOUT = 3 * time.Second //handle timeout
)

type LogHandler interface {
	Run()
	Receive(format *LogFormat)
	SyncFormatterReceive(format *LogFormat)
	Notify(*LogSignal)
	Shutdown() <-chan bool
}

var (
	stdColorfulHandler     *StreamHandler
	stdColorfulHandlerOnce sync.Once
	counter                int64 //ensure solo seq for every record in SYNC mode
)

type StreamHandler struct {
	formatter  Formatter
	logChan    chan *LogFormat //消息处理队列
	logBufChan chan []byte     //buf消息处理队列
	sigChan    chan *LogSignal //信号处理队列
	sigCb      map[int][]func(*LogSignal) error
	w          io.Writer
	n          int64

	shutdown  bool
	shutdownC chan bool

	runonce      sync.Once
	sync.RWMutex       // ensures atomic writes; protects the following fields
	counter      int64 //ensure solo seq for every record in ASYNC mode
}

func NewStreamHandler() *StreamHandler {
	s := &StreamHandler{
		logChan:    make(chan *LogFormat, 10*1024),
		logBufChan: make(chan []byte, 10*1024),
		sigChan:    make(chan *LogSignal, 32),
		shutdownC:  make(chan bool, 1),
		sigCb:      make(map[int][]func(*LogSignal) error),
	}
	s.setSigCb(SignalShutdown, s.shutdownCb)

	return s
}

func NewStdColorfulHandler() *StreamHandler {
	stdColorfulHandlerOnce.Do(func() {
		stdColorfulHandler = NewStreamHandler()
		stdColorfulHandler.setFormatter(NewSimpleFormatter("", Ldate|Lmicroseconds|Lshortfile|Lcolor))
		stdColorfulHandler.setWriter(os.Stdout)
	})

	return stdColorfulHandler
}

func (sh *StreamHandler) shutdownCb(signal *LogSignal) error {
	if signal.Action == SignalShutdown {
		sh.Lock()
		sh.shutdown = true
		sh.Unlock()
	}
	return nil
}

func (sh *StreamHandler) Shutdown() <-chan bool {
	signal := &LogSignal{
		Action: SignalShutdown,
	}
	sh.sigChan <- signal

	return sh.shutdownC
}

func (sh *StreamHandler) Notify(sig *LogSignal) {
	sh.sigChan <- sig
}

func (sh *StreamHandler) setSigCb(Action int, cb func(*LogSignal) error) {
	sh.sigCb[Action] = append(sh.sigCb[Action], cb)
}

func (sh *StreamHandler) setupFormatter(mode int) {
	if mode == MOD_JSON {
		sh.setFormatter(DefaultJSONFormatter)
	} else {
		sh.setFormatter(DefaultFormatter)
	}
}

func (sh *StreamHandler) setFormatter(formatter Formatter) {
	sh.formatter = formatter
}

func (sh *StreamHandler) setWriter(w io.Writer) {
	sh.w = w
}

//Async formatter Receive receive msg for logChan
func (sh *StreamHandler) Receive(r *LogFormat) {
	sh.RLock()
	if sh.shutdown {
		sh.RUnlock()
		log.Println("ERROR StreamHandler has already shutdown, ignore any msg received")
		return
	}
	sh.RUnlock()

	select {
	case sh.logChan <- r:
	default:
		buf, err := DefaultFormatter.Format(r)
		log.Printf("ERROR logChan is full, abandon msg:%s, err: %s\n", buf, err)
	}
}

//sync formatter Receive receive msg for logChan
func (sh *StreamHandler) SyncFormatterReceive(r *LogFormat) {
	sh.RLock()
	if sh.shutdown {
		sh.RUnlock()
		log.Println("ERROR StreamHandler has already shutdown, ignore any msg received")
		return
	}
	sh.RUnlock()

	r.Seq = atomic.AddInt64(&counter, 1)

	var (
		buf []byte
		err error
	)

	if sh.formatter == nil {
		buf, err = DefaultFormatter.Format(r)
	} else {
		buf, err = sh.formatter.Format(r)
	}

	if err != nil {
		log.Printf("ERROR: formatter error:%s\n", err)
		return
	}

	select {
	case sh.logBufChan <- buf:
	default:
		log.Printf("ERROR logBufChan is full, abandon msg:%s", buf)
	}
}

//Signal receive signal for sigChan
func (sh *StreamHandler) Signal(signal *LogSignal) {
	sh.RLock()
	if sh.shutdown {
		sh.RUnlock()
		log.Println("ERROR StreamHandler has already shutdown, ignore any msg received")
		return
	}
	sh.RUnlock()

	select {
	case sh.sigChan <- signal:
	default:
		log.Printf("ERROR sigChan is full, abandon signal Action:%d, Payload:%v\n", signal.Action, signal.Payload)
	}
}

// runonce
// 单go routine处理消息和信号,不用锁
// handle之前必须完成所有的准备工作，否则可能有race风险
func (sh *StreamHandler) Run() {
	sh.runonce.Do(func() {
		var (
			timer  = time.NewTimer(HANDLE_TIMEOUT)
			record *LogFormat
		)

		for {
			select {
			case record = <-sh.logChan:
				sh.counter += 1
				record.Seq = sh.counter //只有当前goroutine读写counter,不会RACE
				var (
					buf []byte
					err error
					n   int
				)

				if sh.formatter == nil {
					buf, err = DefaultFormatter.Format(record)
				} else {
					buf, err = sh.formatter.Format(record)
				}

				if err != nil {
					log.Printf("ERROR: formatter error:%s\n", err)
				} else {
					if n, err = sh.w.Write(buf); err != nil {
						log.Printf("ERROR: io.Writer err:%s\n", err)
					} else {
						sh.n += int64(n)

						if n, err = sh.w.Write([]byte("\n")); err != nil {
							log.Printf("ERROR: io.Writer err:%s\n", err)
						} else {
							sh.n += int64(n)
						}
					}
				}
			case buf := <-sh.logBufChan:
				var (
					err error
					n   int
				)

				if n, err = sh.w.Write(buf); err != nil {
					log.Printf("ERROR: io.Writer err:%s\n", err)
				} else {
					sh.n += int64(n)

					if n, err = sh.w.Write([]byte("\n")); err != nil {
						log.Printf("ERROR: io.Writer err:%s\n", err)
					} else {
						sh.n += int64(n)
					}
				}
			case signal := <-sh.sigChan:
				for _, cb := range sh.sigCb[signal.Action] {
					if err := cb(signal); err != nil {
						log.Printf("signal callback ERROR: %s\n", err)
					}
				}
			case <-timer.C:
				//no logChan msg, no sigChan msg, default here
				sh.RLock()
				if sh.shutdown {
					sh.RUnlock()
					sh.shutdownC <- true
					goto END
				} else {
					sh.RUnlock()
					timer.Reset(HANDLE_TIMEOUT)
				}
			}
		}
	END:
	})
}

func FileHandlerAdapterWithConfig(c *LogConfig) LogHandler {
	if c.Rotate {
		return NewRotatingFileHandler(c)
	} else {
		return NewFileHandler(c)
	}
}

type FileHandler struct {
	LogFile
	*StreamHandler
}

func NewFileHandler(c *LogConfig) *FileHandler {
	fh := &FileHandler{
		LogFile: LogFile{
			path: c.Path,
			file: c.File,
		},
		StreamHandler: NewStreamHandler(),
	}

	//setup formatter
	fh.setupFormatter(c.Mode)

	if err := fh.OpenFile(); err != nil {
		log.Fatal(err)
	} else {
		fh.setWriter(fh.Writer())
	}

	fh.setSigCb(SignalReopen, func(signal *LogSignal) (err error) {
		if signal.Action == SignalReopen {
			return fh.Reopen()
		}
		return nil
	})

	return fh
}

func (fh *FileHandler) Reopen() (err error) {
	//Reopen
	if err = fh.Close(); err != nil {
		return
	}

	if err = fh.OpenFile(); err != nil {
		return
	}

	fh.setWriter(fh.Writer())

	return
}

type RotatingFileHandler struct {
	*FileHandler
	rotator    FileRotator
	rotateChan chan error //used for logging rotate
}

func NewRotatingFileHandler(c *LogConfig) LogHandler {
	rfh := &RotatingFileHandler{
		FileHandler: NewFileHandler(c),
		rotateChan:  make(chan error),
	}

	switch c.RotatingFileHandler {
	case SIZED_ROTATING_FILE_HANDLER:
		log.Fatal("unexpected RotatingFileHandler:", c.RotatingFileHandler)
	case CRON_ROTATING_FILE_HANDLER:
		rfh.rotator = NewCronFileRotator(c)
	case TIMED_ROTATING_FILE_HANDLER:
		rfh.rotator = NewTimedFileRotator(c)
	default:
		log.Fatal("unexpected RotatingFileHandler:", c.RotatingFileHandler)
	}

	rfh.setSigCb(SignalRotate, rfh.RotateAndReopen)

	go func() {
		for {
			<-rfh.rotator.Monitor()
			rfh.sigChan <- &LogSignal{
				Action:  SignalRotate,
				Payload: rfh.rotator.NextFilePath(),
			}
			<-rfh.rotateChan

			rfh.rotator.NextRotate()
		}
	}()

	return rfh
}

func (rfh *RotatingFileHandler) RotateAndReopen(signal *LogSignal) (err error) {
	rfh.rotator.Rotate()
	rfh.rotateChan <- nil

	return rfh.Reopen()
}
