package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ut0mt8/goChecker/check"
	"github.com/ut0mt8/goChecker/check/check_http"
	"github.com/ut0mt8/goChecker/check/check_tcp"
)

// GetConfig ...
func GetConfig(cfgFile string) ([]check.Checker, error) {
	var checks check.Checks

	_, err := toml.DecodeFile(cfgFile, &checks)
	if err != nil {
		return nil, err
	}

	res := []check.Checker{}
	for _, c := range checks {
		err = c.Validate()
		if err != nil {
			return nil, err
		}
		switch c.Type {
		case "http":
			res = append(res, &check_http.CheckHTTP{Check: c})
		case "tcp":
			res = append(res, &check_tcp.CheckTCP{Check: c})
		}
	}

	return res, nil
}
