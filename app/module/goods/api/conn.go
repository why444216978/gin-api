package api

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/errgroup"

	"github.com/why444216978/gin-api/app/module/goods/respository"
	goodsService "github.com/why444216978/gin-api/app/module/goods/service"
	"github.com/why444216978/gin-api/app/resource"
	"github.com/why444216978/gin-api/app/response"
)

func Do(c *gin.Context) {
	var (
		err   error
		goods respository.Test
	)

	ctx := c.Request.Context()

	defer func() {
		if err != nil {
			resource.ServiceLogger.Error(ctx, err.Error())
			response.Response(c, response.CodeServer, goods, err.Error())
			return
		}
		response.Response(c, response.CodeSuccess, goods, "")
	}()

	goods, err = goodsService.Instance.CrudGoods(ctx)
	if err != nil {
		return
	}

	data := &Data{}

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		goods.Name = "golang"
		_, err = goodsService.Instance.GetGoodsName(ctx, 1)
		return
	})
	g.Go(func() (err error) {
		err = resource.RedisCache.GetData(ctx, "cache_key", time.Hour, time.Hour, GetDataA, data)
		return
	})
	g.Go(func() (err error) {
		result := make([]*redis.StringCmd, 0)
		pipe := resource.RedisDefault.Pipeline()
		result = append(result, pipe.Get(ctx, "test"))
		result = append(result, pipe.Get(ctx, "test"))
		result = append(result, pipe.Get(ctx, "test"))
		_, _ = pipe.Exec(ctx)
		return
	})
	err = g.Wait()
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
