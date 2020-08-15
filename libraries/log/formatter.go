package log

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"time"
)

const (
	// Bits or'ed together to control what's printed.
	// There is no control over the order they appear (the order listed
	// here) or the format they present (as described in the comments).
	// The prefix is followed by a colon only when Llongfile or Lshortfile
	// is specified.
	// For example, flags Ldate | Ltime (or LstdFlags) produce,
	//	2009/01/23 01:23:23 message
	// while flags Ldate | Ltime | Lmicroseconds | Llongfile produce,
	//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lcolor                        // whether colorful output to terminal or not
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

var (
	DEFAULT_COLOR = map[LogLevel]int{
		DEBUG: 34, //BLUE,
		INFO:  32, //GREEN,
		WARN:  33, //YELLOW,
		ERROR: 31, //RED
	}

	DefaultFormatter     = NewSimpleFormatter("", LstdFlags)
	DefaultJSONFormatter = &JSONFormatter{}
)

type Formatter interface {
	Format(format *LogFormat) ([]byte, error)
}

type SimpleFormatter struct {
	flag   int    // properties
	prefix string // prefix to write at beginning of each line
}

func NewSimpleFormatter(prefix string, flag int) *SimpleFormatter {
	return &SimpleFormatter{prefix: prefix, flag: flag}
}

func (f *SimpleFormatter) Format(record *LogFormat) ([]byte, error) {
	var (
		buf []byte = []byte{}
	)

	if (f.flag | Lcolor) != 0 {
		buf = append(buf, []byte(fmt.Sprintf("\033[%dm", DEFAULT_COLOR[record.Level]))...)
	}

	//f.formatHeader(&buf, time.Time(record.MilliSecond), record.File, record.Line)
	buf = append(buf, []byte(fmt.Sprintf("%s", record.Msg))...)

	if (f.flag | Lcolor) != 0 {
		buf = append(buf, []byte("\033[0m")...)
	}

	return buf, nil
}

// Cheap integer to fixed-width decimal ASCII.  Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (f *SimpleFormatter) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	*buf = append(*buf, f.prefix...)
	if f.flag&LUTC != 0 {
		t = t.UTC()
	}
	if f.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if f.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if f.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if f.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if f.flag&(Lshortfile|Llongfile) != 0 {
		if f.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

type JSONFormatter struct {
}

// the encoding/json package will do HTML-escape automatically
// for example, "a=b&c=d" will be encoded as "a=b\u0026c=d"
//
// https://github.com/golang/go/issues/8592
//
func (f *JSONFormatter) Format(record *LogFormat) (buf []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%s\n", debug.Stack())
			log.Printf("JSONFormatter run time panic: %v\n", r)

			//set error
			if re, ok := r.(runtime.Error); ok {
				err = re
			} else if s, ok := r.(string); ok {
				err = fmt.Errorf("%s", s)
			} else {
				err = r.(error)
			}
		}
	}()

	return json.Marshal(record)
}
