package test_service

import (
	"gin-api/models/test_model"
	"gin-api/resource"
	"sync"

	"github.com/gin-gonic/gin"
)

type TestInterface interface {
	GetFirstRow(ctx *gin.Context, oCache bool) (ret test_model.Test, err error)
}

type TestService struct {
	model test_model.TestInterface
}

var (
	instance     TestInterface
	instanceOnce sync.Once
)

func New() TestInterface {
	instanceOnce.Do(func() {
		instance = &TestService{
			model: test_model.New(resource.TestDB.MasterOrm(), resource.TestDB.SlaveOrm()),
		}
	})
	return instance
}

func (srv *TestService) GetFirstRow(ctx *gin.Context, oCache bool) (ret test_model.Test, err error) {
	return srv.model.GetFirst()
}
