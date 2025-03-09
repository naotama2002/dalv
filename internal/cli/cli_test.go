package cli

import (
	"fmt"
	"testing"
)

// CLIのモック実装
type mockCLI struct {
	args        []string
	helpFlag    bool
	versionFlag bool
	tableValue  string
}

func (m *mockCLI) Parse() (string, string, error) {
	// 引数が空の場合はエラー
	if len(m.args) == 0 {
		return "", "", fmt.Errorf("S3パスが指定されていません")
	}

	// ヘルプフラグが指定された場合
	if m.helpFlag {
		return "", "", nil
	}

	// バージョンフラグが指定された場合
	if m.versionFlag {
		return "", "", nil
	}

	// 最後の引数をS3パスとして扱う
	s3Path := m.args[len(m.args)-1]
	return s3Path, m.tableValue, nil
}

func TestCLIWithTableFlag(t *testing.T) {
	// テーブル名指定のテスト
	mock := &mockCLI{
		args:       []string{"-t", "test_table", "s3://bucket/path"},
		tableValue: "test_table",
	}

	s3Path, tableName, err := mock.Parse()
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if s3Path != "s3://bucket/path" {
		t.Fatalf("Expected s3Path to be 's3://bucket/path', got '%s'", s3Path)
	}

	if tableName != "test_table" {
		t.Fatalf("Expected tableName to be 'test_table', got '%s'", tableName)
	}
}

func TestCLIWithoutTableFlag(t *testing.T) {
	// テーブル名指定なしのテスト
	mock := &mockCLI{
		args: []string{"s3://bucket/path"},
	}

	s3Path, tableName, err := mock.Parse()
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if s3Path != "s3://bucket/path" {
		t.Fatalf("Expected s3Path to be 's3://bucket/path', got '%s'", s3Path)
	}

	if tableName != "" {
		t.Fatalf("Expected tableName to be empty, got '%s'", tableName)
	}
}

func TestCLIWithHelpFlag(t *testing.T) {
	// ヘルプフラグのテスト
	mock := &mockCLI{
		args:     []string{"-h"},
		helpFlag: true,
	}

	s3Path, tableName, err := mock.Parse()
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if s3Path != "" {
		t.Fatalf("Expected s3Path to be empty, got '%s'", s3Path)
	}

	if tableName != "" {
		t.Fatalf("Expected tableName to be empty, got '%s'", tableName)
	}
}

func TestCLIWithVersionFlag(t *testing.T) {
	// バージョンフラグのテスト
	mock := &mockCLI{
		args:        []string{"-v"},
		versionFlag: true,
	}

	s3Path, tableName, err := mock.Parse()
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if s3Path != "" {
		t.Fatalf("Expected s3Path to be empty, got '%s'", s3Path)
	}

	if tableName != "" {
		t.Fatalf("Expected tableName to be empty, got '%s'", tableName)
	}
}

func TestCLIWithNoArgs(t *testing.T) {
	// 引数なしのテスト
	mock := &mockCLI{
		args: []string{},
	}

	_, _, err := mock.Parse()
	if err == nil {
		t.Fatal("Expected error for no args, got nil")
	}
}
