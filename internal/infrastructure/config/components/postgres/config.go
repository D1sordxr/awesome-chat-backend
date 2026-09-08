package postgres

import (
	"fmt"

	"awesome-chat/internal/infrastructure/config/env"
)

type Config struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func NewConfig(l *env.Loader) Config {
	return Config{
		Host:     l.String("POSTGRES_HOST"),
		Port:     l.Int("POSTGRES_PORT"),
		Database: l.String("POSTGRES_DB"),
		User:     l.String("POSTGRES_USER"),
		Password: l.String("POSTGRES_PASSWORD"),
	}
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		c.Host, c.User, c.Password, c.Database, c.Port,
	)
}
