package main

import "testing"

func TestSanitizeMessage(t *testing.T) {
	tests := []struct{
		name string 
		input []byte
		expected string
	} {
		{"no special characters", []byte("hello world"), "hello world"},
		{"yes special characters", []byte("hello\rworld\n"), "hello_world_"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			got := sanitizeMessage(tt.input)
			if got != tt.expected {
				t.Errorf("santizeMessage(message) = %q, want %q", got, tt.expected)
			}
		})
	}
}

