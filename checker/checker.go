package checker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ut0mt8/goChecker/checker/check"
	"github.com/ut0mt8/goChecker/checker/check/check_http"
	"github.com/ut0mt8/goChecker/checker/check/check_tcp"
	"time"
)

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

func StartChecker(c check.Check) {
	log.Printf("checker %s started\n", c.Name)
	ticker := time.NewTicker(time.Duration(c.Interval) * time.Millisecond)

	for range ticker.C {
		cr := make(chan check.CheckResponse, 1)

		switch c.Type {
		case "http":
			go check_http.RunCheck(c, cr)
		case "tcp":
			go check_tcp.RunCheck(c, cr)
		}

		select {
		case r := <-cr:
			isupMetric.WithLabelValues(c.Name).Set(float64(r.IsUp))
			successMetric.WithLabelValues(c.Name).Inc()
			durationMetric.WithLabelValues(c.Name).Set(float64(r.Duration))
			log.Printf("check %s %s %s %v\n", c.Name, c.Target, r.Status, r.Duration)
		case <-time.After(time.Duration(c.Timeout) * time.Millisecond):
			isupMetric.WithLabelValues(c.Name).Set(0)
			failedMetric.WithLabelValues(c.Name).Inc()
			durationMetric.WithLabelValues(c.Name).Set(float64(c.Timeout))
			log.Printf("check %s %s timeout\n", c.Name, c.Target)
		}
	}

	log.Printf("checker %s ended\n", c.Name)
}
