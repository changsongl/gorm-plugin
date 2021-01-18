package main

import (
	"github.com/changsongl/gorm-plugin/query"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type test struct {
	Id   int64  `gorm:"column:id" json:"id"`     // 实例id
	Test string `gorm:"column:test" json:"test"` // 告警等级
}

func (test) TableName() string {
	return "test"
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.Use(query.New(query.SlowQueryCallback(query.SlowQueryConfig{
		NamePrefix:    "1111",
		Namespace:     "ns",
		SlowThreshold: time.Millisecond,
	})))
	db.Raw("SELECT id FROM test WHERE id = ?", 3).Scan(&test{})
	db.Create(&test{Test: "hahaha"})
}
