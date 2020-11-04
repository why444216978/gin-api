package log

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

type LogFormat struct {
	LogId     string    `json:"logid"`
	HttpCode  int       `json:"http_code"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	//MilliSecond millts      `json:"millisecond"`
	HumanTime string   `json:"human_time"`
	Level     LogLevel `json:"level"`
	//File        string      `json:"file"`
	//Line        int         `json:"line"`
	//Func        string      `json:"func"`
	Msg         interface{} `json:"msg"`
	Trace       interface{} `json:"trace,omitempty"`
	Seq         int64       `json:"seq"` //用于日志排序
	LatencyTime int64       `json:"latency_time"`
	TimeUnit    string      `json:"time_unit"`
	Method      string      `json:"method"`
	StatusCode  int         `json:"status_code"`
	CallerIp    string      `json:"caller_ip"`
	HostIp      string      `json:"host_ip"`
	Port        int         `json:"port"`
	Product     string      `json:"product"`
	Module      string      `json:"module"`
	//ServiceId  string        `json:"service_id"`
	//InstanceId string        `json:"instance_id"`
	UriPath string     `json:"uri_path"`
	//XHop    *xhop.XHop `json:"x_hop"`
	//Tag        string        `json:"tag"`
	Env string `json:"env"`
	//SVersion   string        `json:"sversion"`
	//Stag       app.Stag      `json:"stag"`
	Request *gin.Context `json:"-"`
}

func NewLog() *LogFormat {
	logHeader := &LogFormat{
		LogId: NewObjectId().Hex(),
	}
	//logHeader.XHop = xhop.NewXHop()

	return logHeader
}

func (h *LogFormat) Dup() *LogFormat {
	if h == nil {
		return NewLog()
	}

	return &LogFormat{
		LogId:    h.LogId,
		CallerIp: h.CallerIp,
		HostIp:   h.HostIp,
		Product:  h.Product,
		//Module:     h.Module,
		UriPath: h.UriPath,
		//ServiceId:  h.ServiceId,
		//InstanceId: h.InstanceId,
		//XHop: h.XHop.Dup(),
		//Stag:       h.Stag,
		//Tag:        h.Tag,
		Env: h.Env,
		//SVersion:   h.SVersion,
	}
}

func (h *LogFormat) AddTag(tag ...string) {
	return
	/*if h == nil {
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
	}*/
}

func (h *LogFormat) SetTag(tag ...string) {
	return
	/*if h == nil {
		return
	}

	h.Tag = ""
	for _, t := range tag {
		h.AddTag(t)
	}*/
}

func (h *LogFormat) GetAppKey() string {
	return ""
	/*if h == nil {
		return ""
	} else {
		return fmt.Sprintf("%s_%s_%s", h.Product, h.Module, h.Env)
	}*/
}

type RPCRecord struct {
	StatusCode int    `json:"status_code,omitempty"`
	RequestUrl string `json:"request_url,omitempty"`
}
