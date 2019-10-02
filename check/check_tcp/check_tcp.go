package check_tcp

import (
	"log"
	"net"
	"time"

	"github.com/ut0mt8/goChecker/check"
)

// CheckTCP implemetation
type CheckTCP struct {
	check.Check
}

// Run implemetation
func (c *CheckTCP) Run(cr chan check.CheckResponse) {
	var isUp int
	var status string
	var duration time.Duration

	start := time.Now() //time.Now().UnixNano()
	_, err := net.DialTimeout("tcp", c.Target, time.Duration(c.Timeout)*time.Millisecond)

	if err != nil {
		isUp = 0
		status = "connection failed"
		duration = 0
	} else {
		isUp = 1
		status = "connection success"
		duration = time.Since(start) //time.Duration(time.Now().UnixNano() - start)
	}

	cr <- check.CheckResponse{
		IsUp:     isUp,
		Status:   status,
		Duration: duration,
	}
}

// Start ...
func (c *CheckTCP) Start() {
	log.Printf("check thread %s started\n", c.Name)
	ticker := time.NewTicker(time.Duration(c.Interval) * time.Millisecond)

	for range ticker.C {
		cr := make(chan check.CheckResponse, 1)
		go c.Run(cr)

		select {
		case r := <-cr:
			check.SetUpMetricsSuccess(c.Name, float64(r.IsUp), float64(r.Duration))
			log.Printf("check %s %s %s %v\n", c.Name, c.Target, r.Status, r.Duration)
		case <-time.After(time.Duration(c.Timeout) * time.Millisecond):
			check.SetUpMetricsFailed(c.Name, 0, float64(c.Timeout))
			log.Printf("check %s %s timeout\n", c.Name, c.Target)
		}
	}
	log.Printf("check thread %s ended\n", c.Name)
}
