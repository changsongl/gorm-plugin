package query

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

type Callback interface {
	apply(*gorm.DB)
	getCollector() prometheus.Collector
}

type LogFunc func(cbName string, cost time.Duration, db *gorm.DB)

var DefaultLogFunc = func(cbName string, cost time.Duration, db *gorm.DB) {
	query, vars := db.Statement.SQL.String(), db.Statement.Vars
	sql := BindSQL(query, vars)

	fmt.Printf("[slow sql][%s] err = %v, time = %s, db = %s, table = %s, sql = %s\n",
		cbName, db.Error, cost, db.Name(), db.Statement.Table, sql)
}

type cb struct {
	f   func(db *gorm.DB)
	col prometheus.Collector
}

type SlowQueryConfig struct {
	CounterNamespace string
	CounterName      string
	MaxTime          time.Duration
	LogFunction      LogFunc
}

func NewCallback(f func(db *gorm.DB), col prometheus.Collector) Callback {
	return cb{f: f, col: col}
}

func (o cb) apply(db *gorm.DB) {
	o.f(db)
}

func (o cb) getCollector() prometheus.Collector {
	return o.col
}

func SlowQueryCallback(c SlowQueryConfig) Callback {
	slowCounter := newSlowCounter(c.CounterName, c.CounterNamespace)
	cbFunc := func(db *gorm.DB) {
		replaceCreateCallback(c, slowCounter, db)
		replaceUpdateCallback(c, slowCounter, db)
		replaceDeleteCallback(c, slowCounter, db)
		replaceQueryCallback(c, slowCounter, db)
		replaceRawCallback(c, slowCounter, db)
		replaceRowCallback(c, slowCounter, db)
	}

	return NewCallback(cbFunc, slowCounter)
}

func ErrorCallback() Callback {
	return nil
}
