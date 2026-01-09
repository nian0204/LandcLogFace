package LandcLogFace

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

// StdLogger 标准库log适配器
type StdLogger struct {
	level    LogLevel
	fields   []Field
	ctx      context.Context
	logger   *log.Logger
	name     string
}

// NewStdLogger 创建标准库log实例
func NewStdLogger(name string, opts ...Option) *StdLogger {
	options := &LoggerOptions{
		Level:      InfoLevel,
		Format:     "text",
		OutputPath: "stdout",
		Config:     make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(options)
	}

	// 配置输出
	output := os.Stdout
	if options.OutputPath != "stdout" {
		file, err := os.OpenFile(options.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			output = file
		}
	}

	// 创建标准库log实例
	logger := log.New(output, "", log.LstdFlags)

	return &StdLogger{
		level:  options.Level,
		fields: make([]Field, 0),
		ctx:    context.Background(),
		logger: logger,
		name:   name,
	}
}

// SetLevel 设置日志级别
func (s *StdLogger) SetLevel(level LogLevel) {
	s.level = level
}

// GetLevel 获取当前日志级别
func (s *StdLogger) GetLevel() LogLevel {
	return s.level
}

// formatMessage 格式化日志消息
func (s *StdLogger) formatMessage(level LogLevel, msg string, fields []Field) string {
	allFields := append(s.fields, fields...)

	fieldStr := ""
	for _, field := range allFields {
		fieldStr += fmt.Sprintf(" %s=%v", field.Key, field.Value)
	}

	return fmt.Sprintf("[%s] [%s] %s%s", level.String(), s.name, msg, fieldStr)
}

// Debug 输出调试级日志
func (s *StdLogger) Debug(msg string, fields ...Field) {
	if s.level <= DebugLevel {
		s.logger.Println(s.formatMessage(DebugLevel, msg, fields))
	}
}

// Debugf 输出格式化的调试级日志
func (s *StdLogger) Debugf(format string, args ...interface{}) {
	if s.level <= DebugLevel {
		msg := fmt.Sprintf(format, args...)
		s.logger.Println(s.formatMessage(DebugLevel, msg, nil))
	}
}

// Info 输出信息级日志
func (s *StdLogger) Info(msg string, fields ...Field) {
	if s.level <= InfoLevel {
		s.logger.Println(s.formatMessage(InfoLevel, msg, fields))
	}
}

// Infof 输出格式化的信息级日志
func (s *StdLogger) Infof(format string, args ...interface{}) {
	if s.level <= InfoLevel {
		msg := fmt.Sprintf(format, args...)
		s.logger.Println(s.formatMessage(InfoLevel, msg, nil))
	}
}

// Warn 输出警告级日志
func (s *StdLogger) Warn(msg string, fields ...Field) {
	if s.level <= WarnLevel {
		s.logger.Println(s.formatMessage(WarnLevel, msg, fields))
	}
}

// Warnf 输出格式化的警告级日志
func (s *StdLogger) Warnf(format string, args ...interface{}) {
	if s.level <= WarnLevel {
		msg := fmt.Sprintf(format, args...)
		s.logger.Println(s.formatMessage(WarnLevel, msg, nil))
	}
}

// Error 输出错误级日志
func (s *StdLogger) Error(msg string, fields ...Field) {
	if s.level <= ErrorLevel {
		s.logger.Println(s.formatMessage(ErrorLevel, msg, fields))
	}
}

// Errorf 输出格式化的错误级日志
func (s *StdLogger) Errorf(format string, args ...interface{}) {
	if s.level <= ErrorLevel {
		msg := fmt.Sprintf(format, args...)
		s.logger.Println(s.formatMessage(ErrorLevel, msg, nil))
	}
}

// Fatal 输出致命级日志并退出程序
func (s *StdLogger) Fatal(msg string, fields ...Field) {
	if s.level <= FatalLevel {
		s.logger.Println(s.formatMessage(FatalLevel, msg, fields))
		os.Exit(1)
	}
}

