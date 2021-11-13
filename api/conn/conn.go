package conn

import (
	"context"
	"errors"
	"time"

	"github.com/why444216978/gin-api/resource"
	"github.com/why444216978/gin-api/response"
	"github.com/why444216978/gin-api/services/goods/respository"
	goods_service "github.com/why444216978/gin-api/services/goods/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
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

	goods, err = goods_service.Instance.CrudGoods(ctx)
	if err != nil {
		return
	}

	data := &Data{}

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		goods.Name = "golang"
		_, err = goods_service.Instance.GetGoodsName(ctx, 1)
		return
	})
	g.Go(func() (err error) {
		err = resource.RedisCache.GetData(ctx, "cache_key", time.Hour, time.Hour, GetDataA, data)
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
