package test_model

import (
	"gin-api/libraries/mysql"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestInterface interface {
	GetFirst(c *gin.Context) (Test, error)
}

type TestModel struct {
	dbMaster *gorm.DB
	dbSlave  *gorm.DB
}

var (
	instance     TestInterface
	instanceOnce sync.Once
)

func New(master, slave *gorm.DB) TestInterface {
	instanceOnce.Do(func() {
		instance = &TestModel{
			dbMaster: master,
			dbSlave:  slave,
		}
	})
	return instance
}

func (m *TestModel) GetFirst(c *gin.Context) (test Test, err error) {
	err = mysql.WithContext(c.Request.Context(), m.dbSlave).Model(&test).Select("*").First(&test).Error

	if test.Id == 0 {
		err = ErrDataEmpty
	}

	return
}
