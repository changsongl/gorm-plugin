package query

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

type MetricPlugin struct {
	opts []Callback
	cols []prometheus.Collector
}

// New a query metric plugin which can monitor query timing
func New(opts ...Callback) *MetricPlugin {
	m := &MetricPlugin{opts: opts}
	for _, opt := range m.opts {
		m.cols = append(m.cols, opt.getCollector())
	}
	return m
}

// Name for metric plugin
func (m *MetricPlugin) Name() string {
	return "gorm:metric"
}

// Initialize replace gorm callbacks
func (m *MetricPlugin) Initialize(db *gorm.DB) error {
	for _, opt := range m.opts {
		opt.apply(db)
	}

	return nil
}

// MetricsCollectors return a set of collector for prometheus,
// so you can use prometheus.register to register them.
func (m *MetricPlugin) MetricsCollectors() []prometheus.Collector {
	return m.cols
}
