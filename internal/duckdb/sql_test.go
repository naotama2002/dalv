package duckdb

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNewSQLGenerator(t *testing.T) {
	generator := NewSQLGenerator()
	if generator == nil {
		t.Fatal("NewSQLGenerator returned nil")
	}
}

func TestGenerateAWSConfigSQL(t *testing.T) {
	generator := NewSQLGenerator()
	sql := generator.GenerateAWSConfigSQL()

	// 必要な要素が含まれているか確認
	requiredElements := []string{
		"INSTALL aws",
		"LOAD aws",
		"INSTALL httpfs",
		"LOAD httpfs",
		"CREATE SECRET",
		"TYPE S3",
		"PROVIDER CREDENTIAL_CHAIN",
	}

	for _, element := range requiredElements {
		if !strings.Contains(sql, element) {
			t.Errorf("Generated AWS config SQL does not contain '%s'", element)
		}
	}
}

func TestGenerateCreateTableSQL(t *testing.T) {
	generator := NewSQLGenerator()
	tableName := "test_table"
	s3Path := "s3://bucket/path/to/logs/*.log.gz"
	
	sql := generator.GenerateCreateTableSQL(tableName, s3Path)

	// 必要な要素が含まれているか確認
	if !strings.Contains(sql, "CREATE TABLE test_table") {
		t.Error("Generated SQL does not contain CREATE TABLE statement with correct table name")
	}

	if !strings.Contains(sql, s3Path) {
		t.Error("Generated SQL does not contain the provided S3 path")
	}

	// ALBログの重要なカラムが含まれているか確認
	requiredColumns := []string{
		"'type': 'VARCHAR'",
		"'timestamp': 'TIMESTAMP'",
		"'elb': 'VARCHAR'",
		"'client_ip_port': 'VARCHAR'",
		"'target_ip_port': 'VARCHAR'",
		"'request_processing_time': 'DOUBLE'",
		"'elb_status_code': 'INTEGER'",
		"'request': 'VARCHAR'",
		"'user_agent': 'VARCHAR'",
	}

	for _, column := range requiredColumns {
		if !strings.Contains(sql, column) {
			t.Errorf("Generated SQL does not contain column definition '%s'", column)
		}
	}

	// 必要なCSV読み込みオプションが含まれているか確認
	requiredOptions := []string{
		"delim=' '",
		"quote='\"'",
		"escape='\"'",
		"header=False",
		"auto_detect=False",
	}

	for _, option := range requiredOptions {
		if !strings.Contains(sql, option) {
			t.Errorf("Generated SQL does not contain CSV option '%s'", option)
		}
	}
}

func TestGenerateCompleteSQL(t *testing.T) {
	generator := NewSQLGenerator()
	tableName := "test_table"
	s3Path := "s3://bucket/path/to/logs/*.log.gz"
	
	sql := generator.GenerateCompleteSQL(s3Path, tableName)

	// AWS設定SQLが含まれているか確認
	if !strings.Contains(sql, "INSTALL aws") {
		t.Error("Complete SQL does not contain AWS configuration")
	}

	// テーブル作成SQLが含まれているか確認
	if !strings.Contains(sql, "CREATE TABLE test_table") {
		t.Error("Complete SQL does not contain table creation statement")
	}

	// インタラクティブモードのメッセージが含まれているか確認
	if !strings.Contains(sql, "SELECT 'ALBログが正常にロードされました") {
		t.Error("Complete SQL does not contain interactive mode message")
	}

	if !strings.Contains(sql, "test_table") {
		t.Error("Complete SQL does not contain the table name in the message")
	}
}

func TestGenerateCompleteSQL_WithEmptyTableName(t *testing.T) {
	generator := NewSQLGenerator()
	s3Path := "s3://bucket/path/to/logs/*.log.gz"
	
	// テーブル名を空にして呼び出し
	sql := generator.GenerateCompleteSQL(s3Path, "")

	// 自動生成されたテーブル名のパターンを確認
	if !strings.Contains(sql, "alb_logs_") {
		t.Error("Complete SQL does not contain auto-generated table name with 'alb_logs_' prefix")
	}

	// 日付形式のパターンが含まれているか確認
	datePattern := time.Now().Format("20060102")
	if !strings.Contains(sql, datePattern) {
		t.Error("Complete SQL does not contain current date in the auto-generated table name")
	}
}

func TestGenerateTableName(t *testing.T) {
	generator := NewSQLGenerator()
	tableName := generator.generateTableName()

	// テーブル名のプレフィックスを確認
	if !strings.HasPrefix(tableName, "alb_logs_") {
		t.Errorf("Generated table name '%s' does not have 'alb_logs_' prefix", tableName)
	}

	// 日付形式が含まれているか確認
	datePattern := time.Now().Format("20060102")
	if !strings.Contains(tableName, datePattern) {
		t.Errorf("Generated table name '%s' does not contain current date '%s'", tableName, datePattern)
	}

	// 時刻形式が含まれているか確認
	currentHour := time.Now().Hour()
	hourStr := fmt.Sprintf("%02d", currentHour)
	
	if !strings.Contains(tableName, hourStr) {
		t.Errorf("Generated table name '%s' does not contain current hour '%s'", tableName, hourStr)
	}
}
