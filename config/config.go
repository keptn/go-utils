package config

import (
	_ "embed"
	"gopkg.in/yaml.v3"
	"sync"
)

//go:embed config.yaml
var yamlConfig []byte
var cfg KeptnGoUtilsConfig

var doOnce = sync.Once{}

type KeptnGoUtilsConfig struct {
	ShKeptnSpecVersion string `yaml:"shkeptnspecversion"`
}

func GetKeptnGoUtilsConfig() KeptnGoUtilsConfig {
	doOnce.Do(func() {
		err := yaml.Unmarshal(yamlConfig, &cfg)
		if err != nil {
			cfg = KeptnGoUtilsConfig{}
		}
	})
	return cfg
}
