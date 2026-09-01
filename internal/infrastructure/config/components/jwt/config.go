package jwt

type Config struct {
	SecretKey string `yaml:"secret_key" env:"JWT_SECRET_KEY" env-required:"true"`
}
