package LandcLogFace

import (
	"context"
	"time"
)

// LogLevel 定义日志级别
type LogLevel int

const (
	// DebugLevel 调试级别
	DebugLevel LogLevel = iota
	// InfoLevel 信息级别
	InfoLevel
	// WarnLevel 警告级别
	WarnLevel
	// ErrorLevel 错误级别
	ErrorLevel
	// FatalLevel 致命级别
	FatalLevel
	// PanicLevel 恐慌级别
	PanicLevel
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	default:
		return "UNKNOWN"
	}
}

// Field 定义日志字段
type Field struct {
	Key   string
	Value interface{}
}

// Logger 日志门面接口
type Logger interface {
	// SetLevel 设置日志级别
	SetLevel(level LogLevel)
	// GetLevel 获取当前日志级别
	GetLevel() LogLevel

	// Debug 输出调试级日志
	Debug(msg string, fields ...Field)
	// Debugf 输出格式化的调试级日志
	Debugf(format string, args ...interface{})

	// Info 输出信息级日志
	Info(msg string, fields ...Field)
	// Infof 输出格式化的信息级日志
	Infof(format string, args ...interface{})

	// Warn 输出警告级日志
	Warn(msg string, fields ...Field)
	// Warnf 输出格式化的警告级日志
	Warnf(format string, args ...interface{})

	// Error 输出错误级日志
	Error(msg string, fields ...Field)
	// Errorf 输出格式化的错误级日志
	Errorf(format string, args ...interface{})

	// Fatal 输出致命级日志并退出程序
	Fatal(msg string, fields ...Field)
	// Fatalf 输出格式化的致命级日志并退出程序
	Fatalf(format string, args ...interface{})

	// Panic 输出恐慌级日志并触发panic
	Panic(msg string, fields ...Field)
	// Panicf 输出格式化的恐慌级日志并触发panic
	Panicf(format string, args ...interface{})

	// WithFields 添加字段到日志
	WithFields(fields ...Field) Logger
	// WithField 添加单个字段到日志
	WithField(key string, value interface{}) Logger
	// WithContext 添加上下文到日志
	WithContext(ctx context.Context) Logger
	// WithError 添加错误信息到日志
	WithError(err error) Logger
	// WithTime 添加时间到日志
	WithTime(t time.Time) Logger

	// IsDebugEnabled 检查调试级别是否启用
	IsDebugEnabled() bool
	// IsInfoEnabled 检查信息级别是否启用
	IsInfoEnabled() bool
	// IsWarnEnabled 检查警告级别是否启用
	IsWarnEnabled() bool
	// IsErrorEnabled 检查错误级别是否启用
	IsErrorEnabled() bool
	// IsFatalEnabled 检查致命级别是否启用
	IsFatalEnabled() bool
	// IsPanicEnabled 检查恐慌级别是否启用
	IsPanicEnabled() bool

	// Sync 刷新日志缓冲区
	Sync() error
}

// LoggerProvider 日志提供者接口
type LoggerProvider interface {
	// Create 创建日志实例
	Create(name string) Logger
	// CreateWithConfig 根据配置创建日志实例
	CreateWithConfig(name string, config map[string]interface{}) Logger
}

// Option 日志配置选项
type Option func(*LoggerOptions)

// LoggerOptions 日志配置选项
type LoggerOptions struct {
	Level      LogLevel
	Format     string
	OutputPath string
	Config     map[string]interface{}
}

// WithLevel 设置日志级别
func WithLevel(level LogLevel) Option {
	return func(opt *LoggerOptions) {
		opt.Level = level
	}
}

// WithFormat 设置日志格式
func WithFormat(format string) Option {
	return func(opt *LoggerOptions) {
		opt.Format = format
	}
}

// WithOutputPath 设置日志输出路径
func WithOutputPath(path string) Option {
	return func(opt *LoggerOptions) {
		opt.OutputPath = path
	}
}

// WithConfig 设置额外配置
func WithConfig(config map[string]interface{}) Option {
	return func(opt *LoggerOptions) {
		opt.Config = config
	}
}
