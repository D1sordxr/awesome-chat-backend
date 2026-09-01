package redis

type Config struct {
	ClientAddress string `yaml:"client_address" env:"REDIS_ADDRESS" env-required:"true"`
	Password      string `yaml:"password" env:"REDIS_PASSWORD"`
	Channel       string `yaml:"channel" env:"REDIS_CHANNEL"`
}
