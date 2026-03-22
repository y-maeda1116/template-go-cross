# CLI + Desktop Go Template Design

**Date:** 2026-03-23
**Author:** Claude
**Status:** Draft

## Overview

単一リポジトリ内でCLIとデスクトップアプリケーション（Wails）を共存させるGoテンプレート。

## Requirements

| 項目 | 選択 |
|------|------|
| アーキテクチャ | 単一テンプレート（CLI + デスクトップ共存） |
| デスクトップUI | Wails v2 |
| CLI | Cobra |
| ロギング | zap |
| 設定管理 | Viper (YAML/TOML/JSON) |
| ホットリロード | Air |
| テスト | testify, mockgen |
| テストカバレッジ | 80%+ |

## Project Structure

```
.
├── cmd/
│   ├── cli/
│   │   └── main.go           # CLIエントリーポイント (Cobra)
│   └── desktop/
│       └── main.go           # Wailsエントリーポイント
├── internal/
│   ├── core/                 # CLIとデスクトップで共有するビジネスロジック
│   │   ├── service.go
│   │   └── service_test.go
│   ├── ui/                   # Wails UI コンポーネント
│   │   ├── app.go
│   │   └── models/
│   ├── config/               # 設定管理 (Viper)
│   │   ├── config.go
│   │   └── config_test.go
│   ├── logger/               # ロギング (zap)
│   │   ├── logger.go
│   │   └── logger_test.go
│   └── cli/                  # CLIコマンド定義
│       ├── root.go
│       ├── cmd1.go
│       └── cmd2.go
├── frontend/                 # Wailsフロントエンド
│   ├── src/
│   ├── package.json
│   └── wailsjs/
├── pkg/                      # 外部から再利用可能なパッケージ
│   └── utils/
├── test/
│   ├── mocks/                # mockgenで生成されたモック
│   └── testutil/             # テストユーティリティ
├── wails.json                # Wails設定
├── air.toml                  # ホットリロード設定
├── Makefile                  # 拡張されたビルドターゲット
├── go.mod
├── go.sum
├── .env.example              # 環境変数テンプレート
├── config.yaml.example       # 設定ファイルテンプレート
└── .github/workflows/
    ├── test.yml              # テスト CI
    ├── build-cli.yml         # CLIビルド CI
    └── build-desktop.yml     # デスクトップビルド CI
```

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Application Layer                    │
├───────────────────┬─────────────────────────────────────┤
│   CLI (Cobra)     │   Desktop (Wails)                  │
│                   │                                     │
│  cmd/cli/main.go  │  cmd/desktop/main.go               │
│  internal/cli/    │  frontend/* (React/Vue)            │
│                   │  internal/ui/ (Go bindings)        │
└─────────┬─────────┴─────────────────────────────────────┘
          │
          ▼
┌─────────────────────────────────────────────────────────┐
│                Core Business Logic Layer                │
│                                                         │
│  internal/core/ (CLIとデスクトップ両方から呼ばれる)      │
│  - Service interfaces                                   │
│  - Business logic implementations                      │
│  - Domain models                                       │
└─────────┬───────────────────────────────────────────────┘
          │
          ▼
┌─────────────────────────────────────────────────────────┐
│              Infrastructure Layer                        │
│                                                         │
│  internal/config/ (Viper)                               │
│  internal/logger/ (zap)                                 │
│  pkg/utils/                                             │
└─────────────────────────────────────────────────────────┘
```

## Dependencies

```go
require (
    // CLI
    github.com/spf13/cobra v1.8.1

    // 設定
    github.com/spf13/viper v1.19.0

    // ロギング
    go.uber.org/zap v1.27.0
    go.uber.org/zap/zapcore v1.27.0

    // デスクトップ
    github.com/wailsapp/wails/v2 v2.9.1

    // テスト
    github.com/stretchr/testify v1.9.0
    golang.org/x/mock v1.6.0
)
```

## Configuration

**環境変数 (`.env`):**
```bash
APP_NAME=myapp
APP_ENV=development
APP_PORT=8080
LOG_LEVEL=info
CONFIG_PATH=config.yaml
```

**設定ファイル (`config.yaml`):**
```yaml
app:
  name: myapp
  version: 1.0.0
  environment: development

server:
  host: localhost
  port: 8080

logging:
  level: info
  format: json
  output: stdout

cli:
  theme: auto

desktop:
  window:
    width: 1200
    height: 800
    resizable: true
```

## Error Handling

**カスタムエラー型:**

```go
// internal/core/errors.go
type AppError struct {
    Code    string
    Message string
    Err     error
}
```

**エラーハンドリング方針:**

| レイヤー | 役割 |
|---------|------|
| Core | `AppError`でラップして返す |
| CLI | `AppError`をキャッチしてユーザーフレンドリーなメッセージを表示 |
| Desktop | `AppError`をキャッチしてUI通知に変換 |
| Logger | 全レイヤーのエラーをログに記録 |

## Testing

| テストタイプ | 対象 | ツール | カバレッジ目標 |
|-------------|------|-------|--------------|
| Unit Tests | core/, internal/*/ | testify, mockgen | 80%+ |
| Integration Tests | CLIコマンド, Config | testify | 主要パス |
| E2E Tests | Wailsアプリ全体 | Playwright (オプション) | 主要ユースケース |

## Build & Dev

**Makefileターゲット:**

| ターゲット | 説明 |
|-----------|------|
| `build-cli` | CLIをビルド |
| `build-desktop` | デスクトップをビルド |
| `build-all` | 両方をビルド |
| `run-cli` | CLIを実行 |
| `run-desktop` | デスクトップをdevモードで実行 |
| `test` | 全テストを実行 |
| `test-coverage` | カバレッジ付きテスト |
| `test-race` | レース検出付きテスト |
| `mocks` | モックを生成 |
| `fmt` | コードフォーマット |
| `lint` | リンター実行 |
| `clean` | クリーンアップ |

## Implementation Notes

- CoreビジネスロジックはCLIとデスクトップ両方から使用される
- 設定の読み込み優先順位: コマンドライン > 環境変数 > 設定ファイル > デフォルト値
- ログレベル: debug, info, warn, error
- モックは`mockgen`で自動生成
- ホットリロードはAirを使用
