package checker

import (
	"github.com/jarcoal/httpmock"
	"testing"
	"time"
)

func TestCheckNxDomain(t *testing.T) {
	c := Check{
		Name:     "nxdomain",
		Url:      "http://nxdomain.net/",
		Interval: 0,
		Timeout:  100,
	}
	cr := make(chan CheckResponse, 1)
	go RunCheck(c, cr)
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
	c := Check{
		Name:     "404-notfound",
		Url:      "https://test.test/404",
		Interval: 0,
		Timeout:  300,
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://test.test/404",
		httpmock.NewStringResponder(404, "Not found"))

	cr := make(chan CheckResponse, 1)
	go RunCheck(c, cr)

	select {
	case r := <-cr:
		if r.IsUp != 0 || r.StatusCode != 404 {
			t.Errorf("404-notfound test should return down status. Test expected to fail but is passing")
		}
	case <-time.After(400 * time.Millisecond):
		t.Errorf("404-notfound test should return down status immediatly. Test expected to fail but is timeouting")
	}
	close(cr)
}

func TestCheckOk(t *testing.T) {
	c := Check{
		Name:     "200-ok",
		Url:      "https://test.test/200",
		Interval: 0,
		Timeout:  300,
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://test.test/200",
		httpmock.NewStringResponder(200, "It works"))

	cr := make(chan CheckResponse, 1)
	go RunCheck(c, cr)

	select {
	case r := <-cr:
		if r.IsUp != 1 || r.StatusCode != 200 {
			t.Errorf("200-ok test should return up status. Test expected to succeed but is failing")
		}
	case <-time.After(400 * time.Millisecond):
		t.Errorf("200-ok test should return up status immediatly. Test expected to succeed but is timeouting")
	}
	close(cr)
}
