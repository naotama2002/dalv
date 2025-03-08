# dalv ツール仕様書

## 概要

`dalv`は、AWS Application Load Balancer (ALB) のS3ログファイルをDuckDBを使って簡単にクエリするためのコマンドラインツールです。S3バケットパスを指定して起動し、インタラクティブなDuckDBコンソールを提供します。

## 技術スタック

- 言語: TypeScript
- 実行環境: Deno 2.2.3
- 主要依存関係:
  - DuckDB: 1.2.0
  - AWS SDK for JavaScript v3.764.0 / https://github.com/aws/aws-sdk-js-v3
  - Cliffy (コマンドライン引数パーサー): 最新バージョン

## 機能要件

### 1. コマンドライン引数

```
dalv [options] <s3-path>
```

#### 引数

- `<s3-path>`: 必須。AWS ALBログが保存されているS3パス (例: `s3://{S3_BUCKET_NAME}/xxxxxx/AWSLogs/{ACCOUNT_ID}/elasticloadbalancing/{REGION}/2025/03/03/*.log.gz`)

#### オプション

- `-t, --table <name>`: 作成するテーブル名 (デフォルト: 自動生成)
- `-o, --output <path>`: クエリ結果の出力先ファイルパス
- `-f, --format <format>`: 出力フォーマット (csv, json, parquet) (デフォルト: csv)
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

### 5. エラーハンドリング

1. 接続エラー
   - AWS認証情報の問題
   - S3バケットへのアクセス権限の問題
   - ネットワーク接続の問題

2. クエリエラー
   - SQL構文エラー
   - メモリ不足エラー
   - タイムアウトエラー

3. ファイル形式エラー
   - ログファイルの形式が想定と異なる場合

## 非機能要件

### 1. パフォーマンス

- 大量のログファイル（数GB）を効率的に処理できること
- クエリ実行時のメモリ使用量の最適化
- 長時間実行クエリのサポート

### 2. セキュリティ

- AWS認証情報の安全な取り扱い
- 認証情報のログやコンソールへの出力禁止
- 一時的なファイルの安全な管理

### 3. ユーザビリティ

- 明確なエラーメッセージ
- 進行状況の表示
- 初回使用時のガイダンス表示

## プロジェクト構造

```
dalv/
├── .github/
│   └── workflows/
│       └── ci.yml
├── src/
│   ├── cli/
│   │   ├── args.ts         # コマンドライン引数の処理
│   │   ├── console.ts      # インタラクティブコンソール
│   │   └── commands.ts     # 特殊コマンドの実装
│   ├── duckdb/
│   │   ├── client.ts       # DuckDB接続クライアント
│   │   ├── schema.ts       # ALBログスキーマ定義
│   │   └── queries.ts      # 共通クエリ
│   ├── aws/
│   │   ├── auth.ts         # AWS認証
│   │   └── s3.ts           # S3操作
│   ├── utils/
│   │   ├── logger.ts       # ロギング
│   │   ├── formatter.ts    # 出力フォーマット
│   │   └── validator.ts    # 入力検証
│   ├── types/
│   │   └── alb-log.ts      # ALBログ型定義
│   └── main.ts             # エントリーポイント
├── test/
│   ├── cli.test.ts
│   ├── duckdb.test.ts
│   └── aws.test.ts
├── deno.json              # Deno設定
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
