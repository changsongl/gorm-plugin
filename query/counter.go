package query

import "github.com/prometheus/client_golang/prometheus"

const (
	labelDbName    = "db_name"
	labelTableName = "table_name"
)

type slowCounter struct {
	prometheus.CounterVec
}

func newSlowCounter(name, namespace string) *slowCounter {
	slowCounter := slowCounter{*prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      name,
			Namespace: namespace,
			Help:      "gorm-plugin: slow query counter",
		},
		[]string{labelDbName, labelTableName},
	)}

	return &slowCounter
}

func (s *slowCounter) inc(db, table string) {
	labels := getDbAndTableMap(db, table)
	s.With(labels).Inc()
}

func getDbAndTableMap(db, table string) map[string]string {
	return map[string]string{
		labelDbName:    db,
		labelTableName: table,
	}
}
