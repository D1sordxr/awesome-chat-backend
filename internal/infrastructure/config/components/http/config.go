package http

import "time"

type Config struct {
	Port        string        `yaml:"port" env:"HTTP_PORT" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}
