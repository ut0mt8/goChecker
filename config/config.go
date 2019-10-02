package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ut0mt8/goChecker/check"
	"github.com/ut0mt8/goChecker/check/check_http"
	"github.com/ut0mt8/goChecker/check/check_tcp"
)

type Config struct {
	Checks map[string]*check.Check
}

func GetConfig(cfgFile string) (check.Checks, error) {
	var cfg Config
	var checks check.Checks

	_, err := toml.DecodeFile(cfgFile, &cfg)
	if err != nil {
		return nil, err
	}

	for k, c := range cfg.Checks {

		cfg.Checks[k].Name = k
		switch c.Type {
		case "http":
			cfg.Checks[k].Run = check_http.Run
		case "tcp":
			cfg.Checks[k].Run = check_tcp.Run
		}

		err = c.Validate()
		if err != nil {
			return nil, err
		}
		checks = append(checks, *c)
	}

	return checks, nil
}
