package checker

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"net/http"
	"time"
)

type Check struct {
	Name     string
	Url      string
	Interval int
	Timeout  int
}

type Checks struct {
	Check []Check
}

type CheckResponse struct {
	IsUp       int
	Status     string
	StatusCode int
	Duration   time.Duration
}

var (
	isupMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "checker_up",
			Help: "Check status, 1 : up,  0 : down",
		},
		[]string{"check"},
	)
	durationMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "checker_duration",
			Help: "Check last duration time",
		},
		[]string{"check"},
	)
	successMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "checker_success",
			Help: "Number of succeeded checks",
		},
		[]string{"check"},
	)
	failedMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "checker_failed",
			Help: "Number of failed checks",
		},
		[]string{"check"},
	)
)

func init() {
	prometheus.MustRegister(isupMetric)
	prometheus.MustRegister(durationMetric)
	prometheus.MustRegister(successMetric)
	prometheus.MustRegister(failedMetric)
}

func CheckConfig(c Check) error {
	if c.Name == "" {
		return errors.New("name of the probe should be defined")
	}
	if c.Url == "" {
		return errors.New("url of the probe should be defined")
	}
	if c.Interval == 0 {
		return errors.New("interval of the probe should be defined")
	}
	if c.Timeout == 0 {
		return errors.New("duration of the probe should be defined")
	}
	return nil
}

func RunCheck(c Check, cr chan CheckResponse) {
	var isUp int
	var status string
	var statusCode int
	var duration time.Duration

	start := time.Now().UnixNano()
	client := http.Client{Timeout: time.Duration(c.Timeout) * time.Millisecond}

	resp, err := client.Get(c.Url)
	if err != nil {
		isUp = 0
		status = "cannot connect"
		duration = 0
	} else {
		defer resp.Body.Close()
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			isUp = 0
			status = "cannot read body"
			duration = 0
		} else if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
			isUp = 1
			status = resp.Status
			duration = time.Duration(time.Now().UnixNano() - start)
		} else {
			isUp = 0
			status = resp.Status
			duration = time.Duration(time.Now().UnixNano() - start)
		}
		statusCode = resp.StatusCode
	}

	cr <- CheckResponse{
		IsUp:       isUp,
		Status:     status,
		StatusCode: statusCode,
		Duration:   duration,
	}
}

func StartCheck(c Check) {
	log.Printf("checker %s started\n", c.Name)
	ticker := time.NewTicker(time.Duration(c.Interval) * time.Millisecond)

	for t := range ticker.C {
		cr := make(chan CheckResponse, 1)

		go RunCheck(c, cr)

		select {
		case r := <-cr:
			isupMetric.WithLabelValues(c.Name).Set(float64(r.IsUp))
			successMetric.WithLabelValues(c.Name).Inc()
			durationMetric.WithLabelValues(c.Name).Set(float64(r.Duration))
			log.Printf("check %s %s %d %s %v at %v \n", c.Name, c.Url, r.StatusCode, r.Status, r.Duration, t)
		case <-time.After(time.Duration(c.Timeout) * time.Millisecond):
			isupMetric.WithLabelValues(c.Name).Set(0)
			failedMetric.WithLabelValues(c.Name).Inc()
			durationMetric.WithLabelValues(c.Name).Set(float64(c.Timeout))
			log.Printf("check %s %s Timeout at %v \n", c.Name, c.Url, t)
		}
	}

	log.Printf("checker %s ended\n", c.Name)
}
