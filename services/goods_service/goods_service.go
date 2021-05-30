package goods_service

import (
	"gin-api/resource"
	"strconv"

	"github.com/pkg/errors"
	"github.com/why444216978/go-util/conversion"

	"github.com/gin-gonic/gin"

	redigo "github.com/gomodule/redigo/redis"
)

type GoodsInterface interface {
	GetGoodsPrice(ctx *gin.Context, id int) (int, error)

	GetGoodsName(ctx *gin.Context, id int) (string, error)

	GetGoodsInfo(ctx *gin.Context, id int) map[string]interface{}

	BatchGoodsName(ctx *gin.Context, ids []int) []string
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

func (srv *GoodsService) GetGoodsPrice(ctx *gin.Context, id int) (int, error) {
	data, err := redigo.Int(resource.DefaultRedis.Do(ctx, "GET", GOODS_PRICE_KEY+strconv.Itoa(id)))
	if err != nil {
		err = errors.Wrap(err, "redis get goods price error：")
	}
	return data, err
}

func (srv *GoodsService) GetGoodsName(ctx *gin.Context, id int) (string, error) {
	data, err := redigo.String(resource.DefaultRedis.Do(ctx, "GET", GOODS_NAME_KEY+strconv.Itoa(id)))
	if err != nil {
		err = errors.Wrap(err, "redis get goods price error：")
	}
	return data, err
}

func (srv *GoodsService) GetGoodsInfo(ctx *gin.Context, id int) map[string]interface{} {
	data, _ := redigo.String(resource.DefaultRedis.Do(ctx, "GET", GOODS_NAME_KEY+strconv.Itoa(id)))
	ret, _ := conversion.JsonToMap(data)
	return ret
}

func (srv *GoodsService) BatchGoodsName(ctx *gin.Context, ids []int) []string {
	var args []interface{}
	for _, v := range ids {
		args = append(args, GOODS_NAME_KEY+strconv.Itoa(v))
	}

	data, _ := redigo.Strings(resource.DefaultRedis.Do(ctx, "MGET", args...))

	return data
}
