package conn

import (
	"gin-api/libraries/logging"
	"gin-api/response"
	"gin-api/services/goods_service"
	"gin-api/services/test_service"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
)

func Do(c *gin.Context) {
	if err := c.Request.Context().Err(); err != nil {
		logging.ErrorCtx(c, map[string]interface{}{"err": err.Error()})
		response.Response(c, response.CODE_SUCCESS, nil, "")
		return
	}

	goods, _ := test_service.Instance.GetFirstRow(c, true)
	g, _ := errgroup.WithContext(c.Request.Context())
	g.Go(func() (err error) {
		goods.Name, err = goods_service.Instance.GetGoodsName(c, goods.Id)
		if err != nil {
			return err
		}
		return nil
	})
	err := g.Wait()
	if err != nil {
		logging.ErrorCtx(c, map[string]interface{}{"err": err.Error()})
		response.Response(c, response.CODE_SERVER, goods, "")
		return
	}

	response.Response(c, response.CODE_SUCCESS, goods, "")
}
