package respository

type Test struct {
	Id      uint   `gorm:"column:id" json:"id"`
	GoodsId uint64 `gorm:"column:goods_id" json:"goods_id"`
	Name    string `gorm:"column:name" json:"name"`
}

func (Test) TableName() string {
	return "test"
}
