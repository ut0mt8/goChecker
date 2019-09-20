package check_http

import (
	"github.com/ut0mt8/goChecker/checker/check"
	"io/ioutil"
	"net/http"
	"time"
)

func RunCheck(c check.Check, cr chan check.CheckResponse) {
	var isUp int
	var status int
	var duration time.Duration

	start := time.Now().UnixNano()
	client := http.Client{Timeout: time.Duration(c.Timeout) * time.Millisecond}

	resp, err := client.Get(c.Target)
	if err != nil {
		isUp = 0
		status = -1
		duration = 0
	} else {
		defer resp.Body.Close()
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			isUp = 0
			status = 0
			duration = 0
		} else if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
			isUp = 1
			status = resp.StatusCode
			duration = time.Duration(time.Now().UnixNano() - start)
		} else {
			isUp = 0
			status = resp.StatusCode
			duration = time.Duration(time.Now().UnixNano() - start)
		}
	}

	cr <- check.CheckResponse{
		IsUp:     isUp,
		Status:   status,
		Duration: duration,
	}
}
