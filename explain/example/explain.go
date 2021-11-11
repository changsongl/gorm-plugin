package main

import (
	"fmt"
	"github.com/changsongl/gorm-plugin/explain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type test struct {
	ID       int64  `gorm:"column:id" json:"id"`
	RoomID   int64  `gorm:"column:room_id" json:"room_id"`
	RoomName string `gorm:"column:room_name" json:"room_name"`
}

func (test) TableName() string {
	return "explain_table"
}

func main() {
	// create gorm v2 db
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// create query plugin
	plugin := explain.New(
		explain.CallBackFuncOption(func(result explain.CallBackResult) {
			fmt.Printf("%+v\n", result)
		}),
		explain.TypeLevelOption(explain.ResultTypeRange),
	)

	// using plugin
	if err := db.Use(plugin); err != nil {
		panic(err.Error())
	}

	db.Where("room_name = 'haha'").First(&test{})
}
