package main

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/naotama2002/dalv/internal/version"
)

// モック用の構造体と関数
type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return m.sys }

func TestInitVersion(t *testing.T) {
	// 元の関数を保存
	origOsExecutable := osExecutable
	origOsStat := osStat
	origReadVersionFromFile := readVersionFromFile
	origVersion := version.Version

	// テスト終了後に元に戻す
	defer func() {
		osExecutable = origOsExecutable
		osStat = origOsStat
		readVersionFromFile = origReadVersionFromFile
		version.Version = origVersion
	}()

	// テストケース
	tests := []struct {
		name            string
		executablePath  string
		executableErr   error
		fileExists      bool
		fileStatErr     error
		versionContent  string
		versionErr      error
		expectedVersion string
	}{
		{
			name:            "正常系: バージョンファイルが存在する場合",
			executablePath:  "/path/to/bin/dalv",
			executableErr:   nil,
			fileExists:      true,
			fileStatErr:     nil,
			versionContent:  "1.0.0",
			versionErr:      nil,
			expectedVersion: "1.0.0",
		},
		{
			name:            "異常系: 実行ファイルのパス取得に失敗した場合",
			executablePath:  "",
			executableErr:   errors.New("executable error"),
			fileExists:      false,
			fileStatErr:     nil,
			versionContent:  "",
			versionErr:      nil,
			expectedVersion: "dev", // 変更されないはず
		},
		{
			name:            "異常系: バージョンファイルが存在しない場合（開発環境）",
			executablePath:  "/path/to/bin/dalv",
			executableErr:   nil,
			fileExists:      false,
			fileStatErr:     os.ErrNotExist,
			versionContent:  "2.0.0",
			versionErr:      nil,
			expectedVersion: "2.0.0", // 開発環境のパスから読み込まれる
		},
		{
			name:            "異常系: バージョンファイルの読み込みに失敗した場合",
			executablePath:  "/path/to/bin/dalv",
			executableErr:   nil,
			fileExists:      true,
			fileStatErr:     nil,
			versionContent:  "",
			versionErr:      errors.New("read error"),
			expectedVersion: "dev", // 変更されないはず
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// バージョンをリセット
			version.Version = "dev"

			// モック関数を設定
			osExecutable = func() (string, error) {
				return tt.executablePath, tt.executableErr
			}

			osStat = func(name string) (os.FileInfo, error) {
				if tt.fileExists {
					return &mockFileInfo{}, nil
				}
				return nil, tt.fileStatErr
			}

			readVersionFromFile = func(filePath string) (string, error) {
				return tt.versionContent, tt.versionErr
			}

			// テスト対象の関数を実行
			initVersion()

			// 結果を検証
			if version.Version != tt.expectedVersion {
				t.Errorf("initVersion() set version to %v, want %v", version.Version, tt.expectedVersion)
			}
		})
	}
}
