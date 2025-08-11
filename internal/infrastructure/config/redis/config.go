package redis

type Config struct {
	ClientAddress string `yaml:"client_address"`
	Password      string `yaml:"password"`
	Channel       string `yaml:"channel"`
}

func (c *Config) GetClientAddress() string {
	return c.ClientAddress
}

func (c *Config) GetPassword() string {
	return c.Password
}

func (c *Config) GetChannel() string {
	return c.Channel
}
