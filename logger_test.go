package LandcLogFace

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestLoggerInterface 测试Logger接口的基本功能
func TestLoggerInterface(t *testing.T) {
	// 测试控制台日志
	logger := NewConsoleLogger("test")

	// 测试日志级别设置
	logger.SetLevel(DebugLevel)
	if logger.GetLevel() != DebugLevel {
		t.Errorf("Expected level DebugLevel, got %v", logger.GetLevel())
	}

	// 测试日志输出（这里只是测试方法调用，不测试输出内容）
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warn message")
	logger.Error("Error message")

	// 测试格式化日志
	logger.Debugf("Debug: %s", "test")
	logger.Infof("Info: %s", "test")
	logger.Warnf("Warn: %s", "test")
	logger.Errorf("Error: %s", "test")

	// 测试字段
	logger.Info("With fields",
		Field{Key: "key1", Value: "value1"},
		Field{Key: "key2", Value: 123},
	)

	// 测试链式调用
	logger.WithField("chain", "value").Info("Chained logger")

	// 测试上下文
	ctx := context.Background()
	logger.WithContext(ctx).Info("With context")

	// 测试错误
	err := errors.New("test error")
	logger.WithError(err).Error("With error")

	// 测试时间
	now := time.Now()
	logger.WithTime(now).Info("With time")

	// 测试日志级别检查
	if !logger.IsDebugEnabled() {
		t.Error("Expected DebugLevel to be enabled")
	}
	if !logger.IsInfoEnabled() {
		t.Error("Expected InfoLevel to be enabled")
	}
	if !logger.IsWarnEnabled() {
		t.Error("Expected WarnLevel to be enabled")
	}
	if !logger.IsErrorEnabled() {
		t.Error("Expected ErrorLevel to be enabled")
	}

	// 测试Sync
	if err := logger.Sync(); err != nil {
		t.Errorf("Sync failed: %v", err)
	}
}

// TestLogFactory 测试日志工厂
func TestLogFactory(t *testing.T) {
	factory := GetLogFactory()

	// 测试默认提供者
	defaultProvider := factory.GetDefaultProvider()
	if defaultProvider != "console" {
		t.Errorf("Expected default provider to be 'console', got '%s'", defaultProvider)
	}

	// 测试设置默认提供者
	factory.SetDefaultProvider("zap")
	if factory.GetDefaultProvider() != "zap" {
		t.Error("Failed to set default provider")
	}

	// 测试创建日志实例
	logger := factory.CreateLogger("test")
	if logger == nil {
		t.Error("Failed to create logger")
	}

	// 测试使用指定提供者创建日志实例
	zapLogger := factory.CreateLoggerWithProvider("test", "zap")
	if zapLogger == nil {
		t.Error("Failed to create zap logger")
	}

	logrusLogger := factory.CreateLoggerWithProvider("test", "logrus")
	if logrusLogger == nil {
		t.Error("Failed to create logrus logger")
	}

	stdLogger := factory.CreateLoggerWithProvider("test", "std")
	if stdLogger == nil {
		t.Error("Failed to create std logger")
	}

	consoleLogger := factory.CreateLoggerWithProvider("test", "console")
	if consoleLogger == nil {
		t.Error("Failed to create console logger")
	}

	// 测试根据配置创建日志实例
	config := map[string]interface{}{
		"provider":   "zap",
		"level":      DebugLevel,
		"format":     "json",
		"outputPath": "stdout",
	}

	configLogger := factory.CreateLoggerWithConfig("test", config)
	if configLogger == nil {
		t.Error("Failed to create logger with config")
	}

	// 测试注册和注销提供者
	customProvider := &CustomLoggerProvider{}
	factory.RegisterProvider("test-custom", customProvider)

	customLogger, exists := factory.GetProvider("test-custom")
	if !exists || customLogger == nil {
		t.Error("Failed to register custom provider")
	}

	factory.UnregisterProvider("test-custom")
	_, exists = factory.GetProvider("test-custom")
	if exists {
		t.Error("Failed to unregister custom provider")
	}

	// 恢复默认提供者
	factory.SetDefaultProvider("console")
}

