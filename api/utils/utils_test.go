package utils

import (
	"testing"
)

func TestTrimIP(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"192.168.1.1\r\n", "192.168.1.1"},
		{"10.0.0.1\n", "10.0.0.1"},
		{"127.0.0.1", "127.0.0.1"},
	}

	for _, test := range tests {
		result := trimIP(test.input)
		if result != test.expected {
			t.Errorf("Expected: %s, got: %s", test.expected, result)
		}
	}
}

func TestIP(t *testing.T) {
	ip := "192.168.1.1\r\n"
	result := trimIP(ip)
	t.Logf("\rRESULT: %v\n", result)
}

func TestHash(t *testing.T) {
	msg := "this is a message"
	hash := hashMessage([]byte(msg), len(msg))
	t.Logf("\rHash: %s\n", hash)
}
