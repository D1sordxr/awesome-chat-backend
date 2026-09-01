package loader_test

import (
	"path/filepath"
	"testing"

	"awesome-chat/internal/infrastructure/config/apps/api"
	"awesome-chat/internal/infrastructure/config/apps/outboxprocessor"
	"awesome-chat/internal/infrastructure/config/apps/topiccreator"
	"awesome-chat/internal/infrastructure/config/apps/worker"
	"awesome-chat/internal/infrastructure/config/apps/wsserver"
	"awesome-chat/internal/infrastructure/config/loader"
)

const repoRoot = "../../../.."

func TestShippedConfigsLoad(t *testing.T) {
	profiles := []string{"prod", "local"}

	apps := map[string]func() error{
		"api": func() error {
			_, err := api.NewConfig()
			return err
		},
		"outbox-processor": func() error {
			_, err := outboxprocessor.NewConfig()
			return err
		},
		"topic-creator": func() error {
			_, err := topiccreator.NewConfig()
			return err
		},
		"worker": func() error {
			_, err := worker.NewConfig()
			return err
		},
		"ws-server": func() error {
			_, err := wsserver.NewConfig()
			return err
		},
	}

	for name, load := range apps {
		for _, profile := range profiles {
			t.Run(name+"/"+profile, func(t *testing.T) {
				path := filepath.Join(repoRoot, "configs", name, profile+".yaml")
				t.Setenv(loader.PathEnv, path)

				if err := load(); err != nil {
					t.Fatalf("loading %s: %v", path, err)
				}
			})
		}
	}
}

func TestKafkaConfigIsPopulated(t *testing.T) {
	t.Setenv(loader.PathEnv, filepath.Join(repoRoot, "configs", "outbox-processor", "prod.yaml"))

	cfg, err := outboxprocessor.NewConfig()
	if err != nil {
		t.Fatalf("NewConfig() = %v", err)
	}
	if len(cfg.MessageBroker.Brokers) == 0 {
		t.Error("MessageBroker.Brokers is empty; the yaml key does not match the struct tag")
	}
	if cfg.MessageBroker.Topic == "" {
		t.Error("MessageBroker.Topic is empty")
	}
}
