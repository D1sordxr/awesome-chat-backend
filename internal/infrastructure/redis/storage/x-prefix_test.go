package storage

import "testing"

func TestNewPrefix(t *testing.T) {
	tests := []struct {
		name     string
		prefixes []Prefix
		expected string
	}{
		{
			name:     "Single prefix",
			prefixes: []Prefix{Message},
			expected: "message",
		},
		{
			name:     "Multiple prefixes",
			prefixes: []Prefix{Voice, Message},
			expected: "voice:message",
		},
		{
			name:     "No prefixes",
			prefixes: []Prefix{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := NewPrefix(tt.prefixes...)
			if res.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, res.String())
			}
		})
	}
}

func TestWithValue(t *testing.T) {
	tests := []struct {
		name     string
		prefix   Prefix
		value    string
		expected string
	}{
		{
			name:     "With non-empty prefix",
			prefix:   Message,
			value:    "123",
			expected: "message:123",
		},
		{
			name:     "With empty prefix",
			prefix:   "",
			value:    "123",
			expected: "123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.prefix.WithValue(tt.value)
			if res != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, res)
			}
		})
	}
}
