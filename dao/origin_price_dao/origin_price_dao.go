package origin_price_dao

import (
	"gin-frame/models/hangqing/origin_price_model"
	"log"
	"sync"

	"gin-frame/libraries/util/conversion"

	"github.com/gin-gonic/gin"
)

type OriginPriceDao struct {
	originPriceModel *origin_price_model.OriginPriceModel
}

var onceOriginPriceDao sync.Once
var originPriceDao *OriginPriceDao

func NewObj() *OriginPriceDao {
	onceOriginPriceDao.Do(func() {
		originPriceDao = &OriginPriceDao{}
		originPriceDao.originPriceModel = origin_price_model.NewOriginPriceModel()
		log.Printf("new origin_price_dao")
	})

	return originPriceDao
}

func (self *OriginPriceDao) GetFirstRow(c *gin.Context, noCache bool) map[string]interface{} {
	dbRes := self.originPriceModel.GetFirst()

	result := make(map[string]interface{})
	if dbRes != nil {
		for _, v := range dbRes {
			result = conversion.StructToMap(v)
			break
		}
	}

	return result
}
