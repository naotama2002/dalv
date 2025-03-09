package duckdb

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Executor はDuckDBを実行するための構造体です
type Executor struct {
	sqlGenerator *SQLGenerator
}

// NewExecutor は新しいDuckDB実行者を作成します
func NewExecutor() *Executor {
	return &Executor{
		sqlGenerator: NewSQLGenerator(),
	}
}

// ExecuteDuckDB はDuckDBを実行します
func (e *Executor) ExecuteDuckDB(s3Path string, tableName string) error {
	// テーブル名が指定されていない場合は生成
	if tableName == "" {
		tableName = e.sqlGenerator.generateTableName()
	}

	// SQLを生成
	sql := e.sqlGenerator.GenerateCompleteSQL(s3Path, tableName)

	// 一時ファイルを作成
	tempDir, err := ioutil.TempDir("", "dalv-")
	if err != nil {
		return fmt.Errorf("一時ディレクトリの作成に失敗しました: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// SQLを一時ファイルに書き込み
	sqlFilePath := filepath.Join(tempDir, "init.sql")
	err = ioutil.WriteFile(sqlFilePath, []byte(sql), 0644)
	if err != nil {
		return fmt.Errorf("SQLファイルの作成に失敗しました: %w", err)
	}

	// DuckDBコマンドを実行
	cmd := exec.Command("duckdb", "-init", sqlFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// コマンドを実行
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("DuckDBの実行に失敗しました: %w", err)
	}

	return nil
}

// CheckDuckDBInstallation はDuckDBがインストールされているかを確認します
func (e *Executor) CheckDuckDBInstallation() error {
	cmd := exec.Command("duckdb", "--version")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("DuckDBがインストールされていないか、実行できません: %w", err)
	}
	return nil
}
