package main

import (
	"gopkg.in/yaml.v2"
	"os"
)

type config struct {
	ApiUrl        string `yaml:"apiUrl"`
	Token         string `yaml:"token"`
	MaxConcurrent int    `yaml:"maxConcurrent"`
}

func (c config) parseConfigFile(filename string) (config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(data, &c)
	return c, err
}
