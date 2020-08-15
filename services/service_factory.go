package service

import (
	"gin-frame/services/location_service"
	"gin-frame/services/origin_price_service"
	"gin-frame/services/product_service"
)

type ServiceFactory struct{}

func (factory *ServiceFactory) GetInstance(name string) map[string]interface{} {
	instances := make(map[string]interface{})

	switch name {
	case "OriginPriceService":
		instances[name] = origin_price_service.NewObj()
	case "LocationService":
		instances[name] = location_service.NewObj()
	case "ProductService":
		instances[name] = product_service.NewObj()
	default:
		panic("service name error")
	}
	return instances
}
