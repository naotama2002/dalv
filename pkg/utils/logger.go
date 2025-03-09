package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel はログレベルを表します
type LogLevel int

const (
	// DEBUG はデバッグレベルのログです
	DEBUG LogLevel = iota
	// INFO は情報レベルのログです
	INFO
	// WARN は警告レベルのログです
	WARN
	// ERROR はエラーレベルのログです
	ERROR
)

// Logger はロギング機能を提供します
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

// NewLogger は新しいロガーを作成します
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stderr, "", 0),
	}
}

// Debug はデバッグレベルのログを出力します
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.log("DEBUG", format, v...)
	}
}

// Info は情報レベルのログを出力します
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.log("INFO", format, v...)
	}
}

// Warn は警告レベルのログを出力します
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WARN {
		l.log("WARN", format, v...)
	}
}

// Error はエラーレベルのログを出力します
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.log("ERROR", format, v...)
	}
}

// log は実際にログを出力します
func (l *Logger) log(level string, format string, v ...interface{}) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	message := fmt.Sprintf(format, v...)
	l.logger.Printf("[%s] %s: %s", timestamp, level, message)
}
