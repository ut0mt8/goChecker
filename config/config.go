package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ut0mt8/goChecker/checker/check"
)

func GetConfig(cfgFile string) (check.Checks, error) {
	var checks check.Checks

	_, err := toml.DecodeFile(cfgFile, &checks)
	if err != nil {
		return checks, err
	}

	for _, c := range checks.Check {
		err = c.CheckConfig()
		if err != nil {
			return checks, err
		}
	}

	return checks, nil
}
