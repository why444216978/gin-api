package test_service

import (
	"gin-api/dao/test_dao"
	"github.com/gin-gonic/gin"
	"log"
)

type TestService struct {
	testDao *test_dao.TestDao
}

//var onceOriginPriceService sync.Once
var testService *TestService

func init() {
	testService = &TestService{}
	testService.testDao = test_dao.GetInstance()
	log.Printf("new test_service")
}

func GetInstance() *TestService {
	return testService
}

func (self *TestService) GetFirstRow(ctx *gin.Context, oCache bool) map[string]interface{} {
	return self.testDao.GetFirstRow(ctx, true)
}
