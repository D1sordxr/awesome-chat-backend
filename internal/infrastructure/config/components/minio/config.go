package minio

type Config struct {
	Endpoint  string `yaml:"endpoint" env:"MINIO_ENDPOINT" env-required:"true"`
	AccessKey string `yaml:"access_key" env:"MINIO_ACCESS_KEY" env-required:"true"`
	SecretKey string `yaml:"secret_key" env:"MINIO_SECRET_KEY" env-required:"true"`
	UseSSL    bool   `yaml:"use_ssl" env:"MINIO_USE_SSL"`
}
