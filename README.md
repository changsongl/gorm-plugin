# gorm-plugins
A set of gorm plugins for providing more features. It is using prometheus
metrics to monitor your SQLs.

I will frequently write gorm plugins in this project. Please star it if you like it.

#### Download package with go mod:
`go get github.com/changsongl/gorm-plugin`

### 1. Query
It is a plugin to monitor query execution time and slow query count 
for tables. We can use prometheus client to expose metrics api, and
having a prometheus server to scrape data from it. Eventually, you can
use `Alertmanager` or `Grafana` to monitor your slow query and query running 
time.


````golang
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
	// create gorm v2 db
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// create query plugin
	plugin := query.New(
		query.SlowQueryCallback(query.Config{
			DBName:        "my_test_db",     // can be empty
			NamePrefix:    "myprefix",       // can be empty
			Namespace:     "mynamespace",    // can be empty
			SlowThreshold: time.Millisecond, // slow query time threshold
		}),
		query.ErrorQueryCallback(query.Config{
			DBName:     "my_test_db",  // can be empty
			NamePrefix: "myprefix",    // can be empty
			Namespace:  "mynamespace", // can be empty
		}),
	)
	// using plugin
	if err := db.Use(plugin); err != nil {
		panic(err.Error())
	}

	// running sqls
	db.Raw("SELECT id FROM test WHERE id = ?", 3).Scan(&test{})
	db.Create(&test{Test: "hahaha"})
	db.Where("id = 123232132").First(&test{})                    // record not found
	db.Raw("SELECT id FROM test2 WHERE id = ?", 3).Scan(&test{}) // error

	// register query plugin collectors to prometheus
	prometheus.MustRegister(plugin.MetricsCollectors()...)

	// run prometheus server
	http.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr: ":8080",
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}



// http://127.0.0.1:8080/metrics

# HELP mynamespace_myprefix_error_count gorm-plugin: error counter
# TYPE mynamespace_myprefix_error_count counter
mynamespace_myprefix_error_count{callback="gorm:row",db_name="my_test_db",table_name=""} 1
# HELP mynamespace_myprefix_query_time gorm-plugin: slow query timeQuery histogram (unit: second)
# TYPE mynamespace_myprefix_query_time histogram
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="",le="0.05"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="",le="0.1"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="",le="0.25"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="",le="0.5"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="",le="1"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="",le="2.5"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="",le="5"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="",le="10"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="",le="+Inf"} 2
mynamespace_myprefix_query_time_sum{db_name="my_test_db",table_name=""} 0.018511481
mynamespace_myprefix_query_time_count{db_name="my_test_db",table_name=""} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="test",le="0.05"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="test",le="0.1"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="test",le="0.25"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="test",le="0.5"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="test",le="1"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="test",le="2.5"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="test",le="5"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="test",le="10"} 2
mynamespace_myprefix_query_time_bucket{db_name="my_test_db",table_name="test",le="+Inf"} 2
mynamespace_myprefix_query_time_sum{db_name="my_test_db",table_name="test"} 0.016793621
mynamespace_myprefix_query_time_count{db_name="my_test_db",table_name="test"} 2
# HELP mynamespace_myprefix_slow_query_count gorm-plugin: slow query counter
# TYPE mynamespace_myprefix_slow_query_count counter
mynamespace_myprefix_slow_query_count{callback="gorm:create",db_name="my_test_db",table_name="test"} 1
mynamespace_myprefix_slow_query_count{callback="gorm:query",db_name="my_test_db",table_name="test"} 1
mynamespace_myprefix_slow_query_count{callback="gorm:row",db_name="my_test_db",table_name=""} 2
````

### 2. Explain SQL
It is a plugin for explain the sql you ran. You can set the requirement of explain result.

````sql
CREATE TABLE `test`.`explain_table` (
`id` INT NOT NULL AUTO_INCREMENT,
`room_id` INT UNSIGNED NOT NULL,
`room_name` VARCHAR(45) NOT NULL,
PRIMARY KEY (`id`),
INDEX `idx_room_id` (`room_id` ASC));

INSERT INTO `test`.`explain_table` (`id`, `room_id`, `room_name`) VALUES (NULL, '1', 'room_name_1');
INSERT INTO `test`.`explain_table` (`id`, `room_id`, `room_name`) VALUES (NULL, '2', 'room_name2');
INSERT INTO `test`.`explain_table` (`room_id`, `room_name`) VALUES ('3', 'room_name2');
INSERT INTO `test`.`explain_table` (`room_id`, `room_name`) VALUES ('4', 'room_name2');
INSERT INTO `test`.`explain_table` (`room_id`, `room_name`) VALUES ('5', 'room_name3');
INSERT INTO `test`.`explain_table` (`room_id`, `room_name`) VALUES ('5', 'room_name5');
INSERT INTO `test`.`explain_table` (`room_id`, `room_name`) VALUES ('6', 'room_name9');
INSERT INTO `test`.`explain_table` (`room_id`, `room_name`) VALUES ('7', 'room_name1');
````

````go
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
````

````shell
{Err:index not pass type requirement (range)  Results:[{Id:1 SelectType:SIMPLE Table:explain_table Type:index PossibleKey: Key:PRIMARY KeyLen:4 Ref: Rows:1 Extra:Using where}] SQL:SELECT * FROM `explain_table` WHERE room_name = 'haha' ORDER BY `explain_table`.`id` LIMIT 1}
````

