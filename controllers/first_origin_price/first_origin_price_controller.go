package first_origin_price

import (
	"gin-frame/controllers"
	"gin-frame/services/location_service"
	"gin-frame/services/origin_price_service"
	"gin-frame/services/product_service"
	"sync"

	"github.com/gin-gonic/gin"
)

type FirstOriginPriceController struct {
	controllers.BaseController
	productService     *product_service.ProductService
	locationService    *location_service.LocationService
	OriginPriceService *origin_price_service.OriginPriceService
	Result             map[string]interface{}
}

func Do(c *gin.Context) {
	instance := new(FirstOriginPriceController)
	instance.Init(c)
	instance.load()
	instance.action()
	instance.setData()
	instance.ResultJson()
}

func (self *FirstOriginPriceController) load() {
	self.OriginPriceService = origin_price_service.NewObj()

	self.locationService = location_service.NewObj()

	self.productService = product_service.NewObj()
}

func (self *FirstOriginPriceController) action() {
	origin := self.OriginPriceService.GetFirstRow(self.C, true)
	self.Data["origin"] = origin

	productId := 0
	locationId := 0

	if origin != nil {
		if origin["product_id"] != nil {
			productId = origin["product_id"].(int)
			if origin["breed_id"] != nil {
				productId = origin["breed_id"].(int)
			}
		}

		if origin["location_id"] != nil {
			locationId = origin["location_id"].(int)
		}
	}

	var wg sync.WaitGroup
	product := make(map[string]interface{})
	location := make(map[string]interface{})
	wg.Add(2)
	go func() {
		defer wg.Done()
		product = self.productService.GetProductDetail(self.C, productId)
	}()

	go func() {
		defer wg.Done()
		location = self.locationService.GetLocationDetail(self.C, locationId)
	}()
	wg.Wait()

	self.Data["product"] = product
	self.Data["location"] = location

}

func (self *FirstOriginPriceController) setData() {

}
