package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ut0mt8/goChecker/checker/check"
	"github.com/ut0mt8/goChecker/checker/check/check_http"
	"github.com/ut0mt8/goChecker/checker/check/check_tcp"
)

func GetConfig(cfgFile string) (check.Checks, error) {
	var checks check.Checks

	_, err := toml.DecodeFile(cfgFile, &checks)
	if err != nil {
		return checks, err
	}

	for i, c := range checks.Check {
		err = check.CheckConfig(c)
		if err != nil {
			return checks, err
		}
		switch c.Type {
		case "http":
			checks.Check[i].Run = check_http.Run
		case "tcp":
			checks.Check[i].Run = check_tcp.Run
		}
	}

	return checks, nil
}
