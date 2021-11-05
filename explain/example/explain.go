package main

import (
	"github.com/changsongl/gorm-plugin/explain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type test struct {
	Id   int64  `gorm:"column:id" json:"id"`
	Test string `gorm:"column:test" json:"test"`
}

func (test) TableName() string {
	return "test"
}

func main() {
	// create gorm v2 db
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// create query plugin
	plugin := explain.New()
	// using plugin
	if err := db.Use(plugin); err != nil {
		panic(err.Error())
	}

	// running sqls
	db.Raw("SELECT id FROM test WHERE id = ?", 3).Scan(&test{})
	db.Create(&test{Test: "hahaha"})
	db.Where("id = 123232132").First(&test{})                    // record not found
	db.Raw("SELECT id FROM test2 WHERE id = ?", 3).Scan(&test{}) // error
}
