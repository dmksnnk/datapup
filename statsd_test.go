package datapup

import (
	"testing"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type fakeStatsDSender struct {
	mock.Mock
}

func (fss *fakeStatsDSender) Count(name string, value int64, tags []string, rate float64) error {
	args := fss.Called(name, value, tags, rate)
	return args.Error(0)
}

func (fss *fakeStatsDSender) Gauge(name string, value float64, tags []string, rate float64) error {
	args := fss.Called(name, value, tags, rate)
	return args.Error(0)
}

func (fss *fakeStatsDSender) Histogram(name string, value float64, tags []string, rate float64) error {
	args := fss.Called(name, value, tags, rate)
	return args.Error(0)
}

func (fss *fakeStatsDSender) ServiceCheck(sc *statsd.ServiceCheck) error {
	args := fss.Called(sc)
	return args.Error(0)
}

type StatsDSuite struct {
	suite.Suite
	fss    fakeStatsDSender
	statsd *StatsD
}

func (suite *StatsDSuite) SetupTest() {
	suite.fss = fakeStatsDSender{}
	suite.statsd = &StatsD{client: &suite.fss}
}

func (suite *StatsDSuite) TestCount() {
	suite.fss.On("Count", "testmetric", int64(1), []string{"key:value"}, 0.5).Return(nil)
	err := suite.statsd.Send("testmetric", int64(1), Count, 0.5, Tag("key", "value"))
	suite.Assert().Nil(err)
	// bad type
	suite.fss.On("Count", "testmetric", int64(1), []string{"key:value"}, 0.5).Return(nil)
	err = suite.statsd.Send("testmetric", "bad param", Count, 0.5, Tag("key", "value"))
	suite.Assert().EqualError(err, "Can't convert bad param to int64")
}

func (suite *StatsDSuite) TestGauge() {
	suite.fss.On("Gauge", "testmetric", 0.3, []string{"key:value"}, 0.5).Return(nil)
	err := suite.statsd.Send("testmetric", 0.3, Gauge, 0.5, Tag("key", "value"))
	suite.Assert().Nil(err)
	// bad type
	suite.fss.On("Gauge", "testmetric", 0.3, []string{"key:value"}, 0.5).Return(nil)
	err = suite.statsd.Send("testmetric", "bad param", Gauge, 0.5, Tag("key", "value"))
	suite.Assert().EqualError(err, "Can't convert bad param to float64")
}

func (suite *StatsDSuite) TestHistogram() {
	suite.fss.On("Histogram", "testmetric", 0.3, []string{"key:value"}, 0.5).Return(nil)
	err := suite.statsd.Send("testmetric", 0.3, Histogram, 0.5, Tag("key", "value"))
	suite.Assert().Nil(err)
	// bad type
	suite.fss.On("Histogram", "testmetric", 0.3, []string{"key:value"}, 0.5).Return(nil)
	err = suite.statsd.Send("testmetric", "bad param", Histogram, 0.5, Tag("key", "value"))
	suite.Assert().EqualError(err, "Can't convert bad param to float64")
}

func (suite *StatsDSuite) TestCheck() {
	suite.fss.On(
		"ServiceCheck",
		&statsd.ServiceCheck{Name: "testmetric", Status: statsd.Critical, Tags: []string{"key:value"}},
	).Return(nil)
	err := suite.statsd.Send("testmetric", statsd.Critical, Check, 0.5, Tag("key", "value"))
	suite.Assert().Nil(err)
	// bad type
	suite.fss.On(
		"ServiceCheck",
		&statsd.ServiceCheck{Name: "testmetric", Status: statsd.Critical, Tags: []string{"key:value"}},
	).Return(nil)
	err = suite.statsd.Send("testmetric", "bad param", Check, 0.5, Tag("key", "value"))
	suite.Assert().EqualError(err, "Can't convert bad param to statsd.ServiceCheckStatus")
}

func TestRunStatsDSuite(t *testing.T) {
	suite.Run(t, new(StatsDSuite))
}
