package duckdb

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestNewExecutor(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("NewExecutor returned nil")
	}

	if executor.sqlGenerator == nil {
		t.Fatal("Executor's sqlGenerator is nil")
	}
}

// モックコマンド実行関数を使用するためのヘルパー関数
type mockCmd struct {
	cmd      string
	args     []string
	exitCode int
	executed bool
}

func (m *mockCmd) execCommand(command string, args ...string) *exec.Cmd {
	m.cmd = command
	m.args = args
	m.executed = true

	// 実際のコマンドは実行せず、モックコマンドを返す
	if m.exitCode != 0 {
		// 失敗するコマンドをシミュレート
		return exec.Command("false")
	}
	// 成功するコマンドをシミュレート
	return exec.Command("true")
}

func TestCheckDuckDBInstallation(t *testing.T) {
	// DuckDBがインストールされているかどうかを確認
	// 実際のコマンド実行をモックに置き換える

	tests := []struct {
		name     string
		exitCode int
		wantErr  bool
	}{
		{
			name:     "DuckDB is installed",
			exitCode: 0,
			wantErr:  false,
		},
		{
			name:     "DuckDB is not installed",
			exitCode: 1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックコマンド実行関数を設定
			mock := &mockCmd{exitCode: tt.exitCode}

			// 元のexec.Command関数を保存
			origExecCommand := execCommand
			// テスト終了後に元に戻す
			defer func() { execCommand = origExecCommand }()

			// モック関数に置き換え
			execCommand = mock.execCommand

			executor := NewExecutor()
			err := executor.CheckDuckDBInstallation()

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckDuckDBInstallation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !mock.executed {
				t.Error("Command was not executed")
			}

			if mock.cmd != "duckdb" {
				t.Errorf("Expected command 'duckdb', got '%s'", mock.cmd)
			}

			if len(mock.args) != 1 || mock.args[0] != "--version" {
				t.Errorf("Expected args ['--version'], got %v", mock.args)
			}
		})
	}
}

// ExecuteDuckDBのテストは複雑なため、部分的なテストを実装
func TestExecuteDuckDB_GeneratesCorrectSQL(t *testing.T) {
	// 元のioutil.TempDir関数を保存
	origTempDir := tempDirFunc
	// テスト終了後に元に戻す
	defer func() { tempDirFunc = origTempDir }()

	// テスト用の一時ディレクトリを作成
	testTempDir, err := os.MkdirTemp("", "dalv-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testTempDir)

	// モックのTempDir関数を設定
	tempDirFunc = func(dir, prefix string) (string, error) {
		return testTempDir, nil
	}

	// 元のioutil.WriteFile関数を保存
	origWriteFile := writeFileFunc
	// テスト終了後に元に戻す
	defer func() { writeFileFunc = origWriteFile }()

	// SQLの内容を検証するためのモック関数
	var capturedSQL []byte
	writeFileFunc = func(filename string, data []byte, perm os.FileMode) error {
		capturedSQL = data
		return nil
	}

	// 元のexec.Command関数を保存
	origExecCommand := execCommand
	// テスト終了後に元に戻す
	defer func() { execCommand = origExecCommand }()

	// コマンド実行をモック化
	mock := &mockCmd{exitCode: 0}
	execCommand = mock.execCommand

	// テスト実行
	executor := NewExecutor()
	s3Path := "s3://test-bucket/logs/*.log.gz"
	tableName := "test_table"

	// ExecuteDuckDBを呼び出す前に、実際のコマンド実行をスキップするようにする
	err = executor.ExecuteDuckDB(s3Path, tableName)
	if err != nil {
		t.Fatalf("ExecuteDuckDB() error = %v", err)
	}

	// 生成されたSQLの内容を検証
	sqlString := string(capturedSQL)

	// 必要な要素が含まれているか確認
	if !contains(sqlString, "INSTALL aws") {
		t.Error("Generated SQL does not contain AWS installation")
	}

	if !contains(sqlString, "CREATE TABLE test_table") {
		t.Error("Generated SQL does not contain table creation with correct name")
	}

	if !contains(sqlString, s3Path) {
		t.Error("Generated SQL does not contain the provided S3 path")
	}

	// コマンドが正しく実行されたか確認
	if !mock.executed {
		t.Error("DuckDB command was not executed")
	}

	if mock.cmd != "duckdb" {
		t.Errorf("Expected command 'duckdb', got '%s'", mock.cmd)
	}

	if len(mock.args) != 2 || mock.args[0] != "-init" {
		t.Errorf("Expected args ['-init', <filepath>], got %v", mock.args)
	}
}

// テスト用のヘルパー関数
func contains(s, substr string) bool {
	return s != "" && substr != "" && strings.Contains(s, substr)
}
