package env_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"awesome-chat/internal/infrastructure/config/env"
)

func TestRequiredValueIsRead(t *testing.T) {
	t.Setenv("SOME_HOST", "db")

	l := env.NewLoader()
	got := l.String("SOME_HOST")

	if err := l.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}
	if got != "db" {
		t.Errorf("String() = %q, want db", got)
	}
}

func TestMissingRequiredValuesAreReportedTogether(t *testing.T) {
	l := env.NewLoader()
	l.String("MISSING_ONE")
	l.Int("MISSING_TWO")

	err := l.Err()
	if err == nil {
		t.Fatal("Err() = nil, want error")
	}
	for _, key := range []string{"MISSING_ONE", "MISSING_TWO"} {
		if !strings.Contains(err.Error(), key) {
			t.Errorf("error %q does not mention %s", err, key)
		}
	}
}

func TestBlankValueCountsAsMissing(t *testing.T) {
	t.Setenv("BLANK_VALUE", "   ")

	l := env.NewLoader()
	l.String("BLANK_VALUE")

	if l.Err() == nil {
		t.Fatal("Err() = nil, want error for blank value")
	}
}

func TestDefaultsApplyWhenUnset(t *testing.T) {
	l := env.NewLoader()

	if got := l.StringDefault("UNSET_STRING", "fallback"); got != "fallback" {
		t.Errorf("StringDefault() = %q, want fallback", got)
	}
	if got := l.BoolDefault("UNSET_BOOL", true); !got {
		t.Errorf("BoolDefault() = false, want true")
	}
	if got := l.DurationDefault("UNSET_DURATION", 4*time.Second); got != 4*time.Second {
		t.Errorf("DurationDefault() = %v, want 4s", got)
	}
	if err := l.Err(); err != nil {
		t.Errorf("Err() = %v, want nil", err)
	}
}

func TestInvalidIntIsReported(t *testing.T) {
	t.Setenv("BAD_PORT", "not-a-number")

	l := env.NewLoader()
	l.Int("BAD_PORT")

	err := l.Err()
	if err == nil {
		t.Fatal("Err() = nil, want error")
	}
	if !strings.Contains(err.Error(), "BAD_PORT") {
		t.Errorf("error %q does not mention the key", err)
	}
}

func TestStringSliceSplitsAndTrims(t *testing.T) {
	t.Setenv("BROKERS", "a:9092, b:9092 ,")

	l := env.NewLoader()
	got := l.StringSlice("BROKERS", ",")

	if err := l.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}
	want := []string{"a:9092", "b:9092"}
	if len(got) != len(want) {
		t.Fatalf("StringSlice() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("StringSlice()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestCustomReportsParseError(t *testing.T) {
	t.Setenv("APP_ENV", "nonsense")

	l := env.NewLoader()
	l.Custom("APP_ENV", func(string) error { return errors.New("bad value") })

	if l.Err() == nil {
		t.Fatal("Err() = nil, want error")
	}
}
