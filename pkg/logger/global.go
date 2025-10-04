package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

// InitGlobalLogger 初始化全局 Logger（项目启动时调用一次）
func InitGlobalLogger(cfg ZapConfig) error {
	var err error
	once.Do(func() {
		globalLogger, err = NewZapLogger(cfg)
	})
	return err
}

// GetLogger 获取全局 Logger
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		// 兜底：若未初始化，返回默认 Logger（避免 nil  panic）
		defaultLogger, _ := zap.NewProduction()
		return defaultLogger
	}
	return globalLogger
}

// Sugar （可选）封装 SugaredLogger，简化使用
func Sugar() *zap.SugaredLogger {
	return GetLogger().Sugar()
}
