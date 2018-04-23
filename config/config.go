package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config file definition
type Config struct {
	Listener Listener  `yaml:"listener"`
	Services []Service `yaml:"services"`
	Autocert Autocert  `yaml:"autocert"`
}

// Autocert configuration
type Autocert struct {
	Cache struct {
		Backend        string            `yaml:"backend"`
		BackendOptions map[string]string `yaml:"backendOptions"`
	} `yaml:"cache"`
}

// Listener section of the configuration
type Listener struct {
	DebugAddr   string              `yaml:"debugAddr" default:"8081"`
	HTTPAddr    string              `yaml:"httpAddr" default:":80"`
	HTTPSAddr   string              `yaml:"httpsAddr" default:":443"`
	Middlewares []map[string]string `yaml:"middlewares"`
}

// Service section of the configuration
type Service struct {
	Frontend struct {
		FQDN                []string          `yaml:"fqdn"`
		HTTPHandler         string            `yaml:"httpHandler"`
		ResponseHTTPHeaders map[string]string `yaml:"responseHTTPHeaders"`
	} `yaml:"frontend"`
	Backend struct {
		URL string `yaml:"url"`
	} `yaml:"backend"`
	Authentication struct {
		Method  string            `yaml:"method"`
		Options map[string]string `yaml:"options"`
	} `yaml:"authentication"`
}

// Load reads YAML configuration file and returns Config
func Load(path string) (*Config, error) {
	spec, err := read(path)
	if err != nil {
		return nil, err
	}

	return parse(spec)
}

func read(path string) ([]byte, error) {
	spec, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

func parse(spec []byte) (*Config, error) {
	var config Config
	err := yaml.UnmarshalStrict(spec, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
