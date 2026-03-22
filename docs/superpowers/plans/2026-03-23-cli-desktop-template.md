# CLI + Desktop Go Template Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 単一リポジトリ内でCLIとデスクトップアプリケーション（Wails）を共存させるGoテンプレートを作成する。

**Architecture:** コアビジネスロジックを共有し、CLI（Cobra）とデスクトップ（Wails）は別々のエントリーポイントから同じコアを呼び出す3層アーキテクチャ（Application → Core → Infrastructure）。

**Tech Stack:** Go 1.22, Cobra, Viper, zap, Wails v2, testify, mockgen, Air

---

## Task 1: 依存関係の追加

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: 依存関係をgo.modに追加**

```bash
go get github.com/spf13/cobra@v1.8.1
go get github.com/spf13/viper@v1.19.0
go get go.uber.org/zap@v1.27.0
go get github.com/wailsapp/wails/v2@v2.9.1
go get github.com/stretchr/testify@v1.9.0
go get golang.org/x/mock@v1.6.0
```

- [ ] **Step 2: go mod tidyを実行**

Run: `go mod tidy`

Expected: go.modとgo.sumが更新される

- [ ] **Step 3: コミット**

```bash
git add go.mod go.sum
git commit -m "deps: add dependencies for CLI + Desktop template"
```

---

## Task 2: カスタムエラー型の作成

**Files:**
- Create: `internal/core/errors.go`
- Create: `internal/core/errors_test.go`

- [ ] **Step 1: テストを先に書く**

```go
// internal/core/errors_test.go
package core

import (
    "errors"
    "testing"
)

func TestAppError_Error(t *testing.T) {
    err := errors.New("original error")
    appErr := NewAppError("CODE001", "test message", err)

    if appErr.Error() != "test message" {
        t.Errorf("Expected 'test message', got '%s'", appErr.Error())
    }
}

func TestAppError_Unwrap(t *testing.T) {
    originalErr := errors.New("original error")
    appErr := NewAppError("CODE001", "test message", originalErr)

    if errors.Unwrap(appErr) != originalErr {
        t.Error("Unwrap should return the original error")
    }
}

func TestAppError_Code(t *testing.T) {
    appErr := NewAppError("CODE001", "test message", nil)

    if appErr.Code != "CODE001" {
        t.Errorf("Expected 'CODE001', got '%s'", appErr.Code)
    }
}
```

- [ ] **Step 2: テストを実行して失敗を確認**

Run: `go test ./internal/core/ -v`

Expected: FAIL with "undefined: NewAppError"

- [ ] **Step 3: エラー型を実装**

```go
// internal/core/errors.go
package core

import (
    "errors"
)

var (
    // 設定エラー
    ErrConfigNotFound = errors.New("configuration not found")
    ErrConfigInvalid  = errors.New("invalid configuration")

    // ビジネスロジックエラー
    ErrNotFound      = errors.New("resource not found")
    ErrInvalidInput  = errors.New("invalid input")
    ErrPermission    = errors.New("permission denied")
)

// AppError アプリケーションエラー
type AppError struct {
    Code    string
    Message string
    Err     error
}

func (e *AppError) Error() string {
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// NewAppError 新しいAppErrorを作成
func NewAppError(code, message string, err error) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Err:     err,
    }
}
```

- [ ] **Step 4: テストを実行して成功を確認**

Run: `go test ./internal/core/ -v`

Expected: PASS

- [ ] **Step 5: コミット**

```bash
git add internal/core/errors.go internal/core/errors_test.go
git commit -m "feat: add custom AppError type for error handling"
```

---

## Task 3: ロガーの作成

**Files:**
- Create: `internal/logger/logger.go`
- Create: `internal/logger/logger_test.go`

- [ ] **Step 1: テストを先に書く**

```go
// internal/logger/logger_test.go
package logger

import (
    "testing"
)

func TestNew(t *testing.T) {
    // Infoレベルでロガーを作成
    log := New("info")

    if log == nil {
        t.Fatal("Expected non-nil logger")
    }
}

func TestLogger_Debug(t *testing.T) {
    log := New("debug")

    // パニックしないことを確認
    log.Debug("debug message", "key", "value")
}

func TestLogger_Info(t *testing.T) {
    log := New("info")

    // パニックしないことを確認
    log.Info("info message", "key", "value")
}

func TestLogger_Error(t *testing.T) {
    log := New("error")

    // パニックしないことを確認
    log.Error("error message", "key", "value")
}

func TestLogger_Sync(t *testing.T) {
    log := New("info")

    // パニックしないことを確認
    log.Sync()
}
```

