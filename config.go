package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	fhu "github.com/valyala/fasthttp/fasthttputil"
	"gopkg.in/yaml.v2"
)

type config struct {
	Listen      string
	ListenPprof string `yaml:"listen_pprof"`

	Target     string
	EnableIPv6 bool `yaml:"enable_ipv6"`

	LogLevel          string `yaml:"log_level"`
	Timeout           time.Duration
	TimeoutShutdown   time.Duration `yaml:"timeout_shutdown"`
	Concurrency       int
	Metadata          bool
	LogResponseErrors bool          `yaml:"log_response_errors"`
	MaxConnDuration   time.Duration `yaml:"max_connection_duration"`

	Auth struct {
		Egress struct {
			Username string
			Password string
		}
	}

	Tenant struct {
		Label       string
		Prefix      string `yaml:"prefix"`
		LabelRemove bool   `yaml:"label_remove"`
		Header      string
		Default     string
		AcceptAll   bool `yaml:"accept_all"`
	}

	pipeIn  *fhu.InmemoryListener
	pipeOut *fhu.InmemoryListener
}

func configParse(b []byte) (*config, error) {
	cfg := &config{}
	if err := yaml.UnmarshalStrict(b, cfg); err != nil {
		return nil, errors.Wrap(err, "Unable to parse config")
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	if cfg.Concurrency == 0 {
		cfg.Concurrency = 512
	}

	if cfg.Tenant.Header == "" {
		cfg.Tenant.Header = "X-Scope-OrgID"
	}

	if cfg.Tenant.Label == "" {
		cfg.Tenant.Label = "__tenant__"
	}

	if cfg.Auth.Egress.Username != "" {
		if cfg.Auth.Egress.Password == "" {
			return nil, fmt.Errorf("egress auth user specified, but the password is not")
		}
	}

	return cfg, nil
}

func configLoad(file string) (*config, error) {
	y, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read config")
	}

	return configParse(y)
}
