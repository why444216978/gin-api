package redis

import (
	"context"
	"gin-api/libraries/jaeger"
	"net/http"

	"github.com/gomodule/redigo/redis"
	opentracing_log "github.com/opentracing/opentracing-go/log"
)

// Do 执行 redis 命令
// NOTE 除非有必要(比如在一个函数内部需要执行多次 redis 操作), 否则请用该函数执行所有的操作, 这样能有效避免忘记释放资源.
func (db *RedisDB) Do(ctx context.Context, header http.Header, commandName string, args ...interface{}) (reply interface{}, err error) {
	if commandName == "PING" {
		return
	}

	sp, _ := jaeger.InjectRedis(ctx, header, commandName, args)
	if sp != nil {
		defer sp.Finish()
	}

	//后defer连接，先释放，统计挥手时间
	conn := db.pool.Get()
	defer conn.Close()

	if err := ctx.Err(); err != nil {
		if err != nil && sp != nil {
			jaeger.SetError(sp, err)
		}
		return nil, err
	}

	reply, err = conn.Do(commandName, args...)
	if err != nil && sp != nil {
		jaeger.SetError(sp, err)
	}

	return
}

// DoLua 执行lua脚本
func (db *RedisDB) DoLua(ctx context.Context, header http.Header, script string, args ...interface{}) (reply interface{}, err error) {
	lua := redis.NewScript(1, script)
	sp, _ := jaeger.InjectRedis(ctx, header, "lua", args)
	if sp != nil {
		defer sp.Finish()
		sp.LogFields(opentracing_log.Object("script", script))
	}

	//后defer连接，先释放，统计挥手时间
	conn := db.pool.Get()
	defer conn.Close()

	if err := ctx.Err(); err != nil {
		if err != nil && sp != nil {
			jaeger.SetError(sp, err)
		}
		return nil, err
	}

	reply, err = lua.Do(db.pool.Get(), args...)
	if err != nil && sp != nil {
		jaeger.SetError(sp, err)
	}

	return
}

func (db *RedisDB) String(ctx context.Context, header http.Header, commandName string, args ...interface{}) (reply string, err error) {
	reply, err = redis.String(db.Do(ctx, header, commandName, args...))
	if err == redis.ErrNil {
		err = nil
	}

	return
}

func (db *RedisDB) Strings(ctx context.Context, header http.Header, commandName string, args ...interface{}) (reply []string, err error) {
	reply, err = redis.Strings(db.Do(ctx, header, commandName, args...))
	if err == redis.ErrNil {
		err = nil
	}

	return
}

func (db *RedisDB) Int(ctx context.Context, header http.Header, commandName string, args ...interface{}) (reply int, err error) {
	reply, err = redis.Int(db.Do(ctx, header, commandName, args...))
	if err == redis.ErrNil {
		err = nil
	}

	return
}
