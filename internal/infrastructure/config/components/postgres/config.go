package postgres

import "fmt"

type Config struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port     int    `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	Database string `yaml:"database" env:"POSTGRES_DB" env-required:"true"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		c.Host, c.User, c.Password, c.Database, c.Port,
	)
}
