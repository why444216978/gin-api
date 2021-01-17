package conn

import (
	"context"
	"gin-api/controllers"
	"gin-api/services/goods_service"
	"gin-api/services/test_service"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
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

	price := 0
	name := ""

	g, _ := errgroup.WithContext(context.TODO())
	g.Go(func() (err error) {
		name, err = self.goodsService.GetGoodsName(self.C, goodsId)
		if err != nil {
			return err
		}
		return nil
	})
	g.Go(func() (err error) {
		price, err = self.goodsService.GetGoodsPrice(self.C, goodsId)
		if err != nil {
			return err
		}
		return nil
	})

	err := g.Wait()
	if err != nil {
		panic(err)
	}

	goods["name"] = name
	goods["price"] = price
	self.Data["goods"] = goods

}

func (self *ConnController) setData() {

}
