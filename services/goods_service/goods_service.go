package goods_service

import (
	"gin-api/resource"
	"strconv"

	"github.com/pkg/errors"
	"github.com/why444216978/go-util/conversion"

	"github.com/gin-gonic/gin"
)

type GoodsInterface interface {
	GetGoodsPrice(c *gin.Context, id int) (int, error)

	GetGoodsName(c *gin.Context, id int) (string, error)

	GetGoodsInfo(c *gin.Context, id int) map[string]interface{}

	BatchGoodsName(c *gin.Context, ids []int) (data []string, err error)
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

func (srv *GoodsService) GetGoodsPrice(c *gin.Context, id int) (int, error) {
	data, err := resource.DefaultRedis.Int(c.Request.Context(), c.Request.Header, "GET", GOODS_PRICE_KEY+strconv.Itoa(id))
	if err != nil {
		err = errors.Wrap(err, "redis get goods price error：")
	}
	return data, err
}

func (srv *GoodsService) GetGoodsName(c *gin.Context, id int) (string, error) {
	data, err := resource.DefaultRedis.String(c.Request.Context(), c.Request.Header, "GET", GOODS_NAME_KEY+strconv.Itoa(id))
	if err != nil {
		err = errors.Wrap(err, "redis get goods price error：")
	}
	return data, err
}

func (srv *GoodsService) GetGoodsInfo(c *gin.Context, id int) map[string]interface{} {
	data, _ := resource.DefaultRedis.String(c.Request.Context(), c.Request.Header, "GET", GOODS_NAME_KEY+strconv.Itoa(id))
	ret, _ := conversion.JsonToMap(data)
	return ret
}

func (srv *GoodsService) BatchGoodsName(c *gin.Context, ids []int) (data []string, err error) {
	var args []interface{}
	for _, v := range ids {
		args = append(args, GOODS_NAME_KEY+strconv.Itoa(v))
	}

	data, err = resource.DefaultRedis.Strings(c.Request.Context(), c.Request.Header, "MGET", args...)
	err = errors.Wrap(err, "redis get goods price error：")

	return
}
