package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Loader struct {
	problems []string
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Err() error {
	if len(l.problems) == 0 {
		return nil
	}

	return fmt.Errorf("environment: %s", strings.Join(l.problems, "; "))
}

func (l *Loader) missing(key string) {
	l.problems = append(l.problems, fmt.Sprintf("%s is required", key))
}

func (l *Loader) invalid(key, value string, err error) {
	l.problems = append(l.problems, fmt.Sprintf("%s=%q is invalid: %v", key, value, err))
}

func (l *Loader) lookup(key string) (string, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", false
	}

	value = strings.TrimSpace(value)
	if value == "" {
		return "", false
	}

	return value, true
}

func (l *Loader) String(key string) string {
	value, ok := l.lookup(key)
	if !ok {
		l.missing(key)

		return ""
	}

	return value
}

func (l *Loader) StringDefault(key, fallback string) string {
	value, ok := l.lookup(key)
	if !ok {
		return fallback
	}

	return value
}

func (l *Loader) Int(key string) int {
	raw, ok := l.lookup(key)
	if !ok {
		l.missing(key)

		return 0
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		l.invalid(key, raw, err)

		return 0
	}

	return value
}

func (l *Loader) BoolDefault(key string, fallback bool) bool {
	raw, ok := l.lookup(key)
	if !ok {
		return fallback
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		l.invalid(key, raw, err)

		return fallback
	}

	return value
}

func (l *Loader) DurationDefault(key string, fallback time.Duration) time.Duration {
	raw, ok := l.lookup(key)
	if !ok {
		return fallback
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		l.invalid(key, raw, err)

		return fallback
	}

	return value
}

func (l *Loader) StringSlice(key, separator string) []string {
	raw, ok := l.lookup(key)
	if !ok {
		l.missing(key)

		return nil
	}

	parts := strings.Split(raw, separator)
	values := make([]string, 0, len(parts))

	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			values = append(values, part)
		}
	}

	if len(values) == 0 {
		l.invalid(key, raw, fmt.Errorf("no values"))

		return nil
	}

	return values
}

func (l *Loader) StringSliceDefault(key, separator string, fallback []string) []string {
	if _, ok := l.lookup(key); !ok {
		return fallback
	}

	return l.StringSlice(key, separator)
}

func (l *Loader) Custom(key string, parse func(string) error) {
	raw, ok := l.lookup(key)
	if !ok {
		l.missing(key)

		return
	}

	if err := parse(raw); err != nil {
		l.invalid(key, raw, err)
	}
}
