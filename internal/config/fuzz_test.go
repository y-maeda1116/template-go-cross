package config

import (
	"os"
	"path/filepath"
	"testing"
)

// FuzzLoad は設定ファイル読み込みが任意入力でパニックや不正状態を生まないことを検証する。
// Scorecard の Fuzzing チェックは *_test.go 内の FuzzXxx(*testing.F) を検出する。
func FuzzLoad(f *testing.F) {
	f.Add([]byte("app:\n  name: test\n  version: 1.0.0\n"))
	f.Add([]byte("server:\n  host: localhost\n  port: 8080\n"))
	f.Add([]byte("logging:\n  level: info\n  format: json\n"))
	f.Add([]byte("desktop:\n  window:\n    width: 800\n    height: 600\n    resizable: true\n"))
	f.Add([]byte("not: valid: yaml:"))
	f.Add([]byte(""))

	f.Fuzz(func(t *testing.T, data []byte) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatalf("failed to write temp config: %v", err)
		}

		cfg, err := Load(path)
		if err != nil {
			// 不正な入力はエラーとして返ればよい
			return
		}
		if cfg == nil {
			t.Fatal("Load returned nil config without error")
		}
	})
}
