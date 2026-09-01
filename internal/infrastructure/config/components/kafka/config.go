package kafka

type Config struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS" env-separator:"," env-required:"true"`
	Topic   string   `yaml:"topic" env:"KAFKA_TOPIC" env-required:"true"`
}
