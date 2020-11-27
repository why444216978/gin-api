package origin_price_service

import (
	"gin-api/dao/origin_price_dao"
	"github.com/gin-gonic/gin"
)

type OriginPriceService struct {
	originPriceDao *origin_price_dao.OriginPriceDao
}

//var onceOriginPriceService sync.Once
var originPriceService *OriginPriceService

func init(){
	originPriceService = &OriginPriceService{}
	originPriceService.originPriceDao = origin_price_dao.NewObj()
}

func NewObj() *OriginPriceService {
	//onceOriginPriceService.Do(func() {
	//	originPriceService = &OriginPriceService{}
	//	originPriceService.originPriceDao = origin_price_dao.NewObj()
	//
	//	log.Printf("new origin_price_service")
	//})

	return originPriceService
}

func (self *OriginPriceService) GetFirstRow(ctx *gin.Context, oCache bool) map[string]interface{} {
	return self.originPriceDao.GetFirstRow(ctx, true)
}
