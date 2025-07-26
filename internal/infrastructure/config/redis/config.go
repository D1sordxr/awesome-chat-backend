package redis

type Config struct {
	ClientAddress string `yaml:"client_address"`
	Channel       string `yaml:"channel"`
}

func (c *Config) GetClientAddress() string {
	return c.ClientAddress
}

func (c *Config) GetChannel() string {
	return c.Channel
}
