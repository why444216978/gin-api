package conn

import (
	"context"
	"errors"
	"time"

	goods_respository "github.com/why444216978/gin-api/internal/goods/respository"
	goods_service "github.com/why444216978/gin-api/internal/goods/service"
	redis_cache "github.com/why444216978/gin-api/library/cache/redis"
	redis_lock "github.com/why444216978/gin-api/library/lock/redis"
	"github.com/why444216978/gin-api/resource"
	"github.com/why444216978/gin-api/response"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/orm"
	"golang.org/x/sync/errgroup"
)

func Do(c *gin.Context) {
	var (
		err   error
		goods goods_respository.Test
	)

	ctx := c.Request.Context()
	db := resource.TestDB.DB

	db = db.WithContext(ctx).Begin()

	defer func() {
		if err != nil {
			db.WithContext(ctx).Rollback()
			// resource.ServiceLogger.Error
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

	_, err = orm.Insert(ctx, db, &goods_respository.Test{
		GoodsId: 333,
		Name:    "a",
	})
	if err != nil {
		return
	}

	where := map[string]interface{}{"goods_id": 333}
	update := map[string]interface{}{"name": 333}

	_, err = orm.Update(ctx, db, &goods_respository.Test{}, where, update)
	if err != nil {
		return
	}

	_, err = orm.Delete(ctx, db, &goods_respository.Test{}, where)
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

	lock, err := redis_lock.New(resource.RedisCache)
	if err != nil {
		return
	}

	cache, err := redis_cache.New(resource.RedisCache, lock)
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
