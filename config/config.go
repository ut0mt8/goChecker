package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ut0mt8/goChecker/checker"
)

func GetConfig(cfgFile string) (checker.Checks, error) {
	var checks checker.Checks

	_, err := toml.DecodeFile(cfgFile, &checks)
	if err != nil {
		return checks, err
	}

	for _, c := range checks.Check {
		err = checker.CheckConfig(c)
		if err != nil {
			return checks, err
		}
	}

	return checks, nil
}
