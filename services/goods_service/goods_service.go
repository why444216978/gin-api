package goods_service

import (
	"gin-api/libraries/redis"
	"github.com/why444216978/go-util/conversion"
	"strconv"

	"github.com/gin-gonic/gin"

	redigo "github.com/gomodule/redigo/redis"
)

type GoodsService struct {
	redis *redis.RedisDB
}

var goods *GoodsService

//var onceServiceLocation sync.Once

const (
	DB_NAME         = "default"
	GOODS_NAME_KEY  = "goods::name::"
	GOODS_PRICE_KEY = "goods::price::"
)

func init() {
	goods = &GoodsService{}
	goods.redis = redis.GetRedis(DB_NAME)
}

func GetInstance() *GoodsService {
	return goods
}

func (self *GoodsService) GetGoodsPrice(ctx *gin.Context, id int) int {
	data, _ := redigo.Int(self.redis.Do(ctx, "GET", GOODS_PRICE_KEY+strconv.Itoa(id)))

	return data
}

func (self *GoodsService) GetGoodsName(ctx *gin.Context, id int) string {
	data, _ := redigo.String(self.redis.Do(ctx, "GET", GOODS_NAME_KEY+strconv.Itoa(id)))

	return data
}

func (self *GoodsService) GetGoodsInfo(ctx *gin.Context, id int) map[string]interface{} {
	data, _ := redigo.String(self.redis.Do(ctx, "GET", GOODS_NAME_KEY+strconv.Itoa(id)))

	return conversion.JsonToMap(data)
}

func (self *GoodsService) BatchGoodsName(ctx *gin.Context, ids []int) []string {
	var args []interface{}
	for _, v := range ids {
		args = append(args, GOODS_NAME_KEY+strconv.Itoa(v))
	}

	data, _ := redigo.Strings(self.redis.Do(ctx, "MGET", args...))

	return data
}
