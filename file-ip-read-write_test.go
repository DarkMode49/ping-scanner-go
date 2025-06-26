package main

import (
	"fmt"
	"testing"
)

func TestFileIPReader(t *testing.T) {
	fileIPReader := &FileIPReader{}
	file, err := fileIPReader.Open("ips.txt")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	tests := []struct {
		line     int64
		expected string
	}{
		{0, "192.168.1.1"},
		{1, "192.168.1.2"},
		{2, "192.168.1.3"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("line=%d", tt.line), func(t *testing.T) {
			ip, err := fileIPReader.ReadIP(tt.line)

			if err != nil {
				t.Fatalf("failed to read IP: %v", err)
			}
			if ip != tt.expected {
				t.Errorf("got %q, want %q", ip, tt.expected)
			}
		})
	}
}
