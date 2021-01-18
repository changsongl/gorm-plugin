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
func newSlowMetric(namePrefix, namespace string) *slowMetric {
	slowCounter := slowMetric{
		counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:      fmt.Sprintf("%s_slow_query_count", namePrefix),
				Namespace: namespace,
				Help:      "gorm-plugin: slow query counter",
			},
			[]string{labelDbName, labelTableName, labelCallbackName},
		),
		histogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:      fmt.Sprintf("%s_query_time", namePrefix),
				Namespace: namespace,
				Help:      "gorm-plugin: slow query timeQuery histogram (unit: second)",
				Buckets:   []float64{.05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{labelDbName, labelTableName},
		),
	}

	return &slowCounter
}

// newErrorMetric return a errorMetric
func newErrorMetric(namePrefix, namespace string) *errorMetric {
	errorCounter := errorMetric{
		counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:      fmt.Sprintf("%s_error_count", namePrefix),
				Namespace: namespace,
				Help:      "gorm-plugin: error counter",
			},
			[]string{labelDbName, labelTableName, labelCallbackName},
		),
	}

	return &errorCounter
}

// incSlowQuery increase slow query counter by 1.
func (s *slowMetric) incSlowQuery(db, table, cbName string) {
	labels := getDbAndTableMap(db, table)
	labels[labelCallbackName] = cbName
	s.counter.With(labels).Inc()
}

// timeQuery set query execution time histogram.
func (s *slowMetric) timeQuery(db, table string, cost time.Duration) {
	labels := getDbAndTableMap(db, table)
	s.histogram.With(labels).Observe(float64(cost) / float64(time.Second))
}

// incErrorQuery increase error query counter by 1.
func (s *errorMetric) incErrorQuery(db, table, cbName string) {
	labels := getDbAndTableMap(db, table)
	labels[labelCallbackName] = cbName
	s.counter.With(labels).Inc()
}

// getDbAndTableMap return a map for prometheus labels.
func getDbAndTableMap(db, table string) map[string]string {
	return map[string]string{
		labelDbName:    db,
		labelTableName: table,
	}
}
