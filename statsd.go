package datapup

import (
	"github.com/DataDog/datadog-go/statsd"
	"github.com/pkg/errors"
)

// StatsDSender is interface which should implement StatsD client
type StatsDSender interface {
	Count(string, int64, []string, float64) error
	Histogram(string, float64, []string, float64) error
	Gauge(string, float64, []string, float64) error
	ServiceCheck(*statsd.ServiceCheck) error
}

// StatsD holds reference for datapup client
type StatsD struct {
	client StatsDSender
}

// NewStatsD creates new StatsD DataDog client
func NewStatsD(addr string) (*StatsD, error) {
	client, err := statsd.New(addr, statsd.WithAsyncUDS())
	if err != nil {
		return nil, errors.Wrap(err, "Can't create statsd client")
	}
	return &StatsD{client: client}, nil
}

// Send sends metric to DataDog
func (s *StatsD) Send(name string, value interface{}, metric MetricType, rate float64, tags ...string) error {
	switch metric {
	case Count:
		v, ok := value.(int64)
		if !ok {
			return errors.Errorf("Can't convert %v to int64", value)
		}
		return s.client.Count(name, v, tags, rate)
	case Gauge:
		v, ok := value.(float64)
		if !ok {
			return errors.Errorf("Can't convert %v to float64", value)
		}
		return s.client.Gauge(name, v, tags, rate)
	case Histogram:
		v, ok := value.(float64)
		if !ok {
			return errors.Errorf("Can't convert %v to float64", value)
		}
		return s.client.Histogram(name, v, tags, rate)
	case Check:
		v, ok := value.(statsd.ServiceCheckStatus)
		if !ok {
			return errors.Errorf("Can't convert %v to statsd.ServiceCheckStatus", value)
		}
		sc := statsd.ServiceCheck{
			Name:   name,
			Status: v,
			Tags:   tags,
		}
		return s.client.ServiceCheck(&sc)
	default:
		return errors.Errorf("Bad metric type %v", metric)
	}
}
