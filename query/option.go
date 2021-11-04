package query

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

// Callback interface for query plugin
type Callback interface {
	apply(*gorm.DB)
	getCollector() []prometheus.Collector
}

// cb implemented Callback with callback function
// and prometheus collectors.
type cb struct {
	f    func(db *gorm.DB)
	cols []prometheus.Collector
}

// Config is for slow query option. Setting counter
// name, namespace and slow query threshold. It will stats
// when slow query execution timeQuery is over SlowThreshold, and
// store in counter and histogram in Namespace and with NamePrefix.
type Config struct {
	DBName        string
	Namespace     string
	NamePrefix    string
	SlowThreshold time.Duration
}

// NewCallback return a Callback interface.
func NewCallback(f func(db *gorm.DB), cols ...prometheus.Collector) Callback {
	return cb{f: f, cols: cols}
}

// apply is a implementation function of Callback for cb
func (o cb) apply(db *gorm.DB) {
	o.f(db)
}

// getCollector is a implementation function of Callback for cb
func (o cb) getCollector() []prometheus.Collector {
	return o.cols
}

// SlowQueryCallback returns a Callback. And replace all kind of Callback
// with slow query stats function.
func SlowQueryCallback(c Config) Callback {
	slowMetric := newSlowMetric(c.NamePrefix, c.Namespace, c.DBName)
	cbFunc := func(db *gorm.DB) {
		s := newSlowCallback(db, c.SlowThreshold, slowMetric)
		replaceAllCallback(s)
	}

	return NewCallback(cbFunc, slowMetric.counter, slowMetric.histogram)
}

// ErrorQueryCallback returns a Callback. And replace all kind of Callback
// with error query stats function.
func ErrorQueryCallback(c Config) Callback {
	errorMetric := newErrorMetric(c.NamePrefix, c.Namespace, c.DBName)
	cbFunc := func(db *gorm.DB) {
		e := newErrorCallback(db, errorMetric)
		replaceAllCallback(e)
	}

	return NewCallback(cbFunc, errorMetric.counter)
}
