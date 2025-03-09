package utils

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger(INFO)
	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}

	if logger.level != INFO {
		t.Errorf("Expected log level to be INFO, got %v", logger.level)
	}

	if logger.logger == nil {
		t.Fatal("Logger's internal logger is nil")
	}
}

func TestLogLevels(t *testing.T) {
	// バッファを使用してログ出力をキャプチャ
	var buf bytes.Buffer
	
	// テスト用のロガーを作成
	logger := &Logger{
		level:  INFO,
		logger: log.New(&buf, "", 0),
	}

	// INFOレベルのログ出力をテスト
	logger.Info("This is an info message: %s", "test")
	output := buf.String()
	if !strings.Contains(output, "INFO") || !strings.Contains(output, "This is an info message: test") {
		t.Errorf("Info log output incorrect, got: %s", output)
	}
	buf.Reset()

	// DEBUGレベルのログは出力されないことをテスト
	logger.Debug("This is a debug message: %s", "test")
	output = buf.String()
	if output != "" {
		t.Errorf("Debug log should not be output at INFO level, got: %s", output)
	}
	buf.Reset()

	// ERRORレベルのログ出力をテスト
	logger.Error("This is an error message: %s", "test")
	output = buf.String()
	if !strings.Contains(output, "ERROR") || !strings.Contains(output, "This is an error message: test") {
		t.Errorf("Error log output incorrect, got: %s", output)
	}
	buf.Reset()

	// ログレベルをDEBUGに変更
	logger.level = DEBUG
	
	// DEBUGレベルのログが出力されることをテスト
	logger.Debug("This is a debug message: %s", "test")
	output = buf.String()
	if !strings.Contains(output, "DEBUG") || !strings.Contains(output, "This is a debug message: test") {
		t.Errorf("Debug log output incorrect, got: %s", output)
	}
	buf.Reset()
}

func TestLoggerWithDifferentLevels(t *testing.T) {
	testCases := []struct {
		level         LogLevel
		debugVisible  bool
		infoVisible   bool
		warnVisible   bool
		errorVisible  bool
	}{
		{DEBUG, true, true, true, true},
		{INFO, false, true, true, true},
		{WARN, false, false, true, true},
		{ERROR, false, false, false, true},
	}

	for _, tc := range testCases {
		var buf bytes.Buffer
		logger := &Logger{
			level:  tc.level,
			logger: log.New(&buf, "", 0),
		}

		// DEBUGメッセージ
		buf.Reset()
		logger.Debug("DEBUG")
		if (buf.Len() > 0) != tc.debugVisible {
			t.Errorf("At level %v, DEBUG message visibility incorrect, expected %v", tc.level, tc.debugVisible)
		}

		// INFOメッセージ
		buf.Reset()
		logger.Info("INFO")
		if (buf.Len() > 0) != tc.infoVisible {
			t.Errorf("At level %v, INFO message visibility incorrect, expected %v", tc.level, tc.infoVisible)
		}

		// ERRORメッセージ
		buf.Reset()
		logger.Error("ERROR")
		if (buf.Len() > 0) != tc.errorVisible {
			t.Errorf("At level %v, ERROR message visibility incorrect, expected %v", tc.level, tc.errorVisible)
		}
	}
}
