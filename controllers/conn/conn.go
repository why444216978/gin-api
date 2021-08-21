package conn

import (
	"context"
	"errors"
	"gin-api/libraries/cache"
	"gin-api/libraries/lock"
	"gin-api/models/test_model"
	"gin-api/resource"
	"gin-api/response"
	"gin-api/services/goods_service"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"

	"github.com/why444216978/go-util/orm"
)

func Do(c *gin.Context) {
	var (
		err   error
		goods test_model.Test
	)

	ctx := c.Request.Context()
	db := resource.TestDB.DB

	db = db.WithContext(ctx).Begin()

	defer func() {
		if err != nil {
			db.WithContext(ctx).Rollback()
			response.Response(c, response.CodeServer, goods, err.Error())
			return
		}
		err = db.WithContext(ctx).Commit().Error
		if err != nil {
			response.Response(c, response.CodeServer, goods, err.Error())
			return
		}

		response.Response(c, response.CodeSuccess, goods, "")
	}()

	err = db.WithContext(ctx).Select("*").First(&goods).Error
	if err != nil {
		return
	}

	_, err = orm.Insert(ctx, db, &test_model.Test{
		GoodsId: 333,
		Name:    "a",
	})
	if err != nil {
		return
	}

	where := map[string]interface{}{"goods_id": 333}
	update := map[string]interface{}{"name": 333}

	_, err = orm.Update(ctx, db, &test_model.Test{}, where, update)
	if err != nil {
		return
	}

	_, err = orm.Delete(ctx, db, &test_model.Test{}, where)
	if err != nil {
		return
	}

	var name string
	err = db.WithContext(ctx).Table("test").Where("id = ?", 1).Select("name").Row().Scan(&name)
	if err != nil {
		return
	}

	err = db.WithContext(ctx).Raw("select * from test where id = 1 limit 1").Scan(&goods).Error
	if err != nil {
		return
	}

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		goods.Name = "golang"
		_, err = goods_service.Instance.GetGoodsName(c, 1)
		return
	})
	err = g.Wait()
	if err != nil {
		return
	}

	data := &Data{}

	lock, err := lock.New(resource.RedisCache)
	if err != nil {
		return
	}

	cache, err := cache.New(resource.RedisCache, lock)
	if err != nil {
		return
	}

	err = cache.GetData(ctx, "key", time.Hour, time.Hour, GetDataA, data)
	if err != nil {
		return
	}

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
