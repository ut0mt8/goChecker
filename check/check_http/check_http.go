package check_http

import (
	"github.com/ut0mt8/goChecker/check"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

func Run(c check.Check, cr chan check.CheckResponse) {
	start := time.Now()
	client := http.Client{Timeout: time.Duration(c.Timeout) * time.Millisecond}

	resp, err := client.Get(c.Target)
	if err != nil {
		c.SendResponse(cr, 0, "connection failed", 0)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.SendResponse(cr, 0, "body read failed", 0)
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		c.SendResponse(cr, 0, resp.Status, time.Since(start))
		return
	}

	if c.BodyMatch != "" {
		re, err := regexp.Compile(c.BodyMatch)
		if err != nil {
			c.SendResponse(cr, 0, "bad body regexp", time.Since(start))
			return
		}
		if !re.Match(body) {
			c.SendResponse(cr, 0, "body not match regexp", time.Since(start))
			return
		}
	}

	c.SendResponse(cr, 1, resp.Status, time.Since(start))
}
