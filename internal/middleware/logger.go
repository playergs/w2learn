package middleware

import (
	"time"
	"w2learn/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 前置过滤取请求时间、URL地址、请求参数
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		// 后置过滤使用 封装zap的logger进行日志打印
		logger.Debug("HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", time.Since(start)),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}