// Fatalf 输出格式化的致命级日志并退出程序
func (s *StdLogger) Fatalf(format string, args ...interface{}) {
	if s.level <= FatalLevel {
		msg := fmt.Sprintf(format, args...)
		s.logger.Println(s.formatMessage(FatalLevel, msg, nil))
		os.Exit(1)
	}
}

// Panic 输出恐慌级日志并触发panic
func (s *StdLogger) Panic(msg string, fields ...Field) {
	if s.level <= PanicLevel {
		msg := s.formatMessage(PanicLevel, msg, fields)
		s.logger.Println(msg)
		panic(msg)
	}
}

// Panicf 输出格式化的恐慌级日志并触发panic
func (s *StdLogger) Panicf(format string, args ...interface{}) {
	if s.level <= PanicLevel {
		msg := fmt.Sprintf(format, args...)
		fullMsg := s.formatMessage(PanicLevel, msg, nil)
		s.logger.Println(fullMsg)
		panic(fullMsg)
	}
}

// WithFields 添加字段到日志
func (s *StdLogger) WithFields(fields ...Field) Logger {
	newLogger := *s
	newLogger.fields = append(newLogger.fields, fields...)
	return &newLogger
}

// WithField 添加单个字段到日志
func (s *StdLogger) WithField(key string, value interface{}) Logger {
	return s.WithFields(Field{Key: key, Value: value})
}

// WithContext 添加上下文到日志
func (s *StdLogger) WithContext(ctx context.Context) Logger {
	newLogger := *s
	newLogger.ctx = ctx
	return &newLogger
}

// WithError 添加错误信息到日志
func (s *StdLogger) WithError(err error) Logger {
	return s.WithField("error", err)
}

// WithTime 添加时间到日志
func (s *StdLogger) WithTime(t time.Time) Logger {
	return s.WithField("time", t)
}

// IsDebugEnabled 检查调试级别是否启用
func (s *StdLogger) IsDebugEnabled() bool {
	return s.level <= DebugLevel
}

// IsInfoEnabled 检查信息级别是否启用
func (s *StdLogger) IsInfoEnabled() bool {
	return s.level <= InfoLevel
}

// IsWarnEnabled 检查警告级别是否启用
func (s *StdLogger) IsWarnEnabled() bool {
	return s.level <= WarnLevel
}

// IsErrorEnabled 检查错误级别是否启用
func (s *StdLogger) IsErrorEnabled() bool {
	return s.level <= ErrorLevel
}

// IsFatalEnabled 检查致命级别是否启用
func (s *StdLogger) IsFatalEnabled() bool {
	return s.level <= FatalLevel
}

// IsPanicEnabled 检查恐慌级别是否启用
func (s *StdLogger) IsPanicEnabled() bool {
	return s.level <= PanicLevel
}

// Sync 刷新日志缓冲区
func (s *StdLogger) Sync() error {
	// 标准库log没有Sync方法，返回nil
	return nil
}

// StdLoggerProvider 标准库log提供者
type StdLoggerProvider struct{}

// NewStdLoggerProvider 创建标准库log提供者
func NewStdLoggerProvider() *StdLoggerProvider {
	return &StdLoggerProvider{}
}

// Create 创建日志实例
func (p *StdLoggerProvider) Create(name string) Logger {
	return NewStdLogger(name)
}

// CreateWithConfig 根据配置创建日志实例
func (p *StdLoggerProvider) CreateWithConfig(name string, config map[string]interface{}) Logger {
	var level LogLevel
	if lvl, ok := config["level"].(LogLevel); ok {
		level = lvl
	} else {
		level = InfoLevel
	}

	var format string
	if fmt, ok := config["format"].(string); ok {
		format = fmt
	} else {
		format = "text"
	}

	var outputPath string
	if path, ok := config["outputPath"].(string); ok {
		outputPath = path
	} else {
		outputPath = "stdout"
	}

	return NewStdLogger(name, 
		WithLevel(level),
		WithFormat(format),
		WithOutputPath(outputPath),
		WithConfig(config),
	)
}
