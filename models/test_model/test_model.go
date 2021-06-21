package test_model

import (
	"gin-api/libraries/mysql"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type TestInterface interface {
	GetFirst(c *gin.Context) (Test, error)
}

type TestModel struct {
	db *gorm.DB
}

var (
	instance     TestInterface
	instanceOnce sync.Once
)

func New(db *gorm.DB) TestInterface {
	instanceOnce.Do(func() {
		instance = &TestModel{
			db: db,
		}
	})
	return instance
}

func (m *TestModel) GetFirst(c *gin.Context) (test Test, err error) {
	err = mysql.WithContext(c.Request.Context(), m.db).Clauses(dbresolver.Write).Model(&test).Select("*").First(&test).Error

	if test.Id == 0 {
		err = ErrDataEmpty
	}

	return
}
