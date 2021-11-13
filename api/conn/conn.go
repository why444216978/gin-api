package conn

import (
	"context"
	"errors"
	"time"

	redis_cache "github.com/why444216978/gin-api/library/cache/redis"
	redis_lock "github.com/why444216978/gin-api/library/lock/redis"
	"github.com/why444216978/gin-api/resource"
	"github.com/why444216978/gin-api/response"
	goods_service "github.com/why444216978/gin-api/services/goods/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func Do(c *gin.Context) {
	ctx := c.Request.Context()

	goods, err := goods_service.Instance.CrudGoods(ctx)
	if err != nil {
		response.Response(c, response.CodeServer, goods, err.Error())
		return
	}

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		goods.Name = "golang"
		_, err = goods_service.Instance.GetGoodsName(ctx, 1)
		return
	})
	err = g.Wait()
	if err != nil {
		return
	}

	data := &Data{}

	lock, err := redis_lock.New(resource.RedisCache)
	if err != nil {
		return
	}

	cache, err := redis_cache.New(resource.RedisCache, lock)
	if err != nil {
		return
	}

	err = cache.GetData(ctx, "key", time.Hour, time.Hour, GetDataA, data)
	if err != nil {
		return
	}

}

type Data struct {
	A string `json:"a"`
}

func GetDataA(ctx context.Context, _data interface{}) (err error) {
	data, ok := _data.(*Data)
	if !ok {
		err = errors.New("err assert")
		return
	}
	data.A = "a"
	return
}
