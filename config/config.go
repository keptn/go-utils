package config

import (
	// allow to load config from file
	_ "embed"
	"gopkg.in/yaml.v3"
	"sync"
)

//go:embed config.yaml
var yamlConfig []byte
var cfg KeptnGoUtilsConfig

var doOnce = sync.Once{}

// KeptnGoUtilsConfig contains config for the keptn go-utils
type KeptnGoUtilsConfig struct {
	ShKeptnSpecVersion string `yaml:"shkeptnspecversion"`
}

// GetKeptnGoUtilsConfig take the config.yaml file and reads it into the KeptnGoUtilsConfig struct
func GetKeptnGoUtilsConfig() KeptnGoUtilsConfig {
	doOnce.Do(func() {
		err := yaml.Unmarshal(yamlConfig, &cfg)
		if err != nil {
			cfg = KeptnGoUtilsConfig{}
		}
	})
	return cfg
}
