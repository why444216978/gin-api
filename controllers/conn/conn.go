package conn

import (
	"context"
	"errors"
	"fmt"
	"gin-api/libraries/logging"
	"gin-api/resource"
	"gin-api/response"
	"gin-api/services/goods_service"
	"gin-api/services/test_service"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
)

func Do(c *gin.Context) {
	goods, _ := test_service.New().GetFirstRow(c, true)
	g, _ := errgroup.WithContext(c.Request.Context())
	g.Go(func() (err error) {
		goods.Name = "golang"
		_, err = goods_service.Instance.BatchGoodsName(c, []int{1, 2})
		if err != nil {
			return err
		}
		return nil
	})
	err := g.Wait()
	if err != nil {
		response.Response(c, response.CODE_SERVER, goods, "")
		return
	}

	resource.Logger.Debug("test conn error msg", logging.MergeHTTPFields(c.Request.Context(), map[string]interface{}{"err": "test err"}))

	data := &Data{}

	err = resource.DefaultRedis.GetData(c.Request.Context(), http.Header{}, "key", 3600, 86400, GetDataA, data)
	fmt.Println(data)
	fmt.Println(err)

	response.Response(c, response.CODE_SUCCESS, goods, "")
}

type Data struct {
	A string `json:"a"`
}

func GetDataA(ctx context.Context, _data interface{}) (err error) {
	data, ok := _data.(*Data)
	if !ok {
		err = errors.New("err assert")
		return
	}
	data.A = "a"
	return
}
