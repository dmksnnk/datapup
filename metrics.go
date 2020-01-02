package datapup

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
)

// Metic hold reference to a client, its name and tags
type Metric struct {
	client *Client
	name   string
	tags   []string
}

// NewMetric creates new metric to send to DataDog, it requires name and optional tags
func (c *Client) NewMetric(name string, tags ...string) *Metric {
	m := Metric{
		client: c,
		name:   name,
		tags:   tags,
	}

	return &m
}

// Tag formats tag into "key:value"
func Tag(key, value string) string {
	return fmt.Sprintf("%s:%s", key, value)
}

// WithTag retruns new metric with additional tag
func (m *Metric) WithTag(t string) *Metric {
	tags := append(m.tags, t)
	return m.client.NewMetric(m.name, tags...)
}

// WithTags retruns new metric with additional tags
func (m *Metric) WithTags(tags ...string) *Metric {
	newTags := append(m.tags, tags...)
	return m.client.NewMetric(m.name, newTags...)
}

// Count tracks how many times something happened per second
func (m *Metric) Count(value int64) error {
	return m.client.Report(m.name, value, Count, 1, m.tags...)
}

// Incr is just Count of 1
func (m *Metric) Incr() error {
	return m.Count(1)
}

// Decr is just Count of -1
func (m *Metric) Decr() error {
	return m.Count(-1)
}

// Gauge measures the value of a metric at a particular time
func (m *Metric) Gauge(value float64) error {
	return m.client.Report(m.name, value, Gauge, *m.client.rate, m.tags...)
}

// Histogram tracks the statistical distribution of a set of values on each host
func (m *Metric) Histogram(value float64) error {
	return m.client.Report(m.name, value, Histogram, *m.client.rate, m.tags...)
}

// Check sends an serviceCheck with status statsd.ServiceCheckStatus
func (m *Metric) Check(status statsd.ServiceCheckStatus) error {
	return m.client.Report(m.name, status, Check, *m.client.rate, m.tags...)
}
