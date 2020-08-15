package dao

import (
	"gin-frame/dao/origin_price_dao"
)

type DaoFactory struct{}

func (factory *DaoFactory) GetInstance(name string) map[string]interface{} {
	instances := make(map[string]interface{})

	switch name {
	case "OriginPriceDao":
		instances[name] = origin_price_dao.NewObj()
	default:
		panic("dao name error")
	}
	return instances
}
