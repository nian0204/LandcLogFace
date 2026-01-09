package LandcLogFace

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

// ConsoleLogger 默认的控制台日志适配器
type ConsoleLogger struct {
	level    LogLevel
	fields   []Field
	ctx      context.Context
	logger   *log.Logger
	name     string
}

// NewConsoleLogger 创建控制台日志实例
func NewConsoleLogger(name string, opts ...Option) *ConsoleLogger {
	options := &LoggerOptions{
		Level:      InfoLevel,
		Format:     "text",
		OutputPath: "stdout",
		Config:     make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(options)
	}

	output := os.Stdout
	if options.OutputPath != "stdout" {
		file, err := os.OpenFile(options.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			output = file
		}
	}

	return &ConsoleLogger{
		level:  options.Level,
		fields: make([]Field, 0),
		ctx:    context.Background(),
		logger: log.New(output, "", 0),
		name:   name,
	}
}

// SetLevel 设置日志级别
func (c *ConsoleLogger) SetLevel(level LogLevel) {
	c.level = level
}

// GetLevel 获取当前日志级别
func (c *ConsoleLogger) GetLevel() LogLevel {
	return c.level
}

// formatMessage 格式化日志消息
func (c *ConsoleLogger) formatMessage(level LogLevel, msg string, fields []Field) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	allFields := append(c.fields, fields...)

	fieldStr := ""
	for _, field := range allFields {
		fieldStr += fmt.Sprintf(" %s=%v", field.Key, field.Value)
	}

	return fmt.Sprintf("%s [%s] [%s] %s%s", timestamp, level.String(), c.name, msg, fieldStr)
}

// Debug 输出调试级日志
func (c *ConsoleLogger) Debug(msg string, fields ...Field) {
	if c.level <= DebugLevel {
		c.logger.Println(c.formatMessage(DebugLevel, msg, fields))
	}
}

// Debugf 输出格式化的调试级日志
func (c *ConsoleLogger) Debugf(format string, args ...interface{}) {
	if c.level <= DebugLevel {
		msg := fmt.Sprintf(format, args...)
		c.logger.Println(c.formatMessage(DebugLevel, msg, nil))
	}
}

// Info 输出信息级日志
func (c *ConsoleLogger) Info(msg string, fields ...Field) {
	if c.level <= InfoLevel {
		c.logger.Println(c.formatMessage(InfoLevel, msg, fields))
	}
}

// Infof 输出格式化的信息级日志
func (c *ConsoleLogger) Infof(format string, args ...interface{}) {
	if c.level <= InfoLevel {
		msg := fmt.Sprintf(format, args...)
		c.logger.Println(c.formatMessage(InfoLevel, msg, nil))
	}
}

// Warn 输出警告级日志
func (c *ConsoleLogger) Warn(msg string, fields ...Field) {
	if c.level <= WarnLevel {
		c.logger.Println(c.formatMessage(WarnLevel, msg, fields))
	}
}

// Warnf 输出格式化的警告级日志
func (c *ConsoleLogger) Warnf(format string, args ...interface{}) {
	if c.level <= WarnLevel {
		msg := fmt.Sprintf(format, args...)
		c.logger.Println(c.formatMessage(WarnLevel, msg, nil))
	}
}

// Error 输出错误级日志
func (c *ConsoleLogger) Error(msg string, fields ...Field) {
	if c.level <= ErrorLevel {
		c.logger.Println(c.formatMessage(ErrorLevel, msg, fields))
	}
}

// Errorf 输出格式化的错误级日志
func (c *ConsoleLogger) Errorf(format string, args ...interface{}) {
	if c.level <= ErrorLevel {
		msg := fmt.Sprintf(format, args...)
		c.logger.Println(c.formatMessage(ErrorLevel, msg, nil))
	}
}

