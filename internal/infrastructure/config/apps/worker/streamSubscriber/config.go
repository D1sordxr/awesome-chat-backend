package streamSubscriber

type Config struct {
	ClientAddress string `yaml:"client_address"`
	Password      string `yaml:"password"`
	// TODO: stream name - group name - consumer id
}

func (c *Config) GetClientAddress() string {
	return c.ClientAddress
}

func (c *Config) GetPassword() string {
	return c.Password
}
