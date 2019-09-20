package check

import (
	"errors"
	"time"
)

type Check struct {
	Name     string
	Interval int
	Target   string
	Timeout  int
}

type Checks struct {
	Check []Check
}

type CheckResponse struct {
	IsUp     int
	Status   int
	Duration time.Duration
}

func CheckConfig(c Check) error {
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
	return nil
}
