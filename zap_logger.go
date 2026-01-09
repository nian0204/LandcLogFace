package LandcLogFace

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger zap日志库适配器
type ZapLogger struct {
	logger *zap.Logger
	level  LogLevel
	fields []Field
	ctx    context.Context
	name   string
}

// NewZapLogger 创建zap日志实例
func NewZapLogger(name string, opts ...Option) *ZapLogger {
	options := &LoggerOptions{
		Level:      InfoLevel,
		Format:     "json",
		OutputPath: "stdout",
		Config:     make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(options)
	}

	// 配置zap
	zapLevel := zapcore.InfoLevel
	switch options.Level {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case InfoLevel:
		zapLevel = zapcore.InfoLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	case FatalLevel:
		zapLevel = zapcore.FatalLevel
	case PanicLevel:
		zapLevel = zapcore.PanicLevel
	}

	// 配置输出
	var outputPaths []string
	if options.OutputPath == "stdout" {
		outputPaths = []string{"stdout"}
	} else {
		outputPaths = []string{options.OutputPath}
	}

	// 创建zap配置
	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLevel),
		Development: false,
		Encoding:    options.Format,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      outputPaths,
		ErrorOutputPaths: []string{"stderr"},
	}

	// 构建logger
	logger, err := zapConfig.Build()
	if err != nil {
		// 如果构建失败，使用默认配置
		logger, _ = zap.NewProduction()
	}

	// 添加名称字段
	logger = logger.Named(name)

	return &ZapLogger{
		logger: logger,
		level:  options.Level,
		fields: make([]Field, 0),
		ctx:    context.Background(),
		name:   name,
	}
}

// SetLevel 设置日志级别
func (z *ZapLogger) SetLevel(level LogLevel) {
	z.level = level
	// 更新zap的日志级别
	// 注意：zap的AtomicLevel需要通过Core来设置，这里简化处理
}

// GetLevel 获取当前日志级别
func (z *ZapLogger) GetLevel() LogLevel {
	return z.level
}

// toZapFields 将自定义字段转换为zap字段
func (z *ZapLogger) toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(z.fields)+len(fields))

	// 添加已有的字段
	for _, field := range z.fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}

	// 添加新的字段
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}

	return zapFields
}

// Debug 输出调试级日志
func (z *ZapLogger) Debug(msg string, fields ...Field) {
	if z.level <= DebugLevel {
		z.logger.Debug(msg, z.toZapFields(fields)...)
	}
}

// Debugf 输出格式化的调试级日志
func (z *ZapLogger) Debugf(format string, args ...interface{}) {
	if z.level <= DebugLevel {
		z.logger.Sugar().Debugf(format, args...)
	}
}

// Info 输出信息级日志
func (z *ZapLogger) Info(msg string, fields ...Field) {
	if z.level <= InfoLevel {
		z.logger.Info(msg, z.toZapFields(fields)...)
	}
}

// Infof 输出格式化的信息级日志
func (z *ZapLogger) Infof(format string, args ...interface{}) {
	if z.level <= InfoLevel {
		z.logger.Sugar().Infof(format, args...)
	}
}

// Warn 输出警告级日志
func (z *ZapLogger) Warn(msg string, fields ...Field) {
	if z.level <= WarnLevel {
		z.logger.Warn(msg, z.toZapFields(fields)...)
	}
}

// Warnf 输出格式化的警告级日志
func (z *ZapLogger) Warnf(format string, args ...interface{}) {
	if z.level <= WarnLevel {
		z.logger.Sugar().Warnf(format, args...)
	}
}

// Error 输出错误级日志
func (z *ZapLogger) Error(msg string, fields ...Field) {
	if z.level <= ErrorLevel {
		z.logger.Error(msg, z.toZapFields(fields)...)
	}
}

// Errorf 输出格式化的错误级日志
func (z *ZapLogger) Errorf(format string, args ...interface{}) {
	if z.level <= ErrorLevel {
		z.logger.Sugar().Errorf(format, args...)
	}
}