// TestGlobalLogger 测试全局日志
func TestGlobalLogger(t *testing.T) {
	// 测试获取全局日志实例
	globalLogger := GetLogger()
	if globalLogger == nil {
		t.Error("Failed to get global logger")
	}

	// 测试全局日志函数
	Debug("Global debug")
	Info("Global info")
	Warn("Global warn")
	Error("Global error")

	// 测试全局格式化函数
	Debugf("Global debug: %s", "test")
	Infof("Global info: %s", "test")
	Warnf("Global warn: %s", "test")
	Errorf("Global error: %s", "test")

	// 测试全局函数带字段
	Info("Global with fields", Field{Key: "key", Value: "value"})
}

// TestZapLogger 测试zap日志适配器
func TestZapLogger(t *testing.T) {
	logger := NewZapLogger("test-zap",
		WithLevel(DebugLevel),
		WithFormat("json"),
	)

	if logger == nil {
		t.Error("Failed to create zap logger")
	}

	// 测试zap日志功能
	logger.Debug("Zap debug")
	logger.Info("Zap info")
	logger.Warn("Zap warn")
	logger.Error("Zap error")

	// 测试字段
	logger.Info("Zap with fields",
		Field{Key: "zap", Value: "test"},
	)

	// 测试Sync
	if err := logger.Sync(); err != nil {
		t.Errorf("Zap sync failed: %v", err)
	}
}

// TestLogrusLogger 测试logrus日志适配器
func TestLogrusLogger(t *testing.T) {
	logger := NewLogrusLogger("test-logrus",
		WithLevel(DebugLevel),
		WithFormat("text"),
	)

	if logger == nil {
		t.Error("Failed to create logrus logger")
	}

	// 测试logrus日志功能
	logger.Debug("Logrus debug")
	logger.Info("Logrus info")
	logger.Warn("Logrus warn")
	logger.Error("Logrus error")

	// 测试字段
	logger.Info("Logrus with fields",
		Field{Key: "logrus", Value: "test"},
	)

	// 测试Sync
	if err := logger.Sync(); err != nil {
		t.Errorf("Logrus sync failed: %v", err)
	}
}

// TestStdLogger 测试标准库日志适配器
func TestStdLogger(t *testing.T) {
	logger := NewStdLogger("test-std",
		WithLevel(DebugLevel),
	)

	if logger == nil {
		t.Error("Failed to create std logger")
	}

	// 测试std日志功能
	logger.Debug("Std debug")
	logger.Info("Std info")
	logger.Warn("Std warn")
	logger.Error("Std error")

	// 测试字段
	logger.Info("Std with fields",
		Field{Key: "std", Value: "test"},
	)

	// 测试Sync
	if err := logger.Sync(); err != nil {
		t.Errorf("Std sync failed: %v", err)
	}
}

// TestConsoleLogger 测试控制台日志适配器
func TestConsoleLogger(t *testing.T) {
	logger := NewConsoleLogger("test-console",
		WithLevel(DebugLevel),
	)

	if logger == nil {
		t.Error("Failed to create console logger")
	}

	// 测试console日志功能
	logger.Debug("Console debug")
	logger.Info("Console info")
	logger.Warn("Console warn")
	logger.Error("Console error")

	// 测试字段
	logger.Info("Console with fields",
		Field{Key: "console", Value: "test"},
	)

	// 测试Sync
	if err := logger.Sync(); err != nil {
		t.Errorf("Console sync failed: %v", err)
	}
}

