package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		wantNil bool
	}{
		{
			name:    "returns config",
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Load(); got == nil && !tt.wantNil {
				t.Error("Load() returned nil")
			}
		})
	}
}

func TestLoad_WithEnvVar(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	cfg := Load()
	if cfg == nil {
		t.Fatal("Load() returned nil")
	}
}
