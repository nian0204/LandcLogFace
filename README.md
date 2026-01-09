# LandcLogFace - Go语言日志门面项目

LandcLogFace是一个Go语言的日志门面（Logging Facade）项目，提供统一的日志接口，支持多种常用日志库的适配，帮助开发者在不同日志实现之间无缝切换，同时保持代码的一致性和可维护性。

## 项目特性

- **统一的日志接口**：定义标准的Logger接口，包含各种日志级别和功能
- **多日志库支持**：实现多种常用日志库的适配器
  - 控制台日志（默认）
  - zap日志库（高性能）
  - logrus日志库（功能丰富）
  - 标准库log（轻量）
- **灵活的配置管理**：支持通过选项函数和配置map进行灵活配置
- **日志工厂**：提供统一的日志实例创建和管理功能
- **全局日志**：提供便捷的全局日志函数
- **结构化字段**：支持添加结构化日志字段
- **上下文支持**：支持添加上下文信息
- **错误处理**：支持添加错误信息
- **时间管理**：支持自定义时间字段
- **日志级别控制**：支持细粒度的日志级别控制
- **可扩展性**：支持自定义日志提供者

## 安装

### 前提条件

- Go 1.24.0或更高版本

### 安装步骤

1. **使用go get安装**

```bash
go get github.com/nian0204/LandcLogFace
```

2. **或直接克隆仓库**

```bash
git clone https://github.com/nian0204/LandcLogFace.git
cd LandcLogFace
go mod tidy
```

## 快速开始

### 基本使用

```go
package main

import (
	"LandcLogFace"
)

func main() {
	// 获取全局日志实例
	logger := LandcLogFace.GetLogger()

	// 输出不同级别的日志
	logger.Debug("这是一条调试日志")
	logger.Info("这是一条信息日志")
	logger.Warn("这是一条警告日志")
	logger.Error("这是一条错误日志")

	// 输出格式化日志
	logger.Infof("Hello, %s!", "world")
	logger.Errorf("Error: %d", 404)

	// 使用字段
	logger.Info("用户登录", 
		LandcLogFace.Field{Key: "user_id", Value: 123},
		LandcLogFace.Field{Key: "ip", Value: "192.168.1.1"},
	)
}
```

### 使用全局函数

```go
package main

import (
	"LandcLogFace"
)

func main() {
	// 使用全局函数输出日志
	LandcLogFace.Debug("全局调试日志")
	LandcLogFace.Info("全局信息日志")
	LandcLogFace.Warn("全局警告日志")
	LandcLogFace.Error("全局错误日志")

	// 使用全局格式化函数
	LandcLogFace.Infof("全局格式化日志: %s", "test")
	LandcLogFace.Errorf("全局错误格式化日志: %d", 500)
}
```

## 详细使用指南

### 1. 选择不同的日志库

LandcLogFace支持多种日志库，你可以根据需要选择合适的日志实现：

```go
package main

import (
	"LandcLogFace"
)

func main() {
	// 使用console提供者（默认）
	consoleLogger := LandcLogFace.GetLogFactory().CreateLoggerWithProvider("app", "console")
	consoleLogger.Info("使用控制台日志")

	// 使用zap提供者（高性能）
	zapLogger := LandcLogFace.GetLogFactory().CreateLoggerWithProvider("app", "zap")
	zapLogger.Info("使用zap日志")

	// 使用logrus提供者（功能丰富）
	logrusLogger := LandcLogFace.GetLogFactory().CreateLoggerWithProvider("app", "logrus")
	logrusLogger.Info("使用logrus日志")

	// 使用std提供者（轻量）
	stdLogger := LandcLogFace.GetLogFactory().CreateLoggerWithProvider("app", "std")
	stdLogger.Info("使用标准库日志")
}
```

### 2. 配置日志实例

你可以通过选项函数或配置map来配置日志实例：

#### 使用选项函数

```go
package main

import (
	"LandcLogFace"
)

func main() {
	// 使用选项函数配置
	logger := LandcLogFace.GetLogFactory().CreateLoggerWithProvider("app", "zap",
		LandcLogFace.WithLevel(LandcLogFace.DebugLevel),
		LandcLogFace.WithFormat("json"),
		LandcLogFace.WithOutputPath("stdout"),
	)

	logger.Debug("调试日志")
}
```

#### 使用配置map

```go
package main

import (
	"LandcLogFace"
)

func main() {
	// 使用配置map
	config := map[string]interface{}{
		"provider":    "zap",
		"level":       LandcLogFace.DebugLevel,
		"format":      "json",
		"outputPath":  "stdout",
	}

	logger := LandcLogFace.GetLogFactory().CreateLoggerWithConfig("app", config)
	logger.Info("带配置的日志")
}
```

### 3. 高级功能

#### 字段使用

```go
package main

import (
	"LandcLogFace"
)

func main() {
	logger := LandcLogFace.GetLogger()

	// 创建带字段的日志实例
	userLogger := logger.WithFields(
		LandcLogFace.Field{Key: "service", Value: "user"},
		LandcLogFace.Field{Key: "version", Value: "1.0.0"},
	)

	// 使用带字段的日志实例
	userLogger.Info("用户注册", 
		LandcLogFace.Field{Key: "username", Value: "john"},
		LandcLogFace.Field{Key: "email", Value: "john@example.com"},
	)

	// 链式添加字段
	userLogger.WithField("status", "success").Info("注册成功")
}
```

