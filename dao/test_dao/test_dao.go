package test_dao

import (
	"gin-api/models/test/test_model"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-util/conversion"
)

type TestDao struct {
	testModel *test_model.TestModel
}

//var onceOriginPriceDao sync.Once
var testDao *TestDao

func init() {
	testDao = &TestDao{}
	testDao.testModel = test_model.GetInstance()
	log.Printf("new test_dao")
}

func GetInstance() *TestDao {
	return testDao
}

func (self *TestDao) GetFirstRow(c *gin.Context, noCache bool) map[string]interface{} {

	dbRes, err := self.testModel.GetFirst()
	if err != nil {
		panic(err)
	}

	result := make(map[string]interface{})
	if dbRes != nil {
		for _, v := range dbRes {
			result = conversion.StructToMap(v)
			break
		}
	}

	return result
}
