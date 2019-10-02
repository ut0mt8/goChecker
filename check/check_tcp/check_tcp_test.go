package check_tcp

import (
	"net"
	"testing"
	"time"

	"github.com/ut0mt8/goChecker/check"
)

func TestCheckKo(t *testing.T) {
	c := check.Check{
		Name:     "localhost:ko",
		Target:   "127.0.0.1:0",
		Interval: 0,
		Timeout:  300,
	}

	cr := make(chan check.CheckResponse, 1)
	go func() {
		(&CheckTCP{Check: c}).Run(cr)
	}()

	select {
	case r := <-cr:
		if r.IsUp != 0 {
			t.Errorf("localhost:ko test should return down status. Test expected to fail but is passing")
		}
	case <-time.After(400 * time.Millisecond):
		t.Errorf("localhost:ko test should return down status immediatly. Test expected to fail but is timeouting")
	}
	close(cr)
}

func TestCheckOk(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	c := check.Check{
		Name:     "localhost:ok",
		Target:   ln.Addr().String(),
		Interval: 0,
		Timeout:  300,
	}
	go func() {
		defer ln.Close()
		ln.Accept()
	}()

	cr := make(chan check.CheckResponse, 1)
	go func() {
		(&CheckTCP{Check: c}).Run(cr)
	}()

	select {
	case r := <-cr:
		if r.IsUp != 1 {
			t.Errorf("localhost:ok test should return up status. Test expected to succeed but is failing")
		}
	case <-time.After(400 * time.Millisecond):
		t.Errorf("localhost:ok test should return up status immediatly. Test expected to succeed but is timeouting")
	}
	close(cr)
}
