package goods_service

import (
	"strconv"

	"github.com/why444216978/gin-api/resource"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

type GoodsInterface interface {
	GetGoodsName(c *gin.Context, id int) (string, error)
}

var Instance GoodsInterface

type GoodsService struct{}

func init() {
	Instance = &GoodsService{}
}

const (
	GOODS_NAME_KEY  = "goods::name::"
	GOODS_PRICE_KEY = "goods::price::"
)

func (srv *GoodsService) GetGoodsName(c *gin.Context, id int) (string, error) {
	data, err := resource.RedisCache.Get(c.Request.Context(), GOODS_NAME_KEY+strconv.Itoa(id)).Result()
	if err == redis.Nil {
		err = nil
	}
	if err != nil {
		err = errors.Wrap(err, "redis get goods price error：")
	}
	return data, err
}
