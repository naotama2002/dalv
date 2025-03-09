package version

import (
	"fmt"
	"os"
	"strings"
)

var (
	// Version はビルド時に -ldflags で上書きされる可能性があります
	Version = "dev"
)

// GetVersion はバージョン情報を返します
func GetVersion() string {
	return Version
}

// ReadVersionFromFile はファイルからバージョン情報を読み込みます
func ReadVersionFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// VersionString はバージョン情報を整形して返します
func VersionString() string {
	return fmt.Sprintf("dalv バージョン %s", Version)
}
