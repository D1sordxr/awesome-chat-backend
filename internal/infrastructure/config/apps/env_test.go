package apps_test

import (
	"testing"

	"awesome-chat/internal/infrastructure/config/apps"
)

func TestParseAppEnvKnownValues(t *testing.T) {
	cases := map[string]apps.AppEnv{
		"local":       apps.Local,
		"development": apps.Development,
		"staging":     apps.Staging,
		"production":  apps.Production,
	}

	for raw, want := range cases {
		got, err := apps.ParseAppEnv(raw)
		if err != nil {
			t.Errorf("ParseAppEnv(%q) = %v, want nil", raw, err)
		}
		if got != want {
			t.Errorf("ParseAppEnv(%q) = %v, want %v", raw, got, want)
		}
		if got.String() != raw {
			t.Errorf("String() = %q, want %q", got.String(), raw)
		}
	}
}

func TestParseAppEnvRejectsUnknown(t *testing.T) {
	got, err := apps.ParseAppEnv("prod")
	if err == nil {
		t.Fatal("ParseAppEnv() = nil error, want error")
	}
	if got != apps.Unknown {
		t.Errorf("ParseAppEnv() = %v, want Unknown", got)
	}
}

func TestZeroValueIsUnknown(t *testing.T) {
	var e apps.AppEnv

	if e != apps.Unknown {
		t.Errorf("zero AppEnv = %v, want Unknown", e)
	}
	if e.IsProduction() {
		t.Error("zero AppEnv reports IsProduction() = true")
	}
}
