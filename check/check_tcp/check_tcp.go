package check_tcp

import (
	"github.com/ut0mt8/goChecker/check"
	"net"
	"time"
)

func Run(c check.Check, cr chan check.CheckResponse) {
	var isUp int
	var status string
	var duration time.Duration

	start := time.Now()
	_, err := net.DialTimeout("tcp", c.Target, time.Duration(c.Timeout)*time.Millisecond)

	if err != nil {
		isUp = 0
		status = "connection failed"
		duration = 0
	} else {
		isUp = 1
		status = "connection success"
		duration = time.Since(start)
	}

	cr <- check.CheckResponse{
		IsUp:     isUp,
		Status:   status,
		Duration: duration,
	}
}
