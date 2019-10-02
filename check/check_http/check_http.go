package check_http

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ut0mt8/goChecker/check"
)

// CheckHTTP implemetation
type CheckHTTP struct {
	check.Check
}

// Run implemantation
func (c *CheckHTTP) Run(cr chan check.CheckResponse) {
	var isUp int
	var status string
	var duration time.Duration

	//start := time.Now().UnixNano()
	start := time.Now()
	client := http.Client{Timeout: time.Duration(c.Timeout) * time.Millisecond}

	resp, err := client.Get(c.Target)
	if err != nil {
		isUp = 0
		status = "connection failed"
		duration = 0
	} else {
		defer resp.Body.Close()
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			isUp = 0
			status = "body read failed"
			duration = 0
		} else if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
			isUp = 1
			status = resp.Status
			duration = time.Since(start) //time.Duration(time.Now().UnixNano() - start)
		} else {
			isUp = 0
			status = resp.Status
			duration = time.Since(start) //time.Duration(time.Now().UnixNano() - start)
		}
	}

	cr <- check.CheckResponse{
		IsUp:     isUp,
		Status:   status,
		Duration: duration,
	}
}

// Start ...
func (c *CheckHTTP) Start() {
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
