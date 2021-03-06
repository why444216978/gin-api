package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"gin-api/libraries/jaeger"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
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

// Do 执行 redis 命令
// NOTE 除非有必要(比如在一个函数内部需要执行多次 redis 操作), 否则请用该函数执行所有的操作, 这样能有效避免忘记释放资源.
func (db *RedisDB) Do(c *gin.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	if commandName == "PING" {
		return
	}

	arg := ""
	for _, v := range args {
		arg = fmt.Sprintf("%v ", v)
	}
	sp, _ := jaeger.InjectRedis(c.Request.Context(), c.Request.Header, commandName, arg)
	if sp != nil {
		defer sp.Finish()
	}

	//后defer连接，先释放，统计挥手时间
	conn := db.pool.Get()
	defer conn.Close()

	if err := c.Request.Context().Err(); err != nil {
		if err != nil && sp != nil {
			jaeger.SetError(sp, err)
		}
		return nil, err
	}

	reply, err = conn.Do(commandName, args...)
	if err != nil && sp != nil {
		jaeger.SetError(sp, err)
	}
	_reply, _ := json.Marshal(reply)
	if sp != nil {
		jaeger.SetResponse(sp, string(_reply))
	}

	return
}

func (db *RedisDB) String(c *gin.Context, commandName string, args ...interface{}) (reply string, err error) {
	reply, err = redis.String(db.Do(c, commandName, args...))
	if err == redis.ErrNil {
		err = nil
	}

	return
}

func (db *RedisDB) Strings(c *gin.Context, commandName string, args ...interface{}) (reply []string, err error) {
	reply, err = redis.Strings(db.Do(c, commandName, args...))
	if err == redis.ErrNil {
		err = nil
	}

	return
}

func (db *RedisDB) Int(c *gin.Context, commandName string, args ...interface{}) (reply int, err error) {
	reply, err = redis.Int(db.Do(c, commandName, args...))
	if err == redis.ErrNil {
		err = nil
	}

	return
}
