package check_http

import (
	"github.com/ut0mt8/goChecker/check"
	"io/ioutil"
	"net/http"
	"time"
)

func Run(c check.Check, cr chan check.CheckResponse) {
	var isUp int
	var status string
	var duration time.Duration

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
			duration = time.Since(start)
		} else {
			isUp = 0
			status = resp.Status
			duration = time.Since(start)
		}
	}

	cr <- check.CheckResponse{
		IsUp:     isUp,
		Status:   status,
		Duration: duration,
	}
}
