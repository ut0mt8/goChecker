package check_http

import (
	"github.com/jarcoal/httpmock"
	"github.com/ut0mt8/goChecker/check"
	"testing"
	"time"
)

func TestCheckNxDomain(t *testing.T) {
	testName := "nxdomain"
	c := check.Check{
		Name:     testName,
		Target:   "http://nxdomain.net/",
		Interval: 0,
		Timeout:  100,
	}
	cr := make(chan check.CheckResponse, 1)
	go Run(c, cr)
	select {
	case r := <-cr:
		if r.IsUp != 0 {
			t.Errorf("%s test should return down status. Test expected to fail but is passing", testName)
		}
	case <-time.After(1000 * time.Millisecond):
		t.Errorf("%s test should return down status immediatly. Test expected to fail but is timeouting", testName)
	}
	close(cr)
}

func TestCheckNotFound(t *testing.T) {
	testName := "404-notfound"
	c := check.Check{
		Name:     testName,
		Target:   "https://test.test/404",
		Interval: 0,
		Timeout:  300,
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://test.test/404", httpmock.NewStringResponder(404, "Not found"))

	cr := make(chan check.CheckResponse, 1)
	go Run(c, cr)

	select {
	case r := <-cr:
		if r.IsUp != 0 || r.Status != "bad status code" {
			t.Errorf("%s test should return down status. Test expected to fail but is passing", testName)
		}
	case <-time.After(400 * time.Millisecond):
		t.Errorf("%s test should return down status immediatly. Test expected to fail but is timeouting", testName)
	}
	close(cr)
}

func TestCheckOk(t *testing.T) {
	testName := "200-ok"
	c := check.Check{
		Name:     testName,
		Target:   "https://test.test/200",
		Interval: 0,
		Timeout:  300,
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://test.test/200", httpmock.NewStringResponder(200, "It works"))

	cr := make(chan check.CheckResponse, 1)
	go Run(c, cr)

	select {
	case r := <-cr:
		if r.IsUp != 1 || r.Status != "200" {
			t.Errorf("%s test should return up status. Test expected to succeed but is failing", testName)
		}
	case <-time.After(400 * time.Millisecond):
		t.Errorf("%s test should return up status immediatly. Test expected to succeed but is timeouting", testName)
	}
	close(cr)
}

func TestCheckBodyMatchKo(t *testing.T) {
	testName := "bodymatch-ko"
	c := check.Check{
		Name:      testName,
		Target:    "https://test.test/match",
		Interval:  0,
		Timeout:   300,
		BodyMatch: "nope",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://test.test/match", httpmock.NewStringResponder(200, "bad match"))

	cr := make(chan check.CheckResponse, 1)
	go Run(c, cr)

	select {
	case r := <-cr:
		if r.IsUp != 0 || r.Status != "body not match regexp" {
			t.Errorf("%s test should return down status. Test expected to fail but is passing", testName)
		}
	case <-time.After(400 * time.Millisecond):
		t.Errorf("%s test should return down status immediatly. Test expected to fail but is timeouting", testName)
	}
	close(cr)
}

func TestCheckBodyMatchOk(t *testing.T) {
	testName := "bodymatch-ok"
	c := check.Check{
		Name:      testName,
		Target:    "https://test.test/match",
		Interval:  0,
		Timeout:   300,
		BodyMatch: "match",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://test.test/match", httpmock.NewStringResponder(200, "good match"))

	cr := make(chan check.CheckResponse, 1)
	go Run(c, cr)

	select {
	case r := <-cr:
		if r.IsUp != 1 || r.Status != "200" {
			t.Errorf("%s test should return up status. Test expected to succeed but is failing", testName)
		}
	case <-time.After(400 * time.Millisecond):
		t.Errorf("%s test should return up status immediatly. Test expected to succeed but is timeouting", testName)
	}
	close(cr)
}

func TestCheckStatusMatchKo(t *testing.T) {
	testName := "statusmatch-ko"
	c := check.Check{
		Name:        testName,
		Target:      "https://test.test/status",
		Interval:    0,
		Timeout:     300,
		StatusMatch: "203",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://test.test/status", httpmock.NewStringResponder(200, "bad status"))

	cr := make(chan check.CheckResponse, 1)
	go Run(c, cr)

	select {
	case r := <-cr:
		if r.IsUp != 0 || r.Status != "bad status code" {
			t.Errorf("%s test should return down status. Test expected to fail but is passing", testName)
		}
	case <-time.After(400 * time.Millisecond):
		t.Errorf("%s test should return down status immediatly. Test expected to fail but is timeouting", testName)
	}
	close(cr)
}

func TestCheckStatusMatchOk(t *testing.T) {
	testName := "statusmatch-ok"
	c := check.Check{
		Name:        testName,
		Target:      "https://test.test/status",
		Interval:    0,
		Timeout:     300,
		StatusMatch: "41[45]",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://test.test/status", httpmock.NewStringResponder(414, "good status"))

	cr := make(chan check.CheckResponse, 1)
	go Run(c, cr)

	select {
	case r := <-cr:
		if r.IsUp != 1 || r.Status != "414" {
			t.Errorf("%s test should return up status. Test expected to succeed but is failing %s", testName, r.Status)
		}
	case <-time.After(400 * time.Millisecond):
		t.Errorf("%s test should return up status immediatly. Test expected to succeed but is timeouting", testName)
	}
	close(cr)
}
