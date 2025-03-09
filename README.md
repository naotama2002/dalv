# dalv

AWS Application Load Balancer (ALB) のS3ログファイルをDuckDBを使って簡単にクエリするためのコマンドラインツールです。

## 概要

`dalv`は、DuckDBを外部コマンドとして利用し、AWS ALBのS3ログを簡単にクエリするためのツールです。S3パスを指定するだけで、必要なセットアップを自動的に行い、DuckDBのインタラクティブコンソールを起動します。

## 前提条件

- Go 1.22以上
- DuckDB（コマンドラインツール）がインストールされていること
- AWS認証情報が設定されていること（環境変数、AWS設定ファイル、またはIAMロール）

## インストール

### Makeを使用したビルドとインストール

```bash
# ソースの取得
git clone https://github.com/naotama2002/dalv.git
cd dalv

# ビルド（bin/dalvが生成されます）
make build

# テストの実行
make test

# リントチェック
make lint

# 全てのステップを実行（クリーン、リント、テスト、ビルド）
make

# インストール（$GOPATH/binにインストールされます）
make install
```

### Goコマンドを直接使用

```bash
# ソースから直接ビルド
git clone https://github.com/naotama2002/dalv.git
cd dalv
go build -o dalv ./cmd/dalv

# インストール（オプション）
go install github.com/naotama2002/dalv/cmd/dalv@latest
```

## 使用方法

```bash
# 基本的な使用方法
dalv "s3://{S3_BUCKET_NAME}/xxxxxx/AWSLogs/{ACCOUNT_ID}/elasticloadbalancing/{REGION}/2025/03/03/*.log.gz"

# テーブル名を指定して実行
dalv -t my_alb_logs "s3://{S3_BUCKET_NAME}/xxxxxx/AWSLogs/{ACCOUNT_ID}/elasticloadbalancing/{REGION}/2025/03/**/*.log.gz"

# ヘルプの表示
dalv -h

# バージョン情報の表示
dalv -v
```

## 動作の仕組み

`dalv`は以下の処理を自動的に行います：

1. DuckDBの初期化
   ```sql
   INSTALL aws;
   LOAD aws;
   INSTALL httpfs;
   LOAD httpfs;
   CREATE SECRET (
       TYPE S3,
       PROVIDER CREDENTIAL_CHAIN
   );
   ```

2. ALBログテーブルの作成
   ```sql
   CREATE TABLE alb_log_20250303 AS
   SELECT *
   FROM read_csv(
       's3://{S3_BUCKET_NAME}/xxxxxx/AWSLogs/{ACCOUNT_ID}/elasticloadbalancing/{REGION}/2025/03/03/*.log.gz',
       columns={
           'type': 'VARCHAR',
           'timestamp': 'TIMESTAMP',
           'elb': 'VARCHAR',
           'client_ip_port': 'VARCHAR',
           'target_ip_port': 'VARCHAR',
           'request_processing_time': 'DOUBLE',
           'target_processing_time': 'DOUBLE',
           'response_processing_time': 'DOUBLE',
           'elb_status_code': 'INTEGER',
           'target_status_code': 'VARCHAR',
           'received_bytes': 'BIGINT',
           'sent_bytes': 'BIGINT',
           'request': 'VARCHAR',
           'user_agent': 'VARCHAR',
           'ssl_cipher': 'VARCHAR',
           'ssl_protocol': 'VARCHAR',
           'target_group_arn': 'VARCHAR',
           'trace_id': 'VARCHAR',
           'domain_name': 'VARCHAR',
           'chosen_cert_arn': 'VARCHAR',
           'matched_rule_priority': 'VARCHAR',
           'request_creation_time': 'TIMESTAMP',
           'actions_executed': 'VARCHAR',
           'redirect_url': 'VARCHAR',
           'error_reason': 'VARCHAR',
           'target_port_list': 'VARCHAR',
           'target_status_code_list': 'VARCHAR',
           'classification': 'VARCHAR',
           'classification_reason': 'VARCHAR',
           'conn_trace_id': 'VARCHAR'
       },
       delim=' ',
       quote='"',
       escape='"',
       header=False,
       auto_detect=False
   );
   ```

3. DuckDBのインタラクティブコンソールの起動

## クエリ例

```sql
-- ステータスコードが431以外のリクエストを表示
SELECT * FROM alb_log_20250303 WHERE elb_status_code != 431;

-- 時間帯別のリクエスト数
SELECT 
    DATE_TRUNC('hour', timestamp) AS hour, 
    COUNT(*) AS request_count 
FROM alb_log_20250303 
GROUP BY hour 
ORDER BY hour;

-- ステータスコード別のリクエスト数
SELECT 
    elb_status_code, 
    COUNT(*) AS count 
FROM alb_log_20250303 
GROUP BY elb_status_code 
ORDER BY count DESC;
```

## ライセンス

MITライセンス
