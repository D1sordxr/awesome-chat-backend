package kafka

type Config struct {
	Brokers []string `yaml:"brokers" binding:"required"`
	Topic   string   `yaml:"topic" binding:"required"`
}