- [ ] **Step 2: テストを実行して失敗を確認**

Run: `go test ./internal/logger/ -v`

Expected: FAIL with "undefined: New"

- [ ] **Step 3: ロガーを実装**

```go
// internal/logger/logger.go
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// Logger 構造化ロガー
type Logger struct {
    *zap.SugaredLogger
}

// New 新しいロガーを作成
func New(level string) *Logger {
    // レベルのパース
    var zapLevel zapcore.Level
    if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
        zapLevel = zapcore.InfoLevel
    }

    // コンフィグ作成
    config := zap.Config{
        Level:            zap.NewAtomicLevelAt(zapLevel),
        Development:       false,
        Encoding:         "json",
        EncoderConfig:     zap.NewProductionEncoderConfig(),
        OutputPaths:      []string{"stdout"},
        ErrorOutputPaths:  []string{"stderr"},
    }

    // ロガー作成
    logger, _ := config.Build()
    return &Logger{
        SugaredLogger: logger.Sugar(),
    }
}

// Debug デバッグログ
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
    l.SugaredLogger.Debugw(msg, keyvals...)
}

// Info 情報ログ
func (l *Logger) Info(msg string, keyvals ...interface{}) {
    l.SugaredLogger.Infow(msg, keyvals...)
}

// Warn 警告ログ
func (l *Logger) Warn(msg string, keyvals ...interface{}) {
    l.SugaredLogger.Warnw(msg, keyvals...)
}

// Error エラーログ
func (l *Logger) Error(msg string, keyvals ...interface{}) {
    l.SugaredLogger.Errorw(msg, keyvals...)
}

// Sync ログをフラッシュ
func (l *Logger) Sync() {
    l.SugaredLogger.Sync()
}
```

- [ ] **Step 4: テストを実行して成功を確認**

Run: `go test ./internal/logger/ -v`

Expected: PASS

- [ ] **Step 5: コミット**

```bash
git add internal/logger/logger.go internal/logger/logger_test.go
git commit -m "feat: add zap-based logger with structured logging"
```

---

## Task 4: 設定管理の実装

**Files:**
- Modify: `internal/config/config.go`
- Modify: `internal/config/config_test.go`

- [ ] **Step 1: 既存のconfig.goを読む**

Run: `cat internal/config/config.go`

- [ ] **Step 2: テストを書く**

```go
// internal/config/config_test.go
package config

import (
    "os"
    "path/filepath"
    "testing"
)

func TestLoad(t *testing.T) {
    // 一時的な設定ファイルを作成
    tmpDir := t.TempDir()
    configPath := filepath.Join(tmpDir, "config.yaml")
    content := `
app:
  name: test-app
  version: 1.0.0
logging:
  level: debug
`
    os.WriteFile(configPath, []byte(content), 0644)

    cfg, err := Load(configPath)

    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    if cfg.App.Name != "test-app" {
        t.Errorf("Expected 'test-app', got '%s'", cfg.App.Name)
    }
}

func TestLoad_FileNotFound(t *testing.T) {
    _, err := Load("/nonexistent/config.yaml")

    if err == nil {
        t.Error("Expected error for non-existent file")
    }
}
```

- [ ] **Step 3: テストを実行して失敗を確認**

Run: `go test ./internal/config/ -v`

Expected: FAIL with "undefined: Config" or "undefined: Load"

- [ ] **Step 4: 設定管理を実装**

```go
// internal/config/config.go
package config

import (
    "fmt"

    "github.com/spf13/viper"
)

// Config 設定構造体
type Config struct {
    App      AppConfig
    Server   ServerConfig
    Logging  LoggingConfig
    CLI      CLIConfig
    Desktop  DesktopConfig
}

// AppConfig アプリケーション設定
type AppConfig struct {
    Name        string
    Version     string
    Environment string
}

// ServerConfig サーバー設定
type ServerConfig struct {
    Host string
    Port int
}

// LoggingConfig ログ設定
type LoggingConfig struct {
    Level  string
    Format string
    Output string
}

// CLIConfig CLI設定
type CLIConfig struct {
    Theme string
}

// DesktopConfig デスクトップ設定
type DesktopConfig struct {
    Window WindowConfig
}

// WindowConfig ウィンドウ設定
type WindowConfig struct {
    Width     int
    Height    int
    Resizable bool
}

// Load 設定ファイルを読み込む
func Load(configPath string) (*Config, error) {
    v := viper.New()

    // 設定ファイルを設定
    v.SetConfigFile(configPath)

    // 設定ファイルを読み込み
    if err := v.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }

    // 構造体にマッピング
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }

    return &cfg, nil
}
```

