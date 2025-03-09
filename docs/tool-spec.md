# dalv ツール仕様書

## 概要

`dalv`は、AWS Application Load Balancer (ALB) のS3ログファイルをDuckDBを使って簡単にクエリするためのコマンドラインツールです。S3バケットパスを指定して、DuckDBを使ってALBログを分析できるようにします。

## 技術スタック

- 言語: Go 1.22
- 主要依存関係:
  - DuckDB: システムにインストールされたバージョン（外部コマンドとして利用）
  - AWS SDK for Go

## 機能要件

### 1. コマンドライン引数

```
dalv [options] <s3-path>
```

#### 引数

- `<s3-path>`: 必須。AWS ALBログが保存されているS3パス (例: `s3://{S3_BUCKET_NAME}/xxxxxx/AWSLogs/{ACCOUNT_ID}/elasticloadbalancing/{REGION}/2025/03/03/*.log.gz`)

#### オプション

- `-t, --table <name>`: 作成するテーブル名 (デフォルト: 自動生成)
- `-h, --help`: ヘルプ情報を表示
- `-v, --version`: バージョン情報を表示

### 2. 初期化処理

1. AWS認証情報の確認と設定
   - 環境変数、AWS設定ファイル、IAMロールから認証情報を自動検出
   - 認証情報が見つからない場合はエラーメッセージを表示

2. DuckDBの初期化
   - DuckDBインスタンスの作成
   - 必要な拡張機能のインストールと読み込み
     - `aws`拡張機能
     - `httpfs`拡張機能
   - AWS認証情報をDuckDBに設定

3. S3パスの検証
   - 指定されたS3パスの形式を検証
   - パスが存在するか確認
   - アクセス権限の確認

### 3. テーブル作成

1. S3ログファイルからテーブルを自動作成
   - ALBログのスキーマに合わせたカラム定義
   - 指定されたS3パスからデータを読み込み
   - テーブル名は引数で指定されたものか、指定がない場合は日付などから自動生成

### 4. インタラクティブコンソール

1. DuckDBのインタラクティブコンソールの提供
   - SQLクエリの入力と実行
   - クエリ結果の表示
   - テーブル一覧の表示コマンド
   - スキーマ情報の表示コマンド
   - ヘルプコマンド

2. コンソール機能
   - コマンド履歴
   - タブ補完
   - シンタックスハイライト
   - 結果の整形表示（テーブル形式）

3. 特殊コマンド
   - `.tables`: 利用可能なテーブル一覧を表示
   - `.schema <table>`: 指定テーブルのスキーマを表示
   - `.export <file> <format>`: クエリ結果をファイルにエクスポート
   - `.quit` または `.exit`: コンソールを終了

### 3. エラーハンドリング

1. システムエラー
   - DuckDBがインストールされていない場合
   - DuckDBのバージョンが互換性がない場合

2. 入力エラー
   - S3パスの形式が不正な場合
   - 必須パラメータが不足している場合

3. 実行時エラー
   - DuckDBの実行に失敗した場合
   - AWS認証情報の問題

## 非機能要件

### 1. パフォーマンス

- DuckDBの処理能力を最大限に活用
- 効率的なSQLクエリの生成

### 2. セキュリティ

- AWS認証情報の安全な取り扱い
- 認証情報のログやコンソールへの出力禁止

### 3. ユーザビリティ

- 明確なエラーメッセージ
- インストール手順の明確化
- 使用例の提供

## プロジェクト構造

```
dalv/
├── .github/
│   └── workflows/
│       └── ci.yml
├── cmd/
│   └── dalv/
│       └── main.go        # エントリーポイント
├── internal/
│   ├── cli/
│   │   └── cli.go         # コマンドライン引数の処理
│   ├── duckdb/
│   │   ├── executor.go    # DuckDB実行ロジック
│   │   └── sql.go         # SQL生成ロジック
│   ├── validator/
│   │   └── validator.go   # 入力検証
│   └── schema/
│       └── alb.go         # ALBログスキーマ定義
├── pkg/
│   └── utils/
│       └── logger.go      # ロギングユーティリティ
├── go.mod                 # Goモジュール定義
├── go.sum                 # 依存関係ハッシュ
├── LICENSE
└── README.md
```

## 開発ロードマップ

### フェーズ1: 基本機能実装

1. プロジェクト構造のセットアップ
2. コマンドライン引数の処理
3. AWS認証と接続
4. DuckDB初期化と基本クエリ

### フェーズ2: インタラクティブコンソール

1. コンソールインターフェース
2. コマンド履歴と補完
3. 特殊コマンド実装

### フェーズ3: 高度な機能

1. 出力フォーマット
2. パフォーマンス最適化
3. エラーハンドリングの強化

### フェーズ4: テストと文書化

1. 単体テスト
2. 統合テスト
3. ユーザードキュメント

## 使用例

### 基本的な使用方法

```bash
# S3パスを指定して起動
dalv s3://my-bucket/logs/AWSLogs/123456789012/elasticloadbalancing/us-east-1/2025/03/03/*.log.gz

# テーブル名を指定
dalv -t alb_logs_march s3://my-bucket/logs/AWSLogs/123456789012/elasticloadbalancing/us-east-1/2025/03/03/*.log.gz

# 結果をファイルに出力
dalv -o results.csv s3://my-bucket/logs/AWSLogs/123456789012/elasticloadbalancing/us-east-1/2025/03/03/*.log.gz
```

### インタラクティブコンソールでの操作

```sql
-- テーブル一覧の表示
.tables

-- スキーマの表示
.schema alb_logs

-- クエリの実行
SELECT elb_status_code, COUNT(*) as count 
FROM alb_logs 
GROUP BY elb_status_code 
ORDER BY count DESC;

-- 結果のエクスポート
.export results.json json

-- 終了
.exit
```

## 制限事項

1. 現在のバージョンでは、ALBログ形式のみサポート
2. 一度に処理できるログサイズに制限あり（マシンのメモリに依存）
3. 複雑なJOIN操作やウィンドウ関数は、大量データ処理時にパフォーマンスが低下する可能性あり
