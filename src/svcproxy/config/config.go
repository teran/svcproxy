package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config file definition
type Config struct {
	Listener Listener  `yaml:"listener"`
	Services []Service `yaml:"services"`
}

// Listener section of the configuration
type Listener struct {
	HTTPAddr  string `yaml:"httpaddr" default:":80"`
	HTTPSAddr string `yaml:"httpsaddr" default:":443"`
}

// Service section of the configuration
type Service struct {
	Frontend struct {
		FQDN string `yaml:"fqdn"`
	} `yaml:"frontend"`
	Backend struct {
		URL string `yaml:"url"`
	} `yaml:"backend"`
}

// Parse reads YAML configuration file and returns Config
func Parse(path string) (*Config, error) {
	spec, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.UnmarshalStrict(spec, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
