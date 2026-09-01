package loader

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const PathEnv = "CONFIG_PATH"

func Read(defaultPath string, cfg any) error {
	path := os.Getenv(PathEnv)
	if path == "" {
		path = defaultPath
	}

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("config file %q: %w", path, err)
	}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return fmt.Errorf("read config %q: %w", path, err)
	}

	return nil
}
