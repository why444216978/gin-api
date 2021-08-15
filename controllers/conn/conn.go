package conn

import (
	"context"
	"errors"
	"gin-api/libraries/cache"
	"gin-api/libraries/lock"
	"gin-api/libraries/logging"
	"gin-api/resource"
	"gin-api/response"
	"gin-api/services/goods_service"
	"gin-api/services/test_service"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
)

func Do(c *gin.Context) {
	goods, _ := test_service.New().GetFirstRow(c, true)
	g, _ := errgroup.WithContext(c.Request.Context())
	g.Go(func() (err error) {
		goods.Name = "golang"
		_, err = goods_service.Instance.GetGoodsName(c, 1)
		return
	})
	err := g.Wait()
	if err != nil {
		response.Response(c, response.CodeServer, goods, "")
		return
	}

	resource.Logger.Debug("test conn error msg", logging.MergeHTTPFields(c.Request.Context(), map[string]interface{}{"err": "test err"}))

	data := &Data{}

	lock, err := lock.New(resource.RedisCache)
	if err != nil {
		response.Response(c, response.CodeServer, goods, err.Error())
		return
	}

	cache, err := cache.New(resource.RedisCache, lock)
	if err != nil {
		response.Response(c, response.CodeServer, goods, err.Error())
		return
	}

	err = cache.GetData(c.Request.Context(), "key", time.Hour, 60, GetDataA, data)
	if err != nil {
		response.Response(c, response.CodeServer, goods, err.Error())
		return
	}

	response.Response(c, response.CodeSuccess, goods, "")
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
