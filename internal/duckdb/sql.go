package duckdb

import (
	"fmt"
	"time"
)

// SQLGenerator はDuckDBのSQLを生成するための構造体です
type SQLGenerator struct{}

// NewSQLGenerator は新しいSQLジェネレーターを作成します
func NewSQLGenerator() *SQLGenerator {
	return &SQLGenerator{}
}

// GenerateCompleteSQL はS3パスからデータを読み込むための完全なSQLを生成します
func (g *SQLGenerator) GenerateCompleteSQL(s3Path string, tableName string) string {
	// テーブル名が指定されていない場合は生成
	if tableName == "" {
		tableName = g.generateTableName()
	}

	// AWS認証設定のSQL
	awsConfigSQL := g.GenerateAWSConfigSQL()

	// ALBログのスキーマを定義し、S3からデータを読み込むSQL
	createTableSQL := g.GenerateCreateTableSQL(tableName, s3Path)

	// 完全なSQLを結合
	return fmt.Sprintf("%s\n\n%s\n\n-- インタラクティブモードのためのメッセージ\nSELECT 'ALBログが正常にロードされました。以下のテーブルに対してクエリを実行できます: %s' AS message;\n", awsConfigSQL, createTableSQL, tableName)
}

// GenerateAWSConfigSQL はAWS認証情報を設定するSQLを生成します
func (g *SQLGenerator) GenerateAWSConfigSQL() string {
	return `-- AWS拡張機能のインストールと認証設定
INSTALL aws;
LOAD aws;
INSTALL httpfs;
LOAD httpfs;
CREATE SECRET secret_s3 (
    TYPE S3,
    PROVIDER CREDENTIAL_CHAIN
);`
}

// GenerateCreateTableSQL はALBログのテーブルを作成するSQLを生成します
func (g *SQLGenerator) GenerateCreateTableSQL(tableName string, s3Path string) string {
	return fmt.Sprintf(`-- ALBログのテーブルを作成
CREATE TABLE %s AS
SELECT *
FROM read_csv(
    '%s',
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
);`, tableName, s3Path)
}

// generateTableName は一意のテーブル名を生成します
func (g *SQLGenerator) generateTableName() string {
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("alb_logs_%s", timestamp)
}
