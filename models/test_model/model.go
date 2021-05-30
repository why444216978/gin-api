package test_model

type Test struct {
	//gorm.Model
	Id       int `gorm:"primary_key"`
	Goods_id int
	Name     string
}

func (Test) TableName() string {
	return "test"
}
