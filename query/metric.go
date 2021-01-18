package query

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// metric labels key
const (
	labelDbName       = "db_name"
	labelTableName    = "table_name"
	labelCallbackName = "callback"
)

// slowMetric has time histogram and slow query counter.
type slowMetric struct {
	counter   *prometheus.CounterVec
	histogram *prometheus.HistogramVec
}

// errorMetric has error counter.
type errorMetric struct {
	counter *prometheus.CounterVec
}

// newSlowMetric return a slowMetric
func newSlowMetric(namePrefix, namespace, dbName string) *slowMetric {
	slowCounter := slowMetric{
		counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        fmt.Sprintf("%s_slow_query_count", namePrefix),
				Namespace:   namespace,
				Help:        "gorm-plugin: slow query counter",
				ConstLabels: getDBConstLabel(dbName),
			},
			[]string{labelTableName, labelCallbackName},
		),
		histogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        fmt.Sprintf("%s_query_time", namePrefix),
				Namespace:   namespace,
				Help:        "gorm-plugin: slow query timeQuery histogram (unit: second)",
				Buckets:     []float64{.05, .1, .25, .5, 1, 2.5, 5, 10},
				ConstLabels: getDBConstLabel(dbName),
			},
			[]string{labelTableName},
		),
	}

	return &slowCounter
}

// newErrorMetric return a errorMetric
func newErrorMetric(namePrefix, namespace, dbName string) *errorMetric {
	errorCounter := errorMetric{
		counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        fmt.Sprintf("%s_error_count", namePrefix),
				Namespace:   namespace,
				Help:        "gorm-plugin: error counter",
				ConstLabels: getDBConstLabel(dbName),
			},
			[]string{labelTableName, labelCallbackName},
		),
	}

	return &errorCounter
}

// incSlowQuery increase slow query counter by 1.
func (s *slowMetric) incSlowQuery(table, cbName string) {
	labels := getDbAndTableMap(table)
	labels[labelCallbackName] = cbName
	s.counter.With(labels).Inc()
}

// timeQuery set query execution time histogram.
func (s *slowMetric) timeQuery(table string, cost time.Duration) {
	labels := getDbAndTableMap(table)
	s.histogram.With(labels).Observe(float64(cost) / float64(time.Second))
}

// incErrorQuery increase error query counter by 1.
func (s *errorMetric) incErrorQuery(table, cbName string) {
	labels := getDbAndTableMap(table)
	labels[labelCallbackName] = cbName
	s.counter.With(labels).Inc()
}

// getDbAndTableMap return a map for prometheus labels.
func getDbAndTableMap(table string) map[string]string {
	return map[string]string{
		labelTableName: table,
	}
}

// getDBConstLabel return label const label
func getDBConstLabel(db string) map[string]string {
	return map[string]string{
		labelDbName: db,
	}
}
