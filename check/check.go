package check

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

type Check struct {
	Name      string
	Interval  int
	Target    string
	Timeout   int
	Type      string
	Run       func(Check, chan CheckResponse)
	BodyMatch string
}

type Checks []Check

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

func (c Check) SendResponse(cr chan CheckResponse, up int, status string, duration time.Duration) {
	cr <- CheckResponse{
		IsUp:     up,
		Status:   status,
		Duration: duration,
	}
}

func (c Check) Start(done chan bool, wg *sync.WaitGroup) {
	log.Printf("> thread %s started\n", c.Name)
	ticker := time.NewTicker(time.Duration(c.Interval) * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			cr := make(chan CheckResponse, 1)
			go c.Run(c, cr)

			select {
			case r := <-cr:
				isupMetric.WithLabelValues(c.Name).Set(float64(r.IsUp))
				successMetric.WithLabelValues(c.Name).Inc()
				durationMetric.WithLabelValues(c.Name).Set(float64(r.Duration))
				log.Printf("[%s] %s : %s (%v)\n", c.Name, c.Target, r.Status, r.Duration)
			case <-time.After(time.Duration(c.Timeout) * time.Millisecond):
				isupMetric.WithLabelValues(c.Name).Set(0)
				failedMetric.WithLabelValues(c.Name).Inc()
				durationMetric.WithLabelValues(c.Name).Set(float64(c.Timeout))
				log.Printf("[%s] %s : timeout\n", c.Name, c.Target)
			}
		case <-done:
			log.Printf("> thread %s ended\n", c.Name)
			wg.Done()
			return
		}
	}

}
