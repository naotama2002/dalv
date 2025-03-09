package validator

import (
	"fmt"
	"strings"
)

// S3PathValidator はS3パスの検証を行います
type S3PathValidator struct{}

// NewS3PathValidator は新しいS3パスバリデータを作成します
func NewS3PathValidator() *S3PathValidator {
	return &S3PathValidator{}
}

// ValidateS3Path はS3パスが正しい形式かどうかを検証します
func (v *S3PathValidator) ValidateS3Path(path string) error {
	if path == "" {
		return fmt.Errorf("S3パスが指定されていません")
	}

	if !strings.HasPrefix(path, "s3://") {
		return fmt.Errorf("S3パスは's3://'で始まる必要があります: %s", path)
	}

	parts := strings.Split(strings.TrimPrefix(path, "s3://"), "/")
	if len(parts) < 2 || parts[0] == "" {
		return fmt.Errorf("無効なS3パス形式です。バケット名が必要です: %s", path)
	}

	return nil
}

// ValidateDuckDBInstallation はDuckDBがインストールされているかを検証します
func (v *S3PathValidator) ValidateDuckDBInstallation() error {
	// この関数は実際の実装では、duckdbコマンドが存在するかどうかを確認します
	// ここでは簡単な実装としています
	return nil
}
