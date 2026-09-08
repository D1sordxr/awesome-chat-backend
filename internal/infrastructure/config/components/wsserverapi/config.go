package wsserverapi

import "awesome-chat/internal/infrastructure/config/env"

type Config struct {
	BroadcastURL string
}

func NewConfig(l *env.Loader) Config {
	return Config{
		BroadcastURL: l.String("WS_SERVER_BROADCAST_URL"),
	}
}
