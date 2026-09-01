package wsserverapi

type Config struct {
	BroadcastURL string `yaml:"broadcast_url" env:"WS_SERVER_BROADCAST_URL" env-required:"true"`
}
