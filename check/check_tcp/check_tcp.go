package check_tcp

import (
	"github.com/ut0mt8/goChecker/check"
	"net"
	"time"
)

func Run(c check.Check, cr chan check.CheckResponse) {
	start := time.Now()
	_, err := net.DialTimeout("tcp", c.Target, time.Duration(c.Timeout)*time.Millisecond)
	if err != nil {
		c.SendResponse(cr, 0, "connection failed", 0)
	}

	c.SendResponse(cr, 1, "connection succeed", time.Since(start))
}
