package grpc

import "time"

type Config struct {
	Port              string        `yaml:"port" env-required:"true"`
	Timeout           time.Duration `yaml:"timeout" env-default:"5s"`             // time.Second
	Time              time.Duration `yaml:"time" env-default:"15m"`               // time.Minute
	MaxConnectionIdle time.Duration `yaml:"max_connection_idle" env-default:"5m"` // time.Minute
	MaxConnectionAge  time.Duration `yaml:"max_connection_age" env-default:"5m"`  // time.Minute
}
