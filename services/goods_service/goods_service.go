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
	redisName     = "default"
	goodsNameKey  = "goods::name::"
	goodsPriceKey = "goods::price::"
)

func init() {
	goods = &GoodsService{}
	goods.redis = redis.GetRedis(redisName)
}

func GetInstance() *GoodsService {
	return goods
}

func (self *GoodsService) GetGoodsPrice(ctx *gin.Context, id int) int {
	data, _ := redigo.Int(self.redis.Do(ctx, "GET", goodsPriceKey+strconv.Itoa(id)))

	return data
}

func (self *GoodsService) GetGoodsName(ctx *gin.Context, id int) string {
	data, _ := redigo.String(self.redis.Do(ctx, "GET", goodsNameKey+strconv.Itoa(id)))

	return data
}

func (self *GoodsService) GetGoodsInfo(ctx *gin.Context, id int) map[string]interface{} {
	data, _ := redigo.String(self.redis.Do(ctx, "GET", goodsNameKey+strconv.Itoa(id)))

	return conversion.JsonToMap(data)
}

func (self *GoodsService) BatchGoodsName(ctx *gin.Context, ids []int) []string {
	var args []interface{}
	for _, v := range ids {
		args = append(args, goodsNameKey+strconv.Itoa(v))
	}

	data, _ := redigo.Strings(self.redis.Do(ctx, "MGET", args...))

	return data
}
