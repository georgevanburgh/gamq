package gamq

// Metric represents StatsD metric
type Metric struct {
	Value int64
	Name  string
	Type  string
}

// NewMetric constructs a new Metric object from the given parameters
func NewMetric(givenName string, givenType string, givenValue int64) *Metric {
	return &Metric{Name: givenName, Value: givenValue, Type: givenType}
}
