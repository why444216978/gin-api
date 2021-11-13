package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	ServiceName    string
	Host           string
	Port           int
	Auth           string
	DB             int
	ConnectTimeout int
	ReadTimeout    int
	WriteTimeout   int
	MaxActive      int
	MaxIdle        int
	IsLog          bool
	ExecTimeout    int64
}

func NewClient(cfg *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	})
}
