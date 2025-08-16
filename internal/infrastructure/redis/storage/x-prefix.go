package storage

import (
	"strings"
)

type Prefix string

const (
	Message Prefix = "message"
	Voice   Prefix = "voice"
)

func NewPrefix(prefixes ...Prefix) Prefix {
	if len(prefixes) == 0 {
		return ""
	}

	ss := make([]string, len(prefixes))
	for i, p := range prefixes {
		ss[i] = string(p)
	}
	return Prefix(strings.Join(ss, ":"))
}

// WithValue example: Message.WithValue("123") → "message:123"
func (p Prefix) WithValue(value string) string {
	if p == "" {
		return value
	}
	return string(p) + ":" + value
}

func (p Prefix) String() string {
	return string(p)
}
