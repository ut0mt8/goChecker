package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ut0mt8/goChecker/checker"
)

func GetConfig(cfgFile string) (checker.Checks, error) {
	var conf checker.Checks

	_, err := toml.DecodeFile(cfgFile, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}
