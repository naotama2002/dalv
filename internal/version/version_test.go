package version

import (
	"os"
	"testing"
)

func TestGetVersion(t *testing.T) {
	// 元のバージョンを保存
	originalVersion := Version
	defer func() {
		// テスト後に元のバージョンを復元
		Version = originalVersion
	}()

	// テスト用にバージョンを設定
	Version = "1.0.0"

	if got := GetVersion(); got != "1.0.0" {
		t.Errorf("GetVersion() = %v, want %v", got, "1.0.0")
	}
}

func TestReadVersionFromFile(t *testing.T) {
	// テスト用の一時ファイルを作成
	tmpFile, err := os.CreateTemp("", "version-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// テスト用のバージョン情報をファイルに書き込む
	testVersion := "2.0.0"
	if _, err := tmpFile.WriteString(testVersion); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// テスト実行
	got, err := ReadVersionFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadVersionFromFile() error = %v", err)
	}
	if got != testVersion {
		t.Errorf("ReadVersionFromFile() = %v, want %v", got, testVersion)
	}

	// 存在しないファイルの場合のテスト
	_, err = ReadVersionFromFile("non-existent-file.txt")
	if err == nil {
		t.Error("ReadVersionFromFile() with non-existent file should return error")
	}
}

func TestVersionString(t *testing.T) {
	// 元のバージョンを保存
	originalVersion := Version
	defer func() {
		// テスト後に元のバージョンを復元
		Version = originalVersion
	}()

	// テスト用にバージョンを設定
	Version = "3.0.0"

	expected := "dalv バージョン 3.0.0"
	if got := VersionString(); got != expected {
		t.Errorf("VersionString() = %v, want %v", got, expected)
	}
}
