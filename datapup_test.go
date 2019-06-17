package datapup

import (
	"reflect"
	"testing"
)

func TestNew_default(t *testing.T) {
	c := New("testservice")

	if c.namespace != "testservice" {
		t.Errorf("Want namespace testservice, get %s", c.namespace)
	}

	if !reflect.DeepEqual(c.tags, []string(nil)) {
		t.Errorf("Want tags []string(nil), get %v", c.tags)
	}

	rate := float64(1)
	if !reflect.DeepEqual(c.rate, &rate) {
		t.Errorf("Want rate 1, get %v", c.rate)
	}

	if want := reflect.TypeOf(&StatsD{}); want != reflect.TypeOf(c.sender) {
		t.Errorf("Want seder to be of type %v , get %v", want, reflect.TypeOf(c.sender))
	}
}

func TestNew_withOptions(t *testing.T) {
	c := New(
		"testservice",
		WithEnvironment("test"),
		WithRate(0),
		WithTag("a", "a"),
		WithTag("b", "b"),
		WithSender(NewLambda()),
	)

	if c.namespace != "testservice" {
		t.Errorf("Want namespace testservice, get %s", c.namespace)
	}

	if !reflect.DeepEqual(c.tags, []string{"env:test", "a:a", "b:b"}) {
		t.Errorf("Want tags []string(nil), get %v", c.tags)
	}

	rate := 0.0
	if !reflect.DeepEqual(c.rate, &rate) {
		t.Errorf("Want rate 0.0, get %v", c.rate)
	}

	if want := reflect.TypeOf(&Lambda{}); want != reflect.TypeOf(c.sender) {
		t.Errorf("Want seder to be of type %v , get %v", want, reflect.TypeOf(c.sender))
	}
}
