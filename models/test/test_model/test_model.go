package test_model

import (
	"fmt"
	"gin-api/libraries/mysql"
	"gin-api/models/base"
)

type Test struct {
	//gorm.Model
	Id       int `gorm:"primary_key"`
	Goods_id int
	Name     string
}

const DB_NAME = "default"

func (Test) TableName() string {
	return "test"
}

type TestModel struct {
	base.BaseModel
}

//var onceOriginPriceModel sync.Once
var testModel *TestModel
var dbInstance *mysql.DB

func init() {
	testModel = &TestModel{}
	dbInstance = testModel.GetConn(DB_NAME)
}

func GetInstance() *TestModel {
	return testModel
}

func (instance *TestModel) GetFirst() []Test {
	test := []Test{}
	orm := dbInstance.SlaveOrm()

	dbRes := orm.First(&test)
	fmt.Println(dbRes)
	fmt.Println(*dbRes)
	fmt.Println(test)

	instance.CheckRes(dbRes)
	return test
}
