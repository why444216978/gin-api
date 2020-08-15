package redis

type Config struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Password    string `json:"password"`
	DB          int    `json:"db"`
	MaxActive   int    `json:"max_active"`
	MaxIdle     int    `json:"max_idle"`
	IsLog       bool
	ExecTimeout int64 `json:"exec_timeout"` //超时打印日志
}
