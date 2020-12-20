package conn

import (
	"gin-api/controllers"
	"gin-api/services/goods_service"
	"gin-api/services/test_service"
	"github.com/gin-gonic/gin"
	"sync"
)

type ConnController struct {
	controllers.BaseController
	goodsService *goods_service.GoodsService
	TestService  *test_service.TestService
	Result       map[string]interface{}
}

func Do(c *gin.Context) {
	instance := new(ConnController)
	instance.Init(c)
	instance.load()
	instance.action()
	instance.setData()
	instance.ResultJson()
}

func (self *ConnController) load() {
	self.TestService = test_service.GetInstance()

	self.goodsService = goods_service.GetInstance()
}

func (self *ConnController) action() {
	goods := self.TestService.GetFirstRow(self.C, true)

	goodsId := goods["goods_id"].(int)

	var wg sync.WaitGroup
	price := 0
	name := ""
	wg.Add(2)
	go func() {
		defer wg.Done()
		name = self.goodsService.GetGoodsName(self.C, goodsId)
	}()

	go func() {
		defer wg.Done()
		price = self.goodsService.GetGoodsPrice(self.C, goodsId)
	}()
	wg.Wait()

	goods["name"] = name
	goods["price"] = price
	self.Data["goods"] = goods

}

func (self *ConnController) setData() {

}
