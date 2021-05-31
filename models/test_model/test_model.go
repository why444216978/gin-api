package test_model

import (
	"sync"

	"gorm.io/gorm"
)

type TestInterface interface {
	GetFirst() (Test, error)
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

func (m *TestModel) GetFirst() (test Test, err error) {
	err = m.dbSlave.Model(&test).Select("*").First(&test).Error

	if test.Id == 0 {
		err = ErrDataEmpty
	}

	return
}
