package test_service

import (
	"gin-api/models/test_model"

	"github.com/gin-gonic/gin"
)

type TestInterface interface {
	GetFirstRow(ctx *gin.Context, oCache bool) (ret test_model.Test, err error)
}

type TestService struct{}

//var onceOriginPriceService sync.Once
var Instance TestInterface

func init() {
	Instance = &TestService{}
}

func (srv *TestService) GetFirstRow(ctx *gin.Context, oCache bool) (ret test_model.Test, err error) {
	return test_model.Instance.GetFirst()
}
