package query

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Handler func(db *gorm.DB)

type Interceptor func(string) func(next Handler) Handler

func slowQueryMetricInterceptor(cbName string, max time.Duration, logFunc LogFunc, counter *slowCounter) func(Handler) Handler {
	return func(originHandler Handler) Handler {
		return func(db *gorm.DB) {
			start := time.Now()
			originHandler(db)
			cost := time.Since(start)

			if cost < max {
				return
			}

			counter.inc(db.Name(), db.Statement.Table)
			if logFunc != nil {
				logFunc(cbName, cost, db)
			}
		}
	}
}

func replaceCreateCallback(c SlowQueryConfig, slowCounter *slowCounter, db *gorm.DB) {
	cb := "gorm:create"
	interFunc := slowQueryMetricInterceptor(cb, c.MaxTime, c.LogFunction, slowCounter)
	err := db.Callback().Create().Replace(cb, interFunc(db.Callback().Create().Get(cb)))
	printCallbackError(err)
}

func replaceDeleteCallback(c SlowQueryConfig, slowCounter *slowCounter, db *gorm.DB) {
	cb := "gorm:delete"
	interFunc := slowQueryMetricInterceptor(cb, c.MaxTime, c.LogFunction, slowCounter)
	err := db.Callback().Delete().Replace(cb, interFunc(db.Callback().Delete().Get(cb)))
	printCallbackError(err)
}

func replaceQueryCallback(c SlowQueryConfig, slowCounter *slowCounter, db *gorm.DB) {
	cb := "gorm:query"
	interFunc := slowQueryMetricInterceptor(cb, c.MaxTime, c.LogFunction, slowCounter)
	err := db.Callback().Query().Replace(cb, interFunc(db.Callback().Query().Get(cb)))
	printCallbackError(err)
}

func replaceUpdateCallback(c SlowQueryConfig, slowCounter *slowCounter, db *gorm.DB) {
	cb := "gorm:update"
	interFunc := slowQueryMetricInterceptor(cb, c.MaxTime, c.LogFunction, slowCounter)
	err := db.Callback().Update().Replace(cb, interFunc(db.Callback().Update().Get(cb)))
	printCallbackError(err)
}

func replaceRowCallback(c SlowQueryConfig, slowCounter *slowCounter, db *gorm.DB) {
	cb := "gorm:row"
	interFunc := slowQueryMetricInterceptor(cb, c.MaxTime, c.LogFunction, slowCounter)
	err := db.Callback().Row().Replace(cb, interFunc(db.Callback().Row().Get(cb)))
	printCallbackError(err)
}

func replaceRawCallback(c SlowQueryConfig, slowCounter *slowCounter, db *gorm.DB) {
	cb := "gorm:raw"
	interFunc := slowQueryMetricInterceptor(cb, c.MaxTime, c.LogFunction, slowCounter)
	err := db.Callback().Raw().Replace(cb, interFunc(db.Callback().Raw().Get(cb)))
	printCallbackError(err)
}

func printCallbackError(err error) {
	if err != nil {
		panic(fmt.Sprintf("SlowQueryCallback plugin failed: %s\n", err.Error()))
	}
}
