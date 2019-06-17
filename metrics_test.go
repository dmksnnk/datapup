package datapup

import (
	"testing"

	"github.com/DataDog/datadog-go/statsd"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type fakeSender struct {
	mock.Mock
}

func (fs *fakeSender) Send(name string, value interface{}, mt MetricType, rate float64, tags ...string) error {
	args := fs.Called(name, value, mt, rate, tags)
	return args.Error(0)
}

type MetricsSuite struct {
	suite.Suite
	fs     fakeSender
	metric *Metric
}

func (suite *MetricsSuite) SetupTest() {
	suite.fs = fakeSender{}
	client := New("testservice", WithSender(&suite.fs))
	suite.metric = client.NewMetric("testmetric", Tag("key", "value"))
}

func TestTag(t *testing.T) {
	if tag := Tag("key", "value"); tag != "key:value" {
		t.Errorf("Want key:value, get %s", tag)
	}
}

func (suite *MetricsSuite) TestCount() {
	suite.fs.On("Send", "testservice.testmetric", int64(2), Count, float64(1), []string{"key:value"}).Return(nil)
	err := suite.metric.Count(2)
	suite.Assert().Nil(err)
	suite.fs.AssertExpectations(suite.T())
}

func (suite *MetricsSuite) TestIncr() {
	suite.fs.On("Send", "testservice.testmetric", int64(1), Count, float64(1), []string{"key:value"}).Return(nil)
	err := suite.metric.Incr()
	suite.Assert().Nil(err)
	suite.fs.AssertExpectations(suite.T())
}

func (suite *MetricsSuite) TestDecr() {
	suite.fs.On("Send", "testservice.testmetric", int64(-1), Count, float64(1), []string{"key:value"}).Return(nil)
	err := suite.metric.Decr()
	suite.Assert().Nil(err)
	suite.fs.AssertExpectations(suite.T())
}

func (suite *MetricsSuite) TestGauge() {
	suite.fs.On("Send", "testservice.testmetric", 1.1, Gauge, float64(1), []string{"key:value"}).Return(nil)
	err := suite.metric.Gauge(1.1)
	suite.Assert().Nil(err)
	suite.fs.AssertExpectations(suite.T())
}

func (suite *MetricsSuite) TestHistogram() {
	suite.fs.On("Send", "testservice.testmetric", 1.1, Histogram, float64(1), []string{"key:value"}).Return(nil)
	err := suite.metric.Histogram(1.1)
	suite.Assert().Nil(err)
	suite.fs.AssertExpectations(suite.T())
}

func (suite *MetricsSuite) TestCheck() {
	suite.fs.On("Send", "testservice.testmetric", statsd.Critical, Check, float64(1), []string{"key:value"}).Return(nil)
	err := suite.metric.Check(statsd.Critical)
	suite.Assert().Nil(err)
	suite.fs.AssertExpectations(suite.T())
}

func TestRunMetricsSuite(t *testing.T) {
	suite.Run(t, new(MetricsSuite))
}
