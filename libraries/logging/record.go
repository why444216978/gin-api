package logging

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	Timestamp   ts          `json:"timestamp"`
	MilliSecond millts      `json:"millisecond"`
	HumanTime   hts         `json:"human_time"`
	Level       LogLevel    `json:"level"`
	File        string      `json:"file"`
	Line        int         `json:"line"`
	Func        string      `json:"func"`
	Msg         interface{} `json:"msg"`
	Trace       interface{} `json:"trace,omitempty"`
	Seq         int64       `json:"seq"` //用于日志排序
	LogHeader
	RPCRecord `json:",omitempty"`
}

type LogHeader struct {
	LogId                 string `json:"logid"`
	CallerIp              string `json:"caller_ip"`
	HostIp                string `json:"host_ip"`
	Port                  int
	Product               string           `json:"product"`
	Module                string           `json:"module"`
	ServiceId             string           `json:"service_id"`
	InstanceId            string           `json:"instance_id"`
	UriPath               string           `json:"uri_path"`
	Tag                   string           `json:"tag"`
	Env                   string           `json:"env"`
	SVersion              string           `json:"sversion"`
	Request               *http.Request    `json:"-"`
	DownstreamErrServices string           `json:"downstream_err_services,omitempty"`
	TraceID               string           `json:"trace_id"`
	SpanID                string           `json:"span_id"`
	GrayMark              string           `json:"gray_mark"`
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
		LogId:      h.LogId,
		CallerIp:   h.CallerIp,
		HostIp:     h.HostIp,
		Product:    h.Product,
		Module:     h.Module,
		ServiceId:  h.ServiceId,
		InstanceId: h.InstanceId,
		UriPath:    h.UriPath,
		Tag:        h.Tag,
		Env:        h.Env,
		SVersion:   h.SVersion,
		Request:    h.Request,
		GrayMark:   h.GrayMark,
	}
}

func (h *LogHeader) AddTag(tag ...string) {
	if h == nil {
		return
	}

	var (
		ss  []string
		set map[string]bool
	)
	if h.Tag != "" {
		ss = strings.Split(h.Tag, ",")
	}
	ss = append(ss, tag...)
	set = make(map[string]bool, len(ss))
	//去重
	for _, s := range ss {
		if s != "" {
			set[s] = true
		}
	}

	if len(set) == 0 {
		h.Tag = ""
	} else {
		ss = make([]string, len(set))
		idx := 0
		for s, _ := range set {
			ss[idx] = s
			idx += 1
		}

		h.Tag = strings.Join(ss, ",")
	}
}

func (h *LogHeader) SetTag(tag ...string) {
	if h == nil {
		return
	}

	h.Tag = ""
	for _, t := range tag {
		h.AddTag(t)
	}
}

func (h *LogHeader) GetAppKey() string {
	if h == nil {
		return ""
	} else {
		return fmt.Sprintf("%s_%s_%s", h.Product, h.Module, h.Env)
	}
}

type RPCRecord struct {
	StatusCode int    `json:"status_code,omitempty"`
	RequestUrl string `json:"request_url,omitempty"`
}