// TestLoggerOptions 测试日志选项
func TestLoggerOptions(t *testing.T) {
	// 测试选项函数
	levelOpt := WithLevel(DebugLevel)
	formatOpt := WithFormat("json")
	outputOpt := WithOutputPath("stdout")
	configOpt := WithConfig(map[string]interface{}{"key": "value"})

	options := &LoggerOptions{}

	// 应用选项
	levelOpt(options)
	formatOpt(options)
	outputOpt(options)
	configOpt(options)

	// 验证选项
	if options.Level != DebugLevel {
		t.Errorf("Expected level DebugLevel, got %v", options.Level)
	}

	if options.Format != "json" {
		t.Errorf("Expected format 'json', got '%s'", options.Format)
	}

	if options.OutputPath != "stdout" {
		t.Errorf("Expected output path 'stdout', got '%s'", options.OutputPath)
	}

	if options.Config["key"] != "value" {
		t.Errorf("Expected config key 'key' to be 'value', got '%v'", options.Config["key"])
	}
}

// TestLogLevelString 测试日志级别字符串表示
func TestLogLevelString(t *testing.T) {
	testCases := []struct {
		level    LogLevel
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{FatalLevel, "FATAL"},
		{PanicLevel, "PANIC"},
		{LogLevel(999), "UNKNOWN"},
	}

	for _, tc := range testCases {
		if tc.level.String() != tc.expected {
			t.Errorf("Expected level %v to be '%s', got '%s'", tc.level, tc.expected, tc.level.String())
		}
	}
}

// TestWithMethods 测试With系列方法
func TestWithMethods(t *testing.T) {
	logger := NewConsoleLogger("test-with")

	// 测试WithField
	logger1 := logger.WithField("key1", "value1")
	if logger1 == nil {
		t.Error("Failed to create logger with field")
	}

	// 测试WithFields
	logger2 := logger.WithFields(
		Field{Key: "key1", Value: "value1"},
		Field{Key: "key2", Value: "value2"},
	)
	if logger2 == nil {
		t.Error("Failed to create logger with fields")
	}

	// 测试WithContext
	ctx := context.Background()
	logger3 := logger.WithContext(ctx)
	if logger3 == nil {
		t.Error("Failed to create logger with context")
	}

	// 测试WithError
	err := errors.New("test error")
	logger4 := logger.WithError(err)
	if logger4 == nil {
		t.Error("Failed to create logger with error")
	}

	// 测试WithTime
	now := time.Now()
	logger5 := logger.WithTime(now)
	if logger5 == nil {
		t.Error("Failed to create logger with time")
	}
}

// TestIsEnabledMethods 测试IsEnabled系列方法
func TestIsEnabledMethods(t *testing.T) {
	// 测试DebugLevel
	debugLogger := NewConsoleLogger("test-debug", WithLevel(DebugLevel))
	if !debugLogger.IsDebugEnabled() {
		t.Error("Debug level should be enabled")
	}
	if !debugLogger.IsInfoEnabled() {
		t.Error("Info level should be enabled")
	}
	if !debugLogger.IsWarnEnabled() {
		t.Error("Warn level should be enabled")
	}
	if !debugLogger.IsErrorEnabled() {
		t.Error("Error level should be enabled")
	}

	// 测试InfoLevel
	infoLogger := NewConsoleLogger("test-info", WithLevel(InfoLevel))
	if infoLogger.IsDebugEnabled() {
		t.Error("Debug level should not be enabled")
	}
	if !infoLogger.IsInfoEnabled() {
		t.Error("Info level should be enabled")
	}

	// 测试ErrorLevel
	errorLogger := NewConsoleLogger("test-error", WithLevel(ErrorLevel))
	if errorLogger.IsDebugEnabled() {
		t.Error("Debug level should not be enabled")
	}
	if errorLogger.IsInfoEnabled() {
		t.Error("Info level should not be enabled")
	}
	if errorLogger.IsWarnEnabled() {
		t.Error("Warn level should not be enabled")
	}
	if !errorLogger.IsErrorEnabled() {
		t.Error("Error level should be enabled")
	}
}
