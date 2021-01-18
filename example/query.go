package main

import (
	"net/http"
	"time"

	"github.com/changsongl/gorm-plugin/query"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	plugin := query.New(
		query.SlowQueryCallback(query.Config{
			DBName:        "my_test_db",
			NamePrefix:    "myprefix",
			Namespace:     "mynamespace",
			SlowThreshold: time.Millisecond,
		}),
		query.ErrorQueryCallback(query.Config{
			DBName:        "my_test_db",
			NamePrefix:    "myprefix",
			Namespace:     "mynamespace",
			SlowThreshold: time.Millisecond,
		}),
	)
	if err := db.Use(plugin); err != nil {
		panic(err.Error())
	}
	db.Raw("SELECT id FROM test WHERE id = ?", 3).Scan(&test{})
	db.Create(&test{Test: "hahaha"})
	db.Where("id = 123232132").First(&test{})                    // record not found
	db.Raw("SELECT id FROM test2 WHERE id = ?", 3).Scan(&test{}) // error

	prometheus.MustRegister(plugin.MetricsCollectors()...)
	http.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr: ":8080",
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}
