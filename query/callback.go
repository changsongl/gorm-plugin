package query

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

// Handler is gorm v2 callback function
type Handler func(db *gorm.DB)

// Interceptor function receives callback name,
// and return a function to wrap the Handler with
// next Handler.
type Interceptor func(string) func(next Handler) Handler

// slowQueryMetricInterceptor return a slow query Interceptor.
func slowQueryMetricInterceptor(max time.Duration, metric *slowMetric) Interceptor {
	return func(cbName string) func(next Handler) Handler {
		return func(originHandler Handler) Handler {
			return func(db *gorm.DB) {
				start := time.Now()
				originHandler(db)
				cost := time.Since(start)
				metric.timeQuery(db.Statement.Table, cost)

				if cost < max {
					return
				}
				metric.incSlowQuery(db.Statement.Table, cbName)
			}
		}
	}
}

// errorQueryMetricInterceptor return a error query Interceptor.
func errorQueryMetricInterceptor(metric *errorMetric) Interceptor {
	return func(cbName string) func(next Handler) Handler {
		return func(originHandler Handler) Handler {
			return func(db *gorm.DB) {
				originHandler(db)
				if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
					metric.incErrorQuery(db.Statement.Table, cbName)
				}
			}
		}
	}
}

type metricCallback interface {
	getDb() *gorm.DB
	getInterceptor() Interceptor
}

// slowCallback is for replace callback with slow metrics.
type slowCallback struct {
	db            *gorm.DB
	slowThreshold time.Duration
	metric        *slowMetric
	interceptor   Interceptor
}

// newSlowCallback return a slowCallback.
func newSlowCallback(db *gorm.DB, slowThreshold time.Duration, metric *slowMetric) metricCallback {
	return &slowCallback{
		db:            db,
		slowThreshold: slowThreshold,
		metric:        metric,
		interceptor:   slowQueryMetricInterceptor(slowThreshold, metric),
	}
}

func (r *slowCallback) getDb() *gorm.DB {
	return r.db
}

func (r *slowCallback) getInterceptor() Interceptor {
	return r.interceptor
}

// errorCallback is for replace callback with error metrics.
type errorCallback struct {
	db          *gorm.DB
	metric      *errorMetric
	interceptor Interceptor
}

// newErrorCallback return a errorCallback.
func newErrorCallback(db *gorm.DB, metric *errorMetric) metricCallback {
	return &errorCallback{
		db:          db,
		metric:      metric,
		interceptor: errorQueryMetricInterceptor(metric),
	}
}

func (e errorCallback) getDb() *gorm.DB {
	return e.db
}

func (e errorCallback) getInterceptor() Interceptor {
	return e.interceptor
}

// replaceAllCallback is for crazy ladder function calls,
// which replace all callbacks.
func replaceAllCallback(m metricCallback) {
	replaceRowCallback(
		replaceRawCallback(
			replaceQueryCallback(
				replaceDeleteCallback(
					replaceUpdateCallback(
						replaceCreateCallback(m))))))
}

// replaceCreateCallback replace create callback.
func replaceCreateCallback(m metricCallback) metricCallback {
	db, interceptor, cb := m.getDb(), m.getInterceptor(), "gorm:create"
	err := db.Callback().Create().Replace(cb, interceptor(cb)(db.Callback().Create().Get(cb)))
	panicCallbackError(err)
	return m
}

// replaceDeleteCallback replace delete callback.
func replaceDeleteCallback(m metricCallback) metricCallback {
	db, interceptor, cb := m.getDb(), m.getInterceptor(), "gorm:delete"
	err := db.Callback().Delete().Replace(cb, interceptor(cb)(db.Callback().Delete().Get(cb)))
	panicCallbackError(err)
	return m
}

// replaceDeleteCallback replace delete callback.
func replaceQueryCallback(m metricCallback) metricCallback {
	db, interceptor, cb := m.getDb(), m.getInterceptor(), "gorm:query"
	err := db.Callback().Query().Replace(cb, interceptor(cb)(db.Callback().Query().Get(cb)))
	panicCallbackError(err)
	return m
}

// replaceUpdateCallback replace update callback.
func replaceUpdateCallback(m metricCallback) metricCallback {
	db, interceptor, cb := m.getDb(), m.getInterceptor(), "gorm:update"
	err := db.Callback().Update().Replace(cb, interceptor(cb)(db.Callback().Update().Get(cb)))
	panicCallbackError(err)
	return m
}

// replaceRowCallback replace row callback.
func replaceRowCallback(m metricCallback) metricCallback {
	db, interceptor, cb := m.getDb(), m.getInterceptor(), "gorm:row"
	err := db.Callback().Row().Replace(cb, interceptor(cb)(db.Callback().Row().Get(cb)))
	panicCallbackError(err)
	return m
}

// replaceRawCallback replace raw callback.
func replaceRawCallback(m metricCallback) metricCallback {
	db, interceptor, cb := m.getDb(), m.getInterceptor(), "gorm:raw"
	err := db.Callback().Raw().Replace(cb, interceptor(cb)(db.Callback().Raw().Get(cb)))
	panicCallbackError(err)
	return m
}

// panicCallbackError panic with error.
func panicCallbackError(err error) {
	if err != nil {
		panic(fmt.Sprintf("SlowQueryCallback plugin failed: %s\n", err.Error()))
	}
}
