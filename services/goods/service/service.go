package goods

import (
	"context"
	"strconv"

	"github.com/why444216978/gin-api/resource"
	"github.com/why444216978/gin-api/services/goods/respository"
	"github.com/why444216978/go-util/orm"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type GoodsInterface interface {
	GetGoodsName(ctx context.Context, id int) (string, error)
	CrudGoods(ctx context.Context) (goods respository.Test, err error)
}

var Instance GoodsInterface

type GoodsService struct{}

func init() {
	Instance = &GoodsService{}
}

const (
	GOODS_NAME_KEY  = "goods::name::"
	GOODS_PRICE_KEY = "goods::price::"
)

func (gs *GoodsService) GetGoodsName(ctx context.Context, id int) (string, error) {
	data, err := resource.RedisCache.Get(ctx, GOODS_NAME_KEY+strconv.Itoa(id)).Result()
	if err == redis.Nil {
		err = nil
	}
	if err != nil {
		err = errors.Wrap(err, "redis get goods price error：")
	}
	return data, err
}

func (gs *GoodsService) CrudGoods(ctx context.Context) (goods respository.Test, err error) {
	db := resource.TestDB.DB.WithContext(ctx).Begin()

	defer func() {
		if err != nil {
			db.WithContext(ctx).Rollback()
			return
		}
		err = db.WithContext(ctx).Commit().Error
	}()

	err = db.WithContext(ctx).Select("*").First(&goods).Error
	if err != nil {
		return
	}

	_, err = orm.Insert(ctx, db, &respository.Test{
		GoodsId: 333,
		Name:    "a",
	})
	if err != nil {
		return
	}

	where := map[string]interface{}{"goods_id": 333}
	update := map[string]interface{}{"name": 333}

	_, err = orm.Update(ctx, db, &respository.Test{}, where, update)
	if err != nil {
		return
	}

	_, err = orm.Delete(ctx, db, &respository.Test{}, where)
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

	return
}
