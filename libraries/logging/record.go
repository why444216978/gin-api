package logging

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type LogLevel int

func (lvl LogLevel) MarshalJSON() ([]byte, error) {
	switch lvl {
	case DEBUG:
		return []byte("\"DEBUG\""), nil
	case INFO:
		return []byte("\"INFO\""), nil
	case WARN:
		return []byte("\"WARN\""), nil
	case ERROR:
		return []byte("\"ERROR\""), nil
	default:
		return []byte("\"\""), nil
	}
}

type ts time.Time

func (t ts) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

type millts time.Time

func (t millts) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).UnixNano()/1000000, 10)), nil
}

type hts time.Time

func (t hts) MarshalJSON() ([]byte, error) {
	bs := []byte(fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05.000")))
	//注意MarshalJSON要求两端有双引号
	bs[len(bs)-5] = ',' //和Python,Java等语言统一

	return bs, nil
}

type Record struct {
	Timestamp   ts       `json:"timestamp"`
	MilliSecond millts   `json:"millisecond"`
	HumanTime   hts      `json:"human_time"`
	Level       LogLevel `json:"level"`
	File        string   `json:"file"`
	Line        int      `json:"line"`
	Func        string   `json:"func"`
	LogHeader
	RPCRecord `json:",omitempty"`
}

type LogHeader struct {
	Header   http.Header `json:"header"`
	Request  interface{} `json:"request"`
	Response interface{} `json:"response,omitempty"`
	HTTPCode int         `json:"http_code"`
	LogId    string      `json:"logid"`
	CallerIp string      `json:"caller_ip"`
	HostIp   string      `json:"host_ip"`
	Port     int         `json:"port"`
	UriPath  string      `json:"uri_path"`
	TraceID  string      `json:"trace_id"`
	SpanID   string      `json:"span_id"`
	Cost     int64       `json:"cost"`
	Module   string      `json:"module"`
	Trace    interface{} `json:"trace,omitempty"`
	Error    interface{} `json:"error,omitempty"`
}

func NewLogHeader() *LogHeader {
	logHeader := &LogHeader{
		LogId: NewObjectId().Hex(),
	}

	return logHeader
}

func (h *LogHeader) Dup() *LogHeader {
	if h == nil {
		return NewLogHeader()
	}

	return &LogHeader{
		LogId:    h.LogId,
		CallerIp: h.CallerIp,
		HostIp:   h.HostIp,
		Request:  h.Request,
	}
}

type RPCRecord struct {
	StatusCode int    `json:"status_code,omitempty"`
	RequestUrl string `json:"request_url,omitempty"`
}
