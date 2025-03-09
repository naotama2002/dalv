package main

import (
	"fmt"
	"os"

	"github.com/naotama2002/dalv/internal/cli"
	"github.com/naotama2002/dalv/internal/duckdb"
	"github.com/naotama2002/dalv/internal/validator"
	"github.com/naotama2002/dalv/pkg/utils"
)

func main() {
	// ロガーの初期化
	logger := utils.NewLogger(utils.INFO)
	logger.Info("dalvを起動しています...")

	// コマンドライン引数の解析
	cliParser := cli.NewCLI(os.Args[1:])
	s3Path, tableName, err := cliParser.Parse()
	if err != nil {
		logger.Error("コマンドライン引数の解析に失敗しました: %v", err)
		os.Exit(1)
	}

	// ヘルプまたはバージョン表示の場合は終了
	if s3Path == "" {
		os.Exit(0)
	}

	// S3パスの検証
	pathValidator := validator.NewS3PathValidator()
	if err := pathValidator.ValidateS3Path(s3Path); err != nil {
		logger.Error("S3パスの検証に失敗しました: %v", err)
		os.Exit(1)
	}

	// DuckDBのインストール確認
	if err := pathValidator.ValidateDuckDBInstallation(); err != nil {
		logger.Error("DuckDBの検証に失敗しました: %v", err)
		fmt.Println("\nDuckDBがインストールされていないようです。")
		fmt.Println("インストール方法: https://duckdb.org/docs/installation/")
		os.Exit(1)
	}

	// DuckDBの実行
	executor := duckdb.NewExecutor()
	if err := executor.CheckDuckDBInstallation(); err != nil {
		logger.Error("DuckDBの検証に失敗しました: %v", err)
		fmt.Println("\nDuckDBがインストールされていないようです。")
		fmt.Println("インストール方法: https://duckdb.org/docs/installation/")
		os.Exit(1)
	}

	logger.Info("DuckDBを起動しています...")
	logger.Info("S3パス: %s", s3Path)
	if tableName != "" {
		logger.Info("テーブル名: %s", tableName)
	}

	if err := executor.ExecuteDuckDB(s3Path, tableName); err != nil {
		logger.Error("DuckDBの実行に失敗しました: %v", err)
		os.Exit(1)
	}

	logger.Info("正常に終了しました")
}