// Fatal 输出致命级日志并退出程序
func (c *ConsoleLogger) Fatal(msg string, fields ...Field) {
	if c.level <= FatalLevel {
		c.logger.Println(c.formatMessage(FatalLevel, msg, fields))
		os.Exit(1)
	}
}

// Fatalf 输出格式化的致命级日志并退出程序
func (c *ConsoleLogger) Fatalf(format string, args ...interface{}) {
	if c.level <= FatalLevel {
		msg := fmt.Sprintf(format, args...)
		c.logger.Println(c.formatMessage(FatalLevel, msg, nil))
		os.Exit(1)
	}
}

// Panic 输出恐慌级日志并触发panic
func (c *ConsoleLogger) Panic(msg string, fields ...Field) {
	if c.level <= PanicLevel {
		msg := c.formatMessage(PanicLevel, msg, fields)
		c.logger.Println(msg)
		panic(msg)
	}
}

// Panicf 输出格式化的恐慌级日志并触发panic
func (c *ConsoleLogger) Panicf(format string, args ...interface{}) {
	if c.level <= PanicLevel {
		msg := fmt.Sprintf(format, args...)
		fullMsg := c.formatMessage(PanicLevel, msg, nil)
		c.logger.Println(fullMsg)
		panic(fullMsg)
	}
}

// WithFields 添加字段到日志
func (c *ConsoleLogger) WithFields(fields ...Field) Logger {
	newLogger := *c
	newLogger.fields = append(newLogger.fields, fields...)
	return &newLogger
}

// WithField 添加单个字段到日志
func (c *ConsoleLogger) WithField(key string, value interface{}) Logger {
	return c.WithFields(Field{Key: key, Value: value})
}

// WithContext 添加上下文到日志
func (c *ConsoleLogger) WithContext(ctx context.Context) Logger {
	newLogger := *c
	newLogger.ctx = ctx
	return &newLogger
}

// WithError 添加错误信息到日志
func (c *ConsoleLogger) WithError(err error) Logger {
	return c.WithField("error", err)
}

// WithTime 添加时间到日志
func (c *ConsoleLogger) WithTime(t time.Time) Logger {
	return c.WithField("time", t)
}

// IsDebugEnabled 检查调试级别是否启用
func (c *ConsoleLogger) IsDebugEnabled() bool {
	return c.level <= DebugLevel
}

// IsInfoEnabled 检查信息级别是否启用
func (c *ConsoleLogger) IsInfoEnabled() bool {
	return c.level <= InfoLevel
}

// IsWarnEnabled 检查警告级别是否启用
func (c *ConsoleLogger) IsWarnEnabled() bool {
	return c.level <= WarnLevel
}

// IsErrorEnabled 检查错误级别是否启用
func (c *ConsoleLogger) IsErrorEnabled() bool {
	return c.level <= ErrorLevel
}

// IsFatalEnabled 检查致命级别是否启用
func (c *ConsoleLogger) IsFatalEnabled() bool {
	return c.level <= FatalLevel
}

// IsPanicEnabled 检查恐慌级别是否启用
func (c *ConsoleLogger) IsPanicEnabled() bool {
	return c.level <= PanicLevel
}

// Sync 刷新日志缓冲区
func (c *ConsoleLogger) Sync() error {
	return nil
}

// ConsoleLoggerProvider 控制台日志提供者
type ConsoleLoggerProvider struct{}

// NewConsoleLoggerProvider 创建控制台日志提供者
func NewConsoleLoggerProvider() *ConsoleLoggerProvider {
	return &ConsoleLoggerProvider{}
}

// Create 创建日志实例
func (p *ConsoleLoggerProvider) Create(name string) Logger {
	return NewConsoleLogger(name)
}

// CreateWithConfig 根据配置创建日志实例
func (p *ConsoleLoggerProvider) CreateWithConfig(name string, config map[string]interface{}) Logger {
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

	return NewConsoleLogger(name, 
		WithLevel(level),
		WithFormat(format),
		WithOutputPath(outputPath),
		WithConfig(config),
	)
}
