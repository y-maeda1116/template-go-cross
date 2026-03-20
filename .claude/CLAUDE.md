# Go Template Project

## プロジェクト概要

Mac、Windows、Linux対応のクロスプラットフォームGoアプリケーションテンプレート。

## プロジェクト構造

```
.
├── .github/workflows/    # GitHub Actions CI/CD
├── .claude/             # Claude Code 設定
├── cmd/                 # アプリケーションエントリーポイント
│   └── app/
│       └── main.go     # メインアプリケーション（signal 待機付き）
├── internal/            # 内部パッケージ（外部からインポート不可）
│   └── config/         # 環境変数設定
├── bin/                # ビルド出力ディレクトリ
├── .gitignore          # Go/Windows/Mac 対応
├── go.mod              # Goモジュール定義
├── go.sum              # 依存関係ロック
├── Makefile            # クロスプラットフォームビルド
├── README.md           # プロジェクトドキュメント
├── LICENSE             # MIT ライセンス
└── env.example         # 環境変数テンプレート
```

## 開発ルール

### 依存関係管理

- `github.com/joho/godotenv` を使用して `.env` ファイルから環境変数を読み込む
- 依存パッケージは `go mod tidy` で管理する

### ビルド

```bash
# 現在のOS向け
make build

# Windows (amd64)
make build-win

# Mac (arm64/Apple Silicon)
make build-mac

# 実行
make run

# クリーンアップ
make clean
```

### テスト

```bash
go test ./...
go test -v -race ./...
```

### コーディング規約

- 標準的なGoプロジェクト構造（`cmd/`, `internal/`）を維持する
- アプリケーションコードは `cmd/app/` に配置
- 内部パッケージは `internal/` 以下に配置
- 外部からインポート不要なパッケージは `internal/` に配置
- signal処理を使用して Ctrl+C で安全に終了させる

### GitHub Actions

- Ubuntu環境でテストを実行
- レースコンディション検出（`-race`）
- カバレッジレポートのアップロード
