package test_model

import (
	"gin-api/resource"
)

type TestInterface interface {
	GetFirst() (Test, error)
}

var Instance TestInterface

type TestModel struct{}

func init() {
	Instance = &TestModel{}
}

func (m *TestModel) GetFirst() (test Test, err error) {
	orm := resource.TestDB.SlaveOrm()
	err = orm.Model(&test).Select("*").First(&test).Error

	return
}
