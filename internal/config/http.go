package config

import (
	"net"
	"os"
	"time"
)

const (
	httpHostEnvName        = "HOST"
	httpPortEnvName        = "PORT"
	httpTimeoutEnvName     = "TIMEOUT"
	httpIdleTimeoutEnvName = "IDLE_TIMEOUT"
)

type HTTPConfig interface {
	Address() string
	Timeout() time.Duration
	IdleTimeout() time.Duration
}

type httpConfig struct {
	host        string `env:"HOST" env-default:"0.0.0.0"`
	port        string `env:"PORT" env-default:"8083"`
	timeout     string `env:"timeout" env-default:"5s"`
	idleTimeout string `env:"idle_timeout" env-default:"60s"`
}

func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		host = "0.0.0.0"
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		port = "8083"
	}

	timeout := os.Getenv(httpTimeoutEnvName)
	if len(timeout) == 0 {
		timeout = "5s"
	}
	idleTimeout := os.Getenv(httpIdleTimeoutEnvName)
	if len(idleTimeout) == 0 {
		idleTimeout = "60s"
	}

	return &httpConfig{
		host:        host,
		port:        port,
		timeout:     timeout,
		idleTimeout: idleTimeout,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *httpConfig) Timeout() time.Duration {
	duration, _ := time.ParseDuration(cfg.timeout)
	return duration
}
func (cfg *httpConfig) IdleTimeout() time.Duration {
	idleTimeout, _ := time.ParseDuration(cfg.idleTimeout)
	return idleTimeout
}
