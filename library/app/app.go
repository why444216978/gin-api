package app

var App struct {
	AppName        string
	AppPort        int
	Pprof          bool
	IsDebug        bool
	ContextTimeout int
	ConnectTimeout int
	WriteTimeout   int
	ReadTimeout    int
}
