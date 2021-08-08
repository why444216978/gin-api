package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	ResultOK = "OK"
)

type Config struct {
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

type RedisDB struct {
	pool *redis.Pool
}

func GetRedis(cfg Config) (db *RedisDB, err error) {
	db = new(RedisDB)
	db.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
				redis.DialPassword(cfg.Auth),
				redis.DialDatabase(cfg.DB),
				redis.DialConnectTimeout(time.Second*time.Duration(cfg.ConnectTimeout)),
				redis.DialReadTimeout(time.Second*time.Duration(cfg.ReadTimeout)),
				redis.DialWriteTimeout(time.Second*time.Duration(cfg.WriteTimeout)),
			)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:     cfg.MaxIdle,   // 最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
		MaxActive:   cfg.MaxActive, // 最大的激活连接数，表示同时最多有N个连接 ，为0事表示没有限制
		IdleTimeout: time.Second,   //最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		Wait:        true,          // 当链接数达到最大后是否阻塞，如果不的话，达到最大后返回错误
	}

	return
}

// ConnPool 返回 redis.Pool.
// 除非必要一般不建议用这个函数, 用本库封装好的函数操作数据库.
func (db *RedisDB) ConnPool() *redis.Pool {
	return db.pool
}

// Close 释放连接资源.
func (db *RedisDB) Close() error {
	if db.pool != nil {
		return db.pool.Close()
	}
	return nil
}
