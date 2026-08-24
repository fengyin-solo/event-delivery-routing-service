// Package logger 提供线程安全的分级日志器。
package logger

import (
	"log"
	"os"
	"sync"
)

// Level 表示日志级别。
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// ParseLevel 从字符串解析日志级别，默认 info。
func ParseLevel(s string) Level {
	switch s {
	case "debug":
		return LevelDebug
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

// String 返回级别字符串。
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "info"
	}
}

// Logger 线程安全分级日志器。
type Logger struct {
	mu    sync.Mutex
	level Level
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
}

// New 创建默认 info 级别的日志器。
func New() *Logger { return NewLevel(LevelInfo) }

// NewLevel 创建指定级别的日志器。
func NewLevel(level Level) *Logger {
	flag := log.LstdFlags | log.Lmicroseconds
	return &Logger{
		level: level,
		debug: log.New(os.Stderr, "[DEBUG] ", flag),
		info:  log.New(os.Stderr, "[INFO ] ", flag),
		warn:  log.New(os.Stderr, "[WARN ] ", flag),
		err:   log.New(os.Stderr, "[ERROR] ", flag),
	}
}

// SetLevel 设置日志级别。
func (l *Logger) SetLevel(level Level) { l.mu.Lock(); defer l.mu.Unlock(); l.level = level }

// Debugf 输出 debug 日志。
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level <= LevelDebug {
		l.debug.Printf(format, args...)
	}
}

// Infof 输出 info 日志。
func (l *Logger) Infof(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level <= LevelInfo {
		l.info.Printf(format, args...)
	}
}

// Warnf 输出 warn 日志。
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level <= LevelWarn {
		l.warn.Printf(format, args...)
	}
}

// Errorf 输出 error 日志。
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level <= LevelError {
		l.err.Printf(format, args...)
	}
}

// Printf 等价于 Infof。
func (l *Logger) Printf(format string, args ...interface{}) { l.Infof(format, args...) }
