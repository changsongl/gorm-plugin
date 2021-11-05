package query

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

// MetricPlugin contains prometheus.Plugin interface,
// and MetricsCollectors function for register metrics
// in prometheus.
type MetricPlugin interface {
	Name() string
	Initialize(db *gorm.DB) error
	MetricsCollectors() []prometheus.Collector
}

// metricPlugin implemented MetricPlugin and prometheus.Plugin
// of gorm v2.
type metricPlugin struct {
	opts []Callback
	cols []prometheus.Collector
}

// New a query metric plugin which can monitor query timing
func New(opts ...Callback) MetricPlugin {
	m := &metricPlugin{opts: opts}
	for _, opt := range m.opts {
		m.cols = append(m.cols, opt.getCollector()...)
	}
	return m
}

// Name for metric plugin
func (m *metricPlugin) Name() string {
	return "gorm-plugin:metric"
}

// Initialize replace gorm callbacks
func (m *metricPlugin) Initialize(db *gorm.DB) error {
	for _, opt := range m.opts {
		opt.apply(db)
	}

	return nil
}

// MetricsCollectors return a set of collector for prometheus,
// so you can use prometheus.register to register them.
func (m *metricPlugin) MetricsCollectors() []prometheus.Collector {
	return m.cols
}
