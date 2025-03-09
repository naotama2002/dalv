package validator

import (
	"strings"
	"testing"
)

func TestNewS3PathValidator(t *testing.T) {
	validator := NewS3PathValidator()
	if validator == nil {
		t.Fatal("NewS3PathValidator returned nil")
	}
}

func TestValidateS3Path_ValidPath(t *testing.T) {
	validator := NewS3PathValidator()
	
	// 有効なS3パスのテストケース
	validPaths := []string{
		"s3://bucket/path",
		"s3://bucket/path/to/logs",
		"s3://bucket/path/to/logs/*.log.gz",
		"s3://my-bucket-name/AWSLogs/123456789012/elasticloadbalancing/ap-northeast-1/2025/03/03/*.log.gz",
	}
	
	for _, path := range validPaths {
		err := validator.ValidateS3Path(path)
		if err != nil {
			t.Errorf("ValidateS3Path failed for valid path '%s': %v", path, err)
		}
	}
}

func TestValidateS3Path_InvalidPath(t *testing.T) {
	validator := NewS3PathValidator()
	
	// 無効なS3パスのテストケース
	invalidPaths := []struct {
		path string
		expectedErrContains string
	}{
		{"", "S3パスが指定されていません"},
		{"bucket/path", "S3パスは's3://'で始まる必要があります"},
		{"s3://", "無効なS3パス形式です。バケット名が必要です"},
		{"s3:///path", "無効なS3パス形式です。バケット名が必要です"},
		{"http://bucket/path", "S3パスは's3://'で始まる必要があります"},
	}
	
	for _, tc := range invalidPaths {
		err := validator.ValidateS3Path(tc.path)
		if err == nil {
			t.Errorf("ValidateS3Path should fail for invalid path '%s'", tc.path)
			continue
		}
		
		if !contains(err.Error(), tc.expectedErrContains) {
			t.Errorf("Error message for path '%s' does not contain expected text. Got: '%s', Expected to contain: '%s'", 
				tc.path, err.Error(), tc.expectedErrContains)
		}
	}
}

func TestValidateDuckDBInstallation(t *testing.T) {
	validator := NewS3PathValidator()
	
	// この関数は実際の実装では、duckdbコマンドが存在するかどうかを確認します
	// テスト環境では常にnilを返すことを確認
	err := validator.ValidateDuckDBInstallation()
	if err != nil {
		t.Errorf("ValidateDuckDBInstallation failed: %v", err)
	}
}

// contains は文字列sに部分文字列subが含まれるかどうかを確認するヘルパー関数
func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