#### 上下文支持

```go
package main

import (
	"LandcLogFace"
	"context"
)

func main() {
	logger := LandcLogFace.GetLogger()

	// 创建上下文
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "123456")

	// 添加上下文到日志
	ctxLogger := logger.WithContext(ctx)
	ctxLogger.Info("带上下文的日志")
}
```

#### 错误处理

```go
package main

import (
	"LandcLogFace"
	"errors"
)

func main() {
	logger := LandcLogFace.GetLogger()

	// 创建错误
	err := errors.New("数据库连接失败")

	// 添加错误到日志
	errLogger := logger.WithError(err)
	errLogger.Error("操作失败")

	// 直接在日志中添加错误字段
	logger.Error("操作失败", LandcLogFace.Field{Key: "error", Value: err})
}
```

#### 时间管理

```go
package main

import (
	"LandcLogFace"
	"time"
)

func main() {
	logger := LandcLogFace.GetLogger()

	// 添加自定义时间
	now := time.Now()
	logger.WithTime(now).Info("带自定义时间的日志")
}
```

#### 日志级别检查

```go
package main

import (
	"LandcLogFace"
)

func main() {
	logger := LandcLogFace.GetLogger()

	// 检查日志级别是否启用
	if logger.IsDebugEnabled() {
		// 执行耗时操作
		data := "耗时操作的结果"
		logger.Debug("调试信息", LandcLogFace.Field{Key: "data", Value: data})
	}

	if logger.IsInfoEnabled() {
		logger.Info("信息日志已启用")
	}
}
```

### 4. 自定义日志提供者

如果你需要使用项目未内置的日志库，可以通过实现`LoggerProvider`接口来添加自定义日志提供者：

```go
package main

import (
	"LandcLogFace"
	"fmt"
)

// CustomLogger 自定义日志实现
type CustomLogger struct {
	name string
}

// NewCustomLogger 创建自定义日志实例
func NewCustomLogger(name string) *CustomLogger {
	return &CustomLogger{name: name}
}

// 实现Logger接口的所有方法...
// （具体实现可参考项目中的示例代码）

// CustomLoggerProvider 自定义日志提供者
type CustomLoggerProvider struct{}

// Create 创建日志实例
func (p *CustomLoggerProvider) Create(name string) LandcLogFace.Logger {
	return NewCustomLogger(name)
}

// CreateWithConfig 根据配置创建日志实例
func (p *CustomLoggerProvider) CreateWithConfig(name string, config map[string]interface{}) LandcLogFace.Logger {
	return NewCustomLogger(name)
}

func main() {
	// 注册自定义提供者
	LandcLogFace.GetLogFactory().RegisterProvider("custom", &CustomLoggerProvider{})

	// 使用自定义提供者
	customLogger := LandcLogFace.GetLogFactory().CreateLoggerWithProvider("app", "custom")
	customLogger.Info("使用自定义日志提供者")

	// 注销自定义提供者
	LandcLogFace.GetLogFactory().UnregisterProvider("custom")
}
```

## 项目结构

```
LandcLogFace/
├── go.mod                # 项目依赖管理
├── logger.go             # 核心接口定义
├── console_logger.go     # 控制台日志适配器
├── zap_logger.go         # zap日志库适配器
├── logrus_logger.go      # logrus日志库适配器
├── std_logger.go         # 标准库log适配器
├── log_factory.go        # 日志工厂和配置管理
├── example.go            # 使用示例
├── logger_test.go        # 测试用例
└── README.md             # 项目文档
```

## 依赖管理

项目使用Go 1.24.0版本，依赖以下第三方库：

| 依赖库 | 版本 | 用途 |
|-------|------|------|
| `go.uber.org/zap` | v1.26.0 | 高性能日志库 |
| `github.com/sirupsen/logrus` | v1.9.3 | 功能丰富的日志库 |

## 测试

项目包含完整的测试用例，验证了所有核心功能的正确性：

```bash
# 运行所有测试
cd LandcLogFace
go test -v ./...
```

测试覆盖以下内容：

- 核心接口功能测试
- 日志工厂测试
- 全局日志测试
- 各种日志适配器测试
- 日志选项测试
- 日志级别测试
- 高级功能测试（字段、上下文、错误处理等）

所有测试用例都已通过，确保项目的可靠性和稳定性。

## 贡献指南

欢迎为LandcLogFace项目贡献代码！如果你有任何改进或新功能的想法，请按照以下步骤进行：

1. **Fork** 项目仓库
2. **Clone** 到本地：`git clone https://github.com/nian0204/LandcLogFace.git`
3. **创建** 特性分支：`git checkout -b feature/your-feature`
4. **实现** 你的功能或修复
5. **编写** 测试用例
6. **运行** 测试：`go test -v ./...`
7. **提交** 代码：`git commit -m "Add your feature"`
8. **推送** 到远程：`git push origin feature/your-feature`
9. **创建** Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 联系方式

如果您有任何问题或建议，请通过以下方式联系我们：

- GitHub Issues：[https://github.com/nian0204/LandcLogFace/issues](https://github.com/nian0204/LandcLogFace/issues)
- 邮箱：your.email@example.com

---

**LandcLogFace** - 让Go语言日志管理更简单、更灵活、更强大！