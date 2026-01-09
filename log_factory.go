package LandcLogFace

import (
	"sync"
)

// LogFactory 日志工厂
type LogFactory struct {
	providers map[string]LoggerProvider
	defaultProvider string
	mu sync.RWMutex
}

// 全局日志工厂实例
var (
	factory *LogFactory
	factoryOnce sync.Once
)

// GetLogFactory 获取全局日志工厂实例
func GetLogFactory() *LogFactory {
	factoryOnce.Do(func() {
		factory = NewLogFactory()
		// 注册默认的日志提供者
		factory.RegisterProvider("console", NewConsoleLoggerProvider())
		factory.RegisterProvider("zap", NewZapLoggerProvider())
		factory.RegisterProvider("logrus", NewLogrusLoggerProvider())
		factory.RegisterProvider("std", NewStdLoggerProvider())
		// 设置默认提供者为console
		factory.SetDefaultProvider("console")
	})
	return factory
}

// NewLogFactory 创建日志工厂实例
func NewLogFactory() *LogFactory {
	return &LogFactory{
		providers: make(map[string]LoggerProvider),
		defaultProvider: "console",
	}
}

// RegisterProvider 注册日志提供者
func (f *LogFactory) RegisterProvider(name string, provider LoggerProvider) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.providers[name] = provider
}

// UnregisterProvider 注销日志提供者
func (f *LogFactory) UnregisterProvider(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.providers, name)
}

// SetDefaultProvider 设置默认日志提供者
func (f *LogFactory) SetDefaultProvider(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.defaultProvider = name
}

// GetDefaultProvider 获取默认日志提供者名称
func (f *LogFactory) GetDefaultProvider() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.defaultProvider
}

// GetProvider 获取指定的日志提供者
func (f *LogFactory) GetProvider(name string) (LoggerProvider, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	provider, exists := f.providers[name]
	return provider, exists
}

// CreateLogger 创建日志实例
func (f *LogFactory) CreateLogger(name string) Logger {
	return f.CreateLoggerWithProvider(name, f.defaultProvider)
}

// CreateLoggerWithProvider 使用指定的提供者创建日志实例
func (f *LogFactory) CreateLoggerWithProvider(name string, providerName string) Logger {
	f.mu.RLock()
	provider, exists := f.providers[providerName]
	f.mu.RUnlock()

	if !exists {
		// 如果指定的提供者不存在，使用默认提供者
		f.mu.RLock()
		provider, exists = f.providers[f.defaultProvider]
		f.mu.RUnlock()
		if !exists {
			// 如果默认提供者也不存在，使用控制台日志
			return NewConsoleLogger(name)
		}
	}

	return provider.Create(name)
}

// CreateLoggerWithConfig 根据配置创建日志实例
func (f *LogFactory) CreateLoggerWithConfig(name string, config map[string]interface{}) Logger {
	// 从配置中获取提供者名称
	providerName := f.defaultProvider
	if pn, ok := config["provider"].(string); ok {
		providerName = pn
	}

	f.mu.RLock()
	provider, exists := f.providers[providerName]
	f.mu.RUnlock()

	if !exists {
		// 如果指定的提供者不存在，使用默认提供者
		f.mu.RLock()
		provider, exists = f.providers[f.defaultProvider]
		f.mu.RUnlock()
		if !exists {
			// 如果默认提供者也不存在，使用控制台日志
			return NewConsoleLogger(name)
		}
	}

	return provider.CreateWithConfig(name, config)
}

// 全局日志实例
var (
	globalLogger Logger
	loggerOnce sync.Once
)

// GetLogger 获取全局日志实例
func GetLogger() Logger {
	loggerOnce.Do(func() {
		globalLogger = GetLogFactory().CreateLogger("global")
	})
	return globalLogger
}

// GetLoggerWithName 获取指定名称的日志实例
func GetLoggerWithName(name string) Logger {
	return GetLogFactory().CreateLogger(name)
}

// GetLoggerWithProvider 获取指定提供者的日志实例
func GetLoggerWithProvider(name string, provider string) Logger {
	return GetLogFactory().CreateLoggerWithProvider(name, provider)
}

// GetLoggerWithConfig 根据配置获取日志实例
func GetLoggerWithConfig(name string, config map[string]interface{}) Logger {
	return GetLogFactory().CreateLoggerWithConfig(name, config)
}

// SetGlobalLogger 设置全局日志实例
func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}

// Debug 全局调试级日志
func Debug(msg string, fields ...Field) {
	GetLogger().Debug(msg, fields...)
}

// Debugf 全局格式化调试级日志
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Info 全局信息级日志
func Info(msg string, fields ...Field) {
	GetLogger().Info(msg, fields...)
}

// Infof 全局格式化信息级日志
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warn 全局警告级日志
func Warn(msg string, fields ...Field) {
	GetLogger().Warn(msg, fields...)
}

// Warnf 全局格式化警告级日志
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Error 全局错误级日志
func Error(msg string, fields ...Field) {
	GetLogger().Error(msg, fields...)
}

// Errorf 全局格式化错误级日志
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatal 全局致命级日志
func Fatal(msg string, fields ...Field) {
	GetLogger().Fatal(msg, fields...)
}

// Fatalf 全局格式化致命级日志
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// Panic 全局恐慌级日志
func Panic(msg string, fields ...Field) {
	GetLogger().Panic(msg, fields...)
}

// Panicf 全局格式化恐慌级日志
func Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(format, args...)
}
