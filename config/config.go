package config

import (
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/creasty/defaults"
)

// Config file definition
type Config struct {
	Listener Listener  `yaml:"listener"`
	Logger   Logger    `yaml:"logger"`
	Services []Service `yaml:"services"`
	Autocert Autocert  `yaml:"autocert"`
}

// Autocert configuration
type Autocert struct {
	Email        string `yaml:"email"`
	DirectoryURL string `yaml:"directoryURL" default:"https://acme-v01.api.letsencrypt.org/directory"`
	Cache        struct {
		Backend        string            `yaml:"backend"`
		BackendOptions map[string]string `yaml:"backendOptions"`
	} `yaml:"cache"`
}

// Listener section of the configuration
type Listener struct {
	Backend struct {
		DualStack             bool          `yaml:"dualStack" default:"true"`
		Timeout               time.Duration `yaml:"timeout" default:"10s"`
		KeepAlive             time.Duration `yaml:"keepAlive" default:"30s"`
		ExpectContinueTimeout time.Duration `yaml:"expectContinueTimeout" default:"5s"`
		IdleConnTimeout       time.Duration `yaml:"idleConnTimeout" default:"10s"`
		MaxIdleConns          int           `yaml:"maxIdleConns" default:"100"`
		ResponseHeaderTimeout time.Duration `yaml:"responseHeaderTimeout" default:"10s"`
		TLSHandshakeTimeout   time.Duration `yaml:"tlsHandshakeTimeout" default:"3s"`
	} `yaml:"backend"`
	DebugAddr string `yaml:"debugAddr" default:"8081"`
	Frontend  struct {
		IdleTimeout       time.Duration `yaml:"idleTimeout" default:"5s"`
		ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout" default:"3s"`
		ReadTimeout       time.Duration `yaml:"readTimeout" default:"10s"`
		WriteTimeout      time.Duration `yaml:"writeTimeout" default:"10s"`
	} `yaml:"frontend"`
	HTTPAddr    string                   `yaml:"httpAddr" default:":80"`
	HTTPSAddr   string                   `yaml:"httpsAddr" default:":443"`
	Middlewares []map[string]interface{} `yaml:"middlewares"`
}

// Logger section of the configuration
type Logger struct {
	Formatter string `yaml:"formatter" default:"text"`
	Level     string `yaml:"level" default:"debug"`
}

// Service section of the configuration
type Service struct {
	Frontend struct {
		FQDN                []string          `yaml:"fqdn"`
		HTTPHandler         string            `yaml:"httpHandler"`
		ResponseHTTPHeaders map[string]string `yaml:"responseHTTPHeaders"`
	} `yaml:"frontend"`
	Backend struct {
		URL                string            `yaml:"url"`
		RequestHTTPHeaders map[string]string `yaml:"requestHTTPHeaders" default:"nil"`
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
	if err := defaults.Set(&config); err != nil {
		return nil, err
	}

	err := yaml.UnmarshalStrict(spec, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
