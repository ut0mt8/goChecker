package check_http

import (
	"github.com/jarcoal/httpmock"
	"github.com/ut0mt8/goChecker/check"
	"testing"
	"time"
)

func TestCheckNxDomain(t *testing.T) {
	c := check.Check{
		Name:     "nxdomain",
		Target:   "http://nxdomain.net/",
		Interval: 0,
		Timeout:  100,
	}
	cr := make(chan check.CheckResponse, 1)
	go Run(c, cr)
	select {
	case r := <-cr:
		if r.IsUp != 0 {
			t.Errorf("nxdomain test should return down status. Test expected to fail but is passing")
		}
	case <-time.After(1000 * time.Millisecond):
		t.Errorf("nxdomain test should return down status immediatly. Test expected to fail but is timeouting")
	}
	close(cr)
}

func TestCheckNotFound(t *testing.T) {
	c := check.Check{
		Name:     "404-notfound",
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
		if r.IsUp != 0 || r.Status != "404" {
			t.Errorf("404-notfound test should return down status. Test expected to fail but is passing")
		}
	case <-time.After(400 * time.Millisecond):
		t.Errorf("404-notfound test should return down status immediatly. Test expected to fail but is timeouting")
	}
	close(cr)
}

func TestCheckOk(t *testing.T) {
	c := check.Check{
		Name:     "200-ok",
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
			t.Errorf("200-ok test should return up status. Test expected to succeed but is failing")
		}
	case <-time.After(400 * time.Millisecond):
		t.Errorf("200-ok test should return up status immediatly. Test expected to succeed but is timeouting")
	}
	close(cr)
}
