package location_service

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

type LocationService struct {
	xhopField string
	redis     *redis.RedisDB
}

var location *LocationService
var onceServiceLocation sync.Once

const (
	redisName         = "location"
	locationDetailKey = "location::id_detail:"
	locationNameKey   = "location::id_name:"
)

func NewObj() *LocationService {
	onceServiceLocation.Do(func() {
		location = &LocationService{}

		location.xhopField = config.GetXhopField()
		location.redis = redis.GetRedis(redisName)

		log.Printf("new service location")
	})
	return location
}

func (self *LocationService) GetLocationDetail(ctx *gin.Context, id int) map[string]interface{} {
	data, _ := redigo.String(self.redis.Do(ctx, location.xhopField, "GET", locationDetailKey+strconv.Itoa(id)))
	fmt.Println(data)

	return conversion.JsonToMap(data)
}

func (self *LocationService) BatchLocationDetail(ctx *gin.Context, ids []int) []string {
	var args []interface{}
	for _, v := range ids {
		args = append(args, locationDetailKey+strconv.Itoa(v))
	}

	data, _ := redigo.Strings(self.redis.Do(ctx, location.xhopField, "MGET", args...))

	return data
}
