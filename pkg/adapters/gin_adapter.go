package adapters

import (
	"fmt"
	"time"

	"github.com/LandcLi/LandcLogFace/pkg/logger"

	"github.com/gin-gonic/gin"
)

// GinLogger 是gin框架的日志适配器
type GinLogger struct {
	log Logger
}

// NewGinLogger 创建一个新的gin日志适配器
func NewGinLogger(log Logger) *GinLogger {
	return &GinLogger{
		log: log,
	}
}

// Logger 返回gin的日志中间件
func (g *GinLogger) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()
		traceID := c.Request.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = fmt.Sprintf("%d", time.Now().UnixNano())
		}

		// 日志字段
		fields := []logger.Field{
			{Key: "status", Value: statusCode},
			{Key: "method", Value: reqMethod},
			{Key: "uri", Value: reqUri},
			{Key: "ip", Value: clientIP},
			{Key: "latency", Value: latencyTime},
			{Key: "timestamp", Value: endTime},
			{Key: "trace_id", Value: traceID},
		}

		// 根据状态码设置日志级别
		switch {
		case statusCode >= 500:
			g.log.Error(fmt.Sprintf("[GIN] %s %s %d %s %s", reqMethod, reqUri, statusCode, latencyTime, clientIP), fields...)
		case statusCode >= 400:
			g.log.Warn(fmt.Sprintf("[GIN] %s %s %d %s %s", reqMethod, reqUri, statusCode, latencyTime, clientIP), fields...)
		case statusCode >= 300:
			g.log.Info(fmt.Sprintf("[GIN] %s %s %d %s %s", reqMethod, reqUri, statusCode, latencyTime, clientIP), fields...)
		default:
			g.log.Info(fmt.Sprintf("[GIN] %s %s %d %s %s", reqMethod, reqUri, statusCode, latencyTime, clientIP), fields...)
		}
	}
}

// Recovery 返回gin的恢复中间件，使用我们的日志门面记录错误
func (g *GinLogger) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				g.log.Error(fmt.Sprintf("[GIN] panic recovered: %v", err),
					logger.Field{Key: "method", Value: c.Request.Method},
					logger.Field{Key: "uri", Value: c.Request.RequestURI},
					logger.Field{Key: "ip", Value: c.ClientIP()},
				)

				// 响应500错误
				c.AbortWithStatus(500)
			}
		}()

		c.Next()
	}
}

// UseWithGin 将日志适配器应用到gin引擎
func UseWithGin(r *gin.Engine, log interface{}) {
	// 类型断言，确保log实现了必要的方法
	if logger, ok := log.(Logger); ok {
		ginLogger := NewGinLogger(logger)
		r.Use(ginLogger.Logger())
		r.Use(ginLogger.Recovery())
	}
}
