package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"awesome-chat/internal/infrastructure/config/components/postgres"
	"awesome-chat/internal/infrastructure/config/loader"
)

type testConfig struct {
	Storage postgres.Config `yaml:"storage"`
}

const validYAML = `storage:
  host: "db"
  user: "u"
  password: "p"
  port: 5432
  database: "chat"
`

func writeConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "cfg.yaml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestReadLoadsDefaultPath(t *testing.T) {
	t.Setenv(loader.PathEnv, "")

	var cfg testConfig
	if err := loader.Read(writeConfig(t, validYAML), &cfg); err != nil {
		t.Fatalf("Read() = %v, want nil", err)
	}
	if cfg.Storage.Host != "db" || cfg.Storage.Port != 5432 {
		t.Errorf("cfg = %+v, want host=db port=5432", cfg.Storage)
	}
}

func TestReadPrefersConfigPathEnv(t *testing.T) {
	override := writeConfig(t, `storage:
  host: "from-env"
  user: "u"
  password: "p"
  port: 1
  database: "chat"
`)
	t.Setenv(loader.PathEnv, override)

	var cfg testConfig
	if err := loader.Read(writeConfig(t, validYAML), &cfg); err != nil {
		t.Fatalf("Read() = %v, want nil", err)
	}
	if cfg.Storage.Host != "from-env" {
		t.Errorf("host = %q, want from-env", cfg.Storage.Host)
	}
}

func TestReadReportsMissingFile(t *testing.T) {
	t.Setenv(loader.PathEnv, "")

	var cfg testConfig
	err := loader.Read(filepath.Join(t.TempDir(), "nope.yaml"), &cfg)
	if err == nil {
		t.Fatal("Read() = nil, want error for missing file")
	}
	if !os.IsNotExist(err) && err.Error() == "" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestReadRejectsMissingRequiredField(t *testing.T) {
	t.Setenv(loader.PathEnv, "")

	path := writeConfig(t, `storage:
  host: "db"
  user: "u"
  port: 5432
  database: "chat"
`)

	var cfg testConfig
	if err := loader.Read(path, &cfg); err == nil {
		t.Fatal("Read() = nil, want error for missing required field")
	}
}

func TestEnvOverridesFileValue(t *testing.T) {
	t.Setenv(loader.PathEnv, "")
	t.Setenv("POSTGRES_HOST", "overridden")

	var cfg testConfig
	if err := loader.Read(writeConfig(t, validYAML), &cfg); err != nil {
		t.Fatalf("Read() = %v, want nil", err)
	}
	if cfg.Storage.Host != "overridden" {
		t.Errorf("host = %q, want overridden", cfg.Storage.Host)
	}
}