// Fatal 输出致命级日志并退出程序
func (z *ZapLogger) Fatal(msg string, fields ...Field) {
	if z.level <= FatalLevel {
		z.logger.Fatal(msg, z.toZapFields(fields)...)
		os.Exit(1)
	}
}

// Fatalf 输出格式化的致命级日志并退出程序
func (z *ZapLogger) Fatalf(format string, args ...interface{}) {
	if z.level <= FatalLevel {
		z.logger.Sugar().Fatalf(format, args...)
		os.Exit(1)
	}
}

// Panic 输出恐慌级日志并触发panic
func (z *ZapLogger) Panic(msg string, fields ...Field) {
	if z.level <= PanicLevel {
		z.logger.Panic(msg, z.toZapFields(fields)...)
	}
}

// Panicf 输出格式化的恐慌级日志并触发panic
func (z *ZapLogger) Panicf(format string, args ...interface{}) {
	if z.level <= PanicLevel {
		z.logger.Sugar().Panicf(format, args...)
	}
}

// WithFields 添加字段到日志
func (z *ZapLogger) WithFields(fields ...Field) Logger {
	newLogger := *z
	newLogger.fields = append(newLogger.fields, fields...)
	// 创建新的zap.Logger实例
	zapFields := make([]zap.Field, 0, len(newLogger.fields))
	for _, field := range newLogger.fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	newLogger.logger = z.logger.With(zapFields...)
	return &newLogger
}

// WithField 添加单个字段到日志
func (z *ZapLogger) WithField(key string, value interface{}) Logger {
	return z.WithFields(Field{Key: key, Value: value})
}

// WithContext 添加上下文到日志
func (z *ZapLogger) WithContext(ctx context.Context) Logger {
	newLogger := *z
	newLogger.ctx = ctx
	return &newLogger
}

// WithError 添加错误信息到日志
func (z *ZapLogger) WithError(err error) Logger {
	return z.WithField("error", err)
}

// WithTime 添加时间到日志
func (z *ZapLogger) WithTime(t time.Time) Logger {
	return z.WithField("time", t)
}

// IsDebugEnabled 检查调试级别是否启用
func (z *ZapLogger) IsDebugEnabled() bool {
	return z.level <= DebugLevel
}

// IsInfoEnabled 检查信息级别是否启用
func (z *ZapLogger) IsInfoEnabled() bool {
	return z.level <= InfoLevel
}

// IsWarnEnabled 检查警告级别是否启用
func (z *ZapLogger) IsWarnEnabled() bool {
	return z.level <= WarnLevel
}

// IsErrorEnabled 检查错误级别是否启用
func (z *ZapLogger) IsErrorEnabled() bool {
	return z.level <= ErrorLevel
}

// IsFatalEnabled 检查致命级别是否启用
func (z *ZapLogger) IsFatalEnabled() bool {
	return z.level <= FatalLevel
}

// IsPanicEnabled 检查恐慌级别是否启用
func (z *ZapLogger) IsPanicEnabled() bool {
	return z.level <= PanicLevel
}

// Sync 刷新日志缓冲区
func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}

// ZapLoggerProvider zap日志提供者
type ZapLoggerProvider struct{}

// NewZapLoggerProvider 创建zap日志提供者
func NewZapLoggerProvider() *ZapLoggerProvider {
	return &ZapLoggerProvider{}
}

// Create 创建日志实例
func (p *ZapLoggerProvider) Create(name string) Logger {
	return NewZapLogger(name)
}

// CreateWithConfig 根据配置创建日志实例
func (p *ZapLoggerProvider) CreateWithConfig(name string, config map[string]interface{}) Logger {
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
		format = "json"
	}

	var outputPath string
	if path, ok := config["outputPath"].(string); ok {
		outputPath = path
	} else {
		outputPath = "stdout"
	}

	return NewZapLogger(name, 
		WithLevel(level),
		WithFormat(format),
		WithOutputPath(outputPath),
		WithConfig(config),
	)
}
