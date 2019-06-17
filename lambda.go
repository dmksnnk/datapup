package datapup

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

type Lambda struct {
	clock Clock
}

func NewLambda() *Lambda {
	return &Lambda{clock: realClock{}}
}

func (l *Lambda) format(name string, value interface{}, metric MetricType, tags ...string) string {
	var strValue string
	switch v := value.(type) {
	case int64:
		strValue = strconv.FormatInt(v, 10)
	case float64:
		// same as https://github.com/DataDog/datadog-go/blob/master/statsd/statsd.go#L206
		strValue = strconv.FormatFloat(v, 'f', 6, 64)
	case statsd.ServiceCheckStatus:
		// as in https://github.com/DataDog/datadog-go/blob/f6e76752dd64e7329d6b314f92a08748a78c2250/statsd/statsd.go#L678
		strValue = strconv.FormatInt(int64(v), 10)
	case string:
		strValue = v
	default:
		// do nothing
	}

	return fmt.Sprintf("MONITORING|%d|%s|%s|%s|#%s\n",
		l.clock.Now().UTC().Unix(),
		strValue,
		metric,
		name,
		strings.Join(tags, ","))
}

func (l *Lambda) Send(name string, value interface{}, metric MetricType, rate float64, tags ...string) error {
	if rate < 1 && rand.Float64() > rate {
		return nil
	}

	msg := l.format(name, value, metric, tags...)
	_, err := fmt.Print(msg)
	return err
}
