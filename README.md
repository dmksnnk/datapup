# Datapup

[![Go Report Card](https://goreportcard.com/badge/github.com/dmksnnk/datapup)](https://goreportcard.com/report/github.com/dmksnnk/datapup)
[![pipeline status](https://gitlab.com/aspidima/datapup/badges/master/pipeline.svg)](https://gitlab.com/aspidima/datapup/commits/master)
[![coverage report](https://gitlab.com/aspidima/datapup/badges/master/coverage.svg)](https://gitlab.com/aspidima/datapup/commits/master)
[![GoDoc](https://img.shields.io/badge/GoDoc-referece-blue.svg?style=flat)](https://godoc.org/github.com/dmksnnk/datapup)

Datapup is metric-centered tool for sending metrics to DataDog.

## Usage

### Client

When creating datapup client, in `WithSender` you may to pass a sender which implements `datapup.Sender` interface. Currently there are 2 senders available:

* `datapup.Lambda` - will print your metrics according to [using-cloudwatch-logs](https://docs.datadoghq.com/integrations/amazon_lambda/#using-cloudwatch-logs)
* `datapup.StatsD` - regular DataDog metrics [client](https://github.com/DataDog/datadog-go)

If not passed, will use default DataDod's `statsd.Client` with default address `localhost:8126` and
["asynchronous" behavior](https://github.com/DataDog/datadog-go#blocking-vs-asynchronous-behavior) (with `statsd.WithAsyncUDS()` option).

Other configurations:

* `WithTag` - add default tag for all metrics, could be passed multiple times to set multiple tags
* `WithEnvironment` - add default `env` tag to all your metrics
* `WithRate` - The sampling rate in [0,1]. For example 0.5 means that half the calls will result in a metric being sent to DataDog. Rate will be 1 if not passed (all metrics will be sent)

#### Client creation example

AWS Lambda:

```go
dd := datapup.New(
    "my.awesome.service",
    datapup.WithSender(datapup.NewLambda()),
    datapup.WithEnvironment("staging"),
    datapup.WithTag("key1", "val1"),
    datapup.WithTag("key2", "val2"),
)
```

StatsD:

```go
dd := datapup.New(
    "my.awesome.service",
    datapup.WithSender(datapup.NewStatsD("localhost:8126")),
    datapup.WithEnvironment("prod"),
    datapup.WithRate(0.1),
)
```

### Metrics

Now, after client creation, you can create a DataDog metric to use:

```go
successMetric := dd.NewMetric("success", datapup.Tag("key3", "val3"))
successMetric.Incr() // shortcut for successMetric.Count(1)
// AWS Lambda client will print
// MONITORING|12345|1|count|my.awesome.service.success|#env:staging,key1:val1,key2:val2,key3:val3
```

`datapup.Tag` will format tag into `key:value` string.

Available metrics:

* Count
* Gauge
* Histogram
* Check
