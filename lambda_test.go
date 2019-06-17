package datapup

import (
	"testing"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

type fakeClock struct{}

func (fakeClock) Now() time.Time {
	return time.Unix(12345, 0)
}

func TestLambda_format(t *testing.T) {
	type args struct {
		name   string
		value  interface{}
		metric MetricType
		tags   []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "int64",
			args: args{
				name:   "test.name",
				value:  int64(1),
				metric: Count,
				tags:   []string{"aaa:aaa", "bbb:bbb"},
			},
			want: "MONITORING|12345|1|count|test.name|#aaa:aaa,bbb:bbb\n",
		},
		{
			name: "float64",
			args: args{
				name:   "test.name",
				value:  float64(1.1),
				metric: Count,
				tags:   []string{"aaa:aaa", "bbb:bbb"},
			},
			want: "MONITORING|12345|1.100000|count|test.name|#aaa:aaa,bbb:bbb\n",
		},
		{
			name: "string",
			args: args{
				name:   "test.name",
				value:  "some value",
				metric: Count,
				tags:   []string{"aaa:aaa", "bbb:bbb"},
			},
			want: "MONITORING|12345|some value|count|test.name|#aaa:aaa,bbb:bbb\n",
		},
		{
			name: "service check",
			args: args{
				name:   "test.name",
				value:  statsd.Warn,
				metric: Check,
				tags:   []string{"aaa:aaa", "bbb:bbb"},
			},
			want: "MONITORING|12345|1|check|test.name|#aaa:aaa,bbb:bbb\n",
		},
		{
			name: "default",
			args: args{
				name:   "test.name",
				value:  []byte("AAA"),
				metric: Count,
				tags:   []string{"aaa:aaa", "bbb:bbb"},
			},
			want: "MONITORING|12345||count|test.name|#aaa:aaa,bbb:bbb\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Lambda{clock: fakeClock{}}
			if got := l.format(tt.args.name, tt.args.value, tt.args.metric, tt.args.tags...); got != tt.want {
				t.Errorf("Lambda.format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSend(t *testing.T) {
	l := NewLambda()
	err := l.Send("test.name", 1, Count, 1, "aaa", "bbb")
	if err != nil {
		t.Errorf("Want err to be nil, but get %v", err)
	}
}