- [ ] **Step 5: テストを実行して成功を確認**

Run: `go test ./internal/config/ -v`

Expected: PASS

- [ ] **Step 6: コミット**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "refactor: enhance config with Viper and comprehensive settings"
```

---

## Task 5: コアビジネスロジックの作成

**Files:**
- Create: `internal/core/service.go`
- Create: `internal/core/service_test.go`

- [ ] **Step 1: テストを先に書く**

```go
// internal/core/service_test.go
package core

import (
    "testing"
)

func TestService_SayHello(t *testing.T) {
    s := NewService()

    result, err := s.SayHello("World")

    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    if result != "Hello, World!" {
        t.Errorf("Expected 'Hello, World!', got '%s'", result)
    }
}

func TestService_SayHello_EmptyName(t *testing.T) {
    s := NewService()

    _, err := s.SayHello("")

    if err == nil {
        t.Error("Expected error for empty name")
    }

    if err != ErrInvalidInput {
        t.Errorf("Expected ErrInvalidInput, got %v", err)
    }
}

func TestService_GetVersion(t *testing.T) {
    s := NewService()

    version := s.GetVersion()

    if version == "" {
        t.Error("Expected non-empty version")
    }
}
```

- [ ] **Step 2: テストを実行して失敗を確認**

Run: `go test ./internal/core/ -v`

Expected: FAIL with "undefined: NewService"

- [ ] **Step 3: サービスを実装**

```go
// internal/core/service.go
package core

// Service コアビジネスロジック
type Service struct {
    version string
}

// NewService 新しいServiceを作成
func NewService() *Service {
    return &Service{
        version: "1.0.0",
    }
}

// SayHello あいさつメッセージを生成
func (s *Service) SayHello(name string) (string, error) {
    if name == "" {
        return "", ErrInvalidInput
    }
    return "Hello, " + name + "!", nil
}

// GetVersion バージョンを取得
func (s *Service) GetVersion() string {
    return s.version
}

// SetVersion バージョンを設定（テスト用）
func (s *Service) SetVersion(version string) {
    s.version = version
}
```

- [ ] **Step 4: テストを実行して成功を確認**

Run: `go test ./internal/core/ -v`

Expected: PASS

- [ ] **Step 5: コミット**

```bash
git add internal/core/service.go internal/core/service_test.go
git commit -m "feat: add core service with SayHello and GetVersion"
```

---

## Task 6: CLIコマンドの作成

**Files:**
- Create: `internal/cli/root.go`
- Create: `internal/cli/version.go`
- Create: `internal/cli/hello.go`
- Create: `cmd/cli/main.go`

- [ ] **Step 1: ルートコマンドを作成**

```go
// internal/cli/root.go
package cli

import (
    "github.com/spf13/cobra"
)

// NewRootCommand ルートコマンドを作成
func NewRootCommand(version string) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "myapp",
        Short: "A CLI and Desktop application",
        Long:  "MyApp is a CLI and Desktop application template built with Go and Wails.",
        Version: version,
    }

    // サブコマンドを追加
    cmd.AddCommand(NewVersionCommand())
    cmd.AddCommand(NewHelloCommand())

    return cmd
}
```

- [ ] **Step 2: バージョンコマンドを作成**

```go
// internal/cli/version.go
package cli

import (
    "fmt"

    "github.com/spf13/cobra"
)

// NewVersionCommand バージョンコマンドを作成
func NewVersionCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "version",
        Short: "Print the version number",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("MyApp v1.0.0")
        },
    }
}
```

- [ ] **Step 3: helloコマンドを作成**

```go
// internal/cli/hello.go
package cli

import (
    "fmt"

    "github.com/spf13/cobra"
)

