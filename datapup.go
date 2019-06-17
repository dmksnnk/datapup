package datapup

import (
	"fmt"
)

type MetricType string

const (
	Count          = MetricType("count")
	Gauge          = MetricType("gauge")
	Histogram      = MetricType("histogram")
	Check          = MetricType("check")
	DefaultAddress = "localhost:8126"
)

// Sender is interface which must implemented to send metrics
type Sender interface {
	Send(string, interface{}, MetricType, float64, ...string) error
}

// A Client is a handle for sending messages to datadog
type Client struct {
	namespace string
	tags      []string
	rate      *float64
	sender    Sender
}

// Op is option to configure datapup client
type Op func(*Client)

// WithTag adds default tags to client
func WithTag(key, value string) Op {
	return func(c *Client) {
		c.tags = append(c.tags, Tag(key, value))
	}
}

// WithEnvironment sets default "env" tag to all metrics
func WithEnvironment(env string) Op {
	return func(c *Client) {
		c.tags = append(c.tags, Tag("env", env))
	}
}

// WithRate configures rate at which send metrics to
func WithRate(rate float64) Op {
	return func(c *Client) {
		c.rate = &rate
	}
}

// WithSender configures sender which sends metrics to datadog
func WithSender(s Sender) Op {
	return func(c *Client) {
		c.sender = s
	}
}

// New creates new datapup client
func New(namespace string, ops ...Op) *Client {
	c := Client{namespace: namespace}

	for _, op := range ops {
		op(&c)
	}
	// if rate wasn't set up
	if c.rate == nil {
		def := float64(1)
		c.rate = &def
	}
	// default client
	if c.sender == nil {
		client, _ := NewStatsD(DefaultAddress)
		c.sender = client
	}

	return &c
}

// Report sends mertic to DataDog
func (c *Client) Report(name string, value interface{}, metric MetricType, rate float64, tags ...string) error {
	fullName := fmt.Sprintf("%s.%s", c.namespace, name) // join metric name
	tags = append(c.tags, tags...)                      // add new tags

	return c.sender.Send(fullName, value, metric, rate, tags...)
}
