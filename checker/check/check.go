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
	Type     string
	Run      func(Check, chan CheckResponse)
}

type Checks struct {
	Check []Check
}

type CheckResponse struct {
	IsUp     int
	Status   string
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
	if c.Type == "" {
		return errors.New("type of the probe should be defined")
	}
	if c.Type != "http" && c.Type != "tcp" {
		return errors.New("type of the probe is uncorrect")
	}
	return nil
}
