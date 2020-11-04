package product_service

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"gin-api/libraries/config"
	"gin-api/libraries/redis"
	"gin-api/libraries/util/conversion"

	"github.com/gin-gonic/gin"

	redigo "github.com/gomodule/redigo/redis"
)

type ProductService struct {
	xhopField string
	redis     *redis.RedisDB
}

var product *ProductService
var onceServiceLocation sync.Once

const (
	redisName        = "product"
	productDetailKey = "product::id_detail:"
	productNameKey   = "product::id_name:"
)

func NewObj() *ProductService {
	onceServiceLocation.Do(func() {
		product = &ProductService{}

		product.xhopField = config.GetXhopField()
		product.redis = redis.GetRedis(redisName)

		log.Printf("new service product")
	})
	return product
}

func (self *ProductService) GetProductDetail(ctx *gin.Context, id int) map[string]interface{} {
	data, _ := redigo.String(self.redis.Do(ctx, product.xhopField, "GET", productDetailKey+strconv.Itoa(id)))
	fmt.Println(data)

	return conversion.JsonToMap(data)
}

func (self *ProductService) BatchProductDetail(ctx *gin.Context, ids []int) []string {
	var args []interface{}
	for _, v := range ids {
		args = append(args, productDetailKey+strconv.Itoa(v))
	}

	data, _ := redigo.Strings(self.redis.Do(ctx, product.xhopField, "MGET", args...))

	return data
}
