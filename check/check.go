package check

import (
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Checker interface
type Checker interface {
	Run(chan CheckResponse)
	Start()
}

// Check ...
type Check struct {
	Name     string
	Interval int
	Target   string
	Timeout  int
	Type     string
	Run      func(Check, chan CheckResponse)
}

// Checks ...
type Checks []Check

// CheckResponse ...
type CheckResponse struct {
	IsUp     int
	Status   string
	Duration time.Duration
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

// SetUpMetricsSuccess ...
func SetUpMetricsSuccess(name string, isUp float64, duration float64) {
	isupMetric.WithLabelValues(name).Set(isUp)
	successMetric.WithLabelValues(name).Inc()
	durationMetric.WithLabelValues(name).Set(duration)
}

// SetUpMetricsFailed ...
func SetUpMetricsFailed(name string, isUp int, duration float64) {
	isupMetric.WithLabelValues(name).Set(0)
	failedMetric.WithLabelValues(name).Inc()
	durationMetric.WithLabelValues(name).Set(duration)
}

// Validate ...
func (c Check) Validate() error {
	if c.Name == "" {
		return errors.New("name of the probe should be defined")
	}
	if c.Interval == 0 {
		return errors.New("interval of the probe should be defined")
	}
	if c.Target == "" {
		return errors.New("target of the probe should be defined")
	}
	if c.Timeout == 0 {
		return errors.New("timeout of the probe should be defined")
	}
	if c.Type == "" {
		return errors.New("type of the probe should be defined")
	}
	if c.Type != "http" && c.Type != "tcp" {
		return errors.New("type of the probe is uncorrect")
	}
	return nil
}

// Start ...
/* func (c Check) Start() {
	log.Printf("check thread %s started\n", c.Name)
	ticker := time.NewTicker(time.Duration(c.Interval) * time.Millisecond)

	for range ticker.C {
		cr := make(chan CheckResponse, 1)
		go c.Run(c, cr)

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

	log.Printf("check thread %s ended\n", c.Name)
} */