// NewHelloCommand helloコマンドを作成
func NewHelloCommand() *cobra.Command {
    var name string

    cmd := &cobra.Command{
        Use:   "hello",
        Short: "Say hello",
        Run: func(cmd *cobra.Command, args []string) {
            if name == "" {
                name = "World"
            }
            fmt.Printf("Hello, %s!\n", name)
        },
    }

    cmd.Flags().StringVarP(&name, "name", "n", "", "Name to greet")

    return cmd
}
```

- [ ] **Step 4: CLIエントリーポイントを作成**

```go
// cmd/cli/main.go
package main

import (
    "os"

    "github.com/user/repo/internal/cli"
)

func main() {
    cmd := cli.NewRootCommand("1.0.0")
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

- [ ] **Step 5: 動作確認**

Run: `go run ./cmd/cli --help`

Expected: ヘルプメッセージが表示される

Run: `go run ./cmd/cli hello --name Claude`

Expected: "Hello, Claude!"と表示される

- [ ] **Step 6: コミット**

```bash
git add internal/cli/ cmd/cli/
git commit -m "feat: add CLI commands with Cobra (root, version, hello)"
```

---

## Task 7: Wailsプロジェクトの初期化

**Files:**
- Create: `wails.json`
- Create: `cmd/desktop/main.go`
- Create: `internal/ui/app.go`

- [ ] **Step 1: Wailsを初期化**

```bash
wails init -n myapp -t react
```

Note: プロンプトで質問がされたらすべてEnterキーでデフォルトを受け入れる

- [ ] **Step 2: wails.jsonを修正**

```json
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "myapp",
  "outputfilename": "MyApp",
  "frontend": {
    "install": "npm install",
    "build": "npm run build",
    "dev": "npm run dev",
    "package": "./frontend/dist"
  },
  "author": {
    "name": "Your Name",
    "email": "your@email.com"
  },
  "info": {
    "companyName": "Your Company",
    "productName": "My App",
    "productVersion": "1.0.0",
    "copyright": "Copyright © 2026",
    "comments": "A CLI + Desktop application"
  },
  "wailsjsdir": "./frontend",
  "version": "2"
}
```

- [ ] **Step 3: デスクトップエントリーポイントを作成**

```go
// cmd/desktop/main.go
package main

import (
    "embed"

    "github.com/wailsapp/wails/v2"
    "github.com/wailsapp/wails/v2/pkg/options"
    "github.com/wailsapp/wails/v2/pkg/options/assetserver"
    "github.com/user/repo/internal/ui"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    // アプリインスタンスを作成
    app := ui.NewApp()

    // デスクトップアプリを開始
    err := wails.Run(&options.App{
        Title:  "My App",
        Width:  1024,
        Height: 768,
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        OnStartup:        app.Startup,
        OnShutdown:      app.Shutdown,
        Bind: []interface{}{
            app,
        },
    })

    if err != nil {
        println("Error:", err.Error())
    }
}
```

- [ ] **Step 4: Wailsアプリバインディングを作成**

```go
// internal/ui/app.go
package ui

// App Wailsアプリ
type App struct{}

// NewApp 新しいAppを作成
func NewApp() *App {
    return &App{}
}

// Startup アプリ起動時に呼ばれる
func (a *App) Startup(ctx context.Context) {
    // 初期化処理
}

// Shutdown アプリ終了時に呼ばれる
func (a *App) Shutdown(ctx context.Context) {
    // クリーンアップ処理
}

// Greet あいさつを返す
func (a *App) Greet(name string) string {
    if name == "" {
        name = "World"
    }
    return "Hello, " + name + "!"
}

// Version バージョンを返す
func (a *App) Version() string {
    return "1.0.0"
}
```

- [ ] **Step 5: frontend/package.jsonを修正**

```json
{
  "name": "myapp-frontend",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "@wailsapp/runtime": "^2.0.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "typescript": "^5.0.0",
    "vite": "^5.0.0"
  }
}
```

- [ ] **Step 6: frontend/src/App.tsxを作成**

```tsx
// frontend/src/App.tsx
import { useState } from 'react'
import { Greet, Version } from '../wailsjs/go/main/App'
import './App.css'

function App() {
  const [name, setName] = useState('')
  const [result, setResult] = useState('')
  const [version, setVersion] = useState('Loading...')

  useState(() => {
    Version().then(v => setVersion(v))
  })

  function doGreet() {
    Greet(name).then(setResult)
  }

  return (
    <div className="App">
      <div className="container">
        <h1>My App v{version}</h1>
        <div className="input-group">
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Enter your name"
          />
          <button onClick={doGreet}>Greet</button>
        </div>
        {result && <div className="result">{result}</div>}
      </div>
    </div>
  )
}

export default App
```

- [ ] **Step 7: frontend/src/App.cssを作成**

```css
/* frontend/src/App.css */
.App {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: #f5f5f5;
}

.container {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

h1 {
  margin-top: 0;
  color: #333;
}

.input-group {
  display: flex;
  gap: 8px;
  margin-top: 1rem;
}

input {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  flex: 1;
}

button {
  padding: 8px 16px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background: #0056b3;
}

.result {
  margin-top: 1rem;
  padding: 1rem;
  background: #e7f3ff;
  border-radius: 4px;
}
```

- [ ] **Step 8: コミット**

```bash
git add wails.json cmd/desktop/ internal/ui/ frontend/
git commit -m "feat: add Wails desktop application with React frontend"
```

---

## Task 8: Makefileの拡張

**Files:**
- Modify: `Makefile`

- [ ] **Step 1: Makefileを拡張**

```makefile
.PHONY: build build-cli build-desktop run-cli run-desktop clean test test-coverage test-race mocks lint fmt help

APP_NAME := myapp
BIN_DIR := bin

# CLI
CLI_MAIN := ./cmd/cli

# Desktop
DESKTOP_MAIN := ./cmd/desktop

# --- ビルド ---

build-cli:
	@echo "Building CLI for current OS..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(CLI_MAIN)

build-desktop:
	@echo "Building Desktop for current OS..."
	@cd frontend && npm install && npm run build
	@wails build

build-all: build-cli build-desktop

# --- 実行 ---

run-cli:
	@echo "Running CLI..."
	@go run $(CLI_MAIN) $(ARGS)

run-desktop:
	@echo "Running Desktop..."
	@wails dev

# --- テスト ---

test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

test-race:
	@echo "Running tests with race detector..."
	@go test -race -v ./...

# --- モック生成 ---

mocks:
	@echo "Generating mocks..."
	@mkdir -p test/mocks
	@mockgen -source=internal/core/service.go -destination=test/mocks/service_mock.go

# --- Lint / Format ---

fmt:
	@go fmt ./...

lint:
	@echo "Installing golangci-lint if needed..."
	@which golangci-lint || (curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest)
	@echo "Running linter..."
	@golangci-lint run ./...

# --- クリーンアップ ---

clean:
	@rm -rf $(BIN_DIR) coverage.out coverage.html test/mocks
	@cd frontend && rm -rf node_modules wailsjs dist

# --- ヘルプ ---

help:
	@echo "Available targets:"
	@echo "  build-cli       - Build CLI for current OS"
	@echo "  build-desktop    - Build Desktop for current OS"
	@echo "  build-all        - Build CLI and Desktop"
	@echo "  run-cli          - Run CLI (use ARGS=\"--help\" for options)"
	@echo "  run-desktop      - Run Desktop in dev mode"
	@echo "  test             - Run all tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  test-race        - Run tests with race detector"
	@echo "  mocks            - Generate mocks"
	@echo "  fmt              - Format Go code"
	@echo "  lint             - Run linter"
	@echo "  clean            - Remove build artifacts"
```

- [ ] **Step 2: helpターゲットを確認**

Run: `make help`

Expected: 全ターゲットの説明が表示される

- [ ] **Step 3: コミット**

```bash
git add Makefile
git commit -m "feat: extend Makefile with CLI/Desktop targets and tooling"
```

---

## Task 9: Air設定の作成

**Files:**
- Create: `air.toml`

- [ ] **Step 1: air.tomlを作成**

```toml
# air.toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ./cmd/cli"
  bin = "tmp/main"
  include_ext = ["go"]
  exclude_dir = ["tmp", "vendor", "test", "frontend", "bin"]
  include_dir = []
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""

[log]
  time = true

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"
```

- [ ] **Step 2: コミット**

```bash
git add air.toml
git commit -m "feat: add Air hot-reload configuration"
```

---

## Task 10: 設定ファイルテンプレートの作成

**Files:**
- Modify: `.env.example`
- Create: `config.yaml.example`

- [ ] **Step 1: .env.exampleを拡張**

```bash
# アプリケーション
APP_NAME=myapp
APP_ENV=development
APP_PORT=8080

# ログレベル (debug, info, warn, error)
LOG_LEVEL=info

# 設定ファイルパス
CONFIG_PATH=config.yaml
```

- [ ] **Step 2: config.yaml.exampleを作成**

```yaml
# アプリケーション設定
app:
  name: myapp
  version: 1.0.0
  environment: development

# サーバー設定
server:
  host: localhost
  port: 8080

# ログ設定
logging:
  level: info
  format: json
  output: stdout

# CLI設定
cli:
  theme: auto

# デスクトップ設定
desktop:
  window:
    width: 1200
    height: 800
    resizable: true
```

- [ ] **Step 3: コミット**

```bash
git add .env.example config.yaml.example
git commit -m "feat: add configuration file templates"
```

---

## Task 11: GitHub Actionsの更新

**Files:**
- Modify: `.github/workflows/test.yml`
- Create: `.github/workflows/build-cli.yml`
- Create: `.github/workflows/build-desktop.yml`

- [ ] **Step 1: test.ymlを更新**

```yaml
# .github/workflows/test.yml
name: Test

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.22'

    - name: Cache Go modules
      uses: actions/cache@v4
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-

    - name: Install dependencies
      run: go mod download

    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Upload coverage
      uses: codecov/codecov-action@v4
      with:
        files: ./coverage.out
        flags: unittests
        name: codecov-umbrella
```

- [ ] **Step 2: build-cli.ymlを作成**

```yaml
# .github/workflows/build-cli.yml
name: Build CLI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  build:
    runs-on: ${{ matrix.os }}

    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ['1.22']

    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: ${{ matrix.go-version }}

    - name: Cache Go modules
      uses: actions/cache@v4
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}

    - name: Install dependencies
      run: go mod download

    - name: Build CLI
      run: go build -o bin/myapp ./cmd/cli

    - name: Test CLI
      run: go test ./...

    - name: Upload artifact
      uses: actions/upload-artifact@v4
      with:
        name: myapp-${{ matrix.os }}-${{ matrix.go-version }}
        path: bin/myapp
```

- [ ] **Step 3: build-desktop.ymlを作成**

```yaml
# .github/workflows/build-desktop.yml
name: Build Desktop

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  build:
    runs-on: ${{ matrix.os }}

    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]

    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.22'

    - name: Set up Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '20'

    - name: Cache Go modules
      uses: actions/cache@v4
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}

    - name: Cache Node modules
      uses: actions/cache@v4
      with:
        path: frontend/node_modules
        key: ${{ runner.os }}-node-${{ hashFiles('**/frontend/package-lock.json') }}

    - name: Install Wails CLI
      run: go install github.com/wailsapp/wails/v2/cmd/wails@latest

    - name: Install dependencies
      run: |
        go mod download
        cd frontend && npm install

    - name: Build Desktop
      run: wails build

    - name: Upload artifact
      uses: actions/upload-artifact@v4
      with:
        name: MyApp-${{ matrix.os }}
        path: build/
```

- [ ] **Step 4: コミット**

```bash
git add .github/workflows/
git commit -m "ci: add GitHub Actions for test, CLI build, and Desktop build"
```

---

## Task 12: モック生成と最終確認

**Files:**
- Test: `test/mocks/`

- [ ] **Step 1: モックを生成**

Run: `make mocks`

Expected: test/mocks/にモックファイルが作成される

- [ ] **Step 2: 全テストを実行**

Run: `make test`

Expected: 全テストがパスする

- [ ] **Step 3: カバレッジを確認**

Run: `make test-coverage`

Expected: カバレッジレポートが生成され、80%以上であること

- [ ] **Step 4: CLI動作確認**

Run: `make run-cli ARGS="--help"`

Expected: ヘルプメッセージが表示される

Run: `make run-cli ARGS="hello --name Test"`

Expected: "Hello, Test!"と表示される

- [ ] **Step 5: デスクトップ動作確認**

Run: `make run-desktop`

Expected: Wails devサーバーが起動し、デスクトップアプリが開く（GUI環境の場合）

- [ ] **Step 6: モックをコミット**

```bash
git add test/mocks/
git commit -m "test: add generated mocks for core service"
```

---

## Task 13: READMEの更新

**Files:**
- Modify: `README.md`

- [ ] **Step 1: READMEを更新**

```markdown
# Go CLI + Desktop Template

A cross-platform Go application template supporting both CLI and Desktop (Wails) interfaces.

## Features

- **Dual Interface:** CLI (Cobra) and Desktop (Wails + React)
- **Shared Core:** Business logic shared between CLI and Desktop
- **Structured Logging:** zap-based logging with multiple levels
- **Configuration:** Viper-based config with YAML/TOML/JSON support
- **Hot Reload:** Air for CLI, Wails dev for Desktop
- **Testing:** testify + mockgen with 80%+ coverage goal
- **CI/CD:** GitHub Actions for testing and building

## Project Structure

```
.
├── cmd/
│   ├── cli/           # CLI entry point (Cobra)
│   └── desktop/       # Wails desktop entry point
├── internal/
│   ├── core/          # Shared business logic
│   ├── ui/            # Wails UI bindings
│   ├── config/        # Configuration (Viper)
│   ├── logger/        # Logging (zap)
│   └── cli/          # CLI commands
├── frontend/         # Wails frontend (React)
├── pkg/              # Reusable packages
├── test/
│   └── mocks/        # Generated mocks
├── Makefile          # Build targets
├── wails.json        # Wails configuration
└── air.toml          # Hot reload configuration
```

## Getting Started

### Prerequisites

- Go 1.22 or later
- Node.js 20 or later (for Desktop)
- Make

### Installation

```bash
# Clone the repository
git clone https://github.com/your-username/your-repo.git
cd your-repo

# Install Go dependencies
go mod download

# Install frontend dependencies (for Desktop)
cd frontend && npm install && cd ..
```

### Configuration

Create environment and config files:

```bash
cp .env.example .env
cp config.yaml.example config.yaml
# Edit .env and config.yaml with your settings
```

## Usage

### CLI

```bash
# Run CLI
make run-cli

# Show help
make run-cli ARGS="--help"

# Say hello
make run-cli ARGS="hello --name World"

# Build CLI
make build-cli
```

### Desktop

```bash
# Run Desktop in dev mode
make run-desktop

# Build Desktop
make build-desktop
```

### Development

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detector
make test-race

# Generate mocks
make mocks

# Format code
make fmt

# Run linter
make lint

# Clean build artifacts
make clean
```

## Architecture

```
Application Layer
├── CLI (Cobra)
└── Desktop (Wails + React)
         ↓
Core Business Logic Layer
└── Shared services
         ↓
Infrastructure Layer
├── Config (Viper)
└── Logger (zap)
```

## License

MIT
```

- [ ] **Step 2: コミット**

```bash
git add README.md
git commit -m "docs: update README with CLI + Desktop template information"
```

---

## Task 14: .gitignoreの更新

**Files:**
- Modify: `.gitignore`

- [ ] **Step 1: .gitignoreを更新**

```gitignore
# Binaries
bin/
tmp/
*.exe
*.dll
*.so
*.dylib
build/

# Test coverage
coverage.out
coverage.html

# Air
tmp/

# Frontend
frontend/node_modules/
frontend/dist/
frontend/wailsjs/

# Environment
.env
.env.local
*.log

# Editor
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db
```

- [ ] **Step 2: コミット**

```bash
git add .gitignore
git commit -m "chore: update .gitignore for CLI + Desktop template"
```

---

## Task 15: 最終確認とタグ付け

**Files:**

- [ ] **Step 1: 全変更を確認**

Run: `git status`

Expected: mainブランチ上に変更がないこと

- [ ] **Step 2: ログを確認**

Run: `git log --oneline -10`

Expected: 一連のコミットが表示されること

- [ ] **Step 3: バージョンタグを作成**

```bash
git tag -a v1.0.0 -m "Initial release: CLI + Desktop template"
```

- [ ] **Step 4: リモートにプッシュ**

Note: ユーザーが手動でプッシュすること

```bash
git push origin main
git push origin v1.0.0
```

---

## Summary

実装完了後、以下の機能が利用可能になります:

| 機能 | コマンド |
|------|----------|
| CLI実行 | `make run-cli` |
| CLIビルド | `make build-cli` |
| デスクトップ実行 | `make run-desktop` |
| デスクトップビルド | `make build-desktop` |
| テスト | `make test` |
| カバレッジ | `make test-coverage` |
| モック生成 | `make mocks` |
| フォーマット | `make fmt` |
| リンター | `make lint` |
| クリーンアップ | `make clean` |
