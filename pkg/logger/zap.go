package logger

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapConfig 日志配置结构体（可从 YAML/JSON 配置文件加载）
type ZapConfig struct {
	Env        string `mapstructure:"env"`         // 环境：dev/prod
	Level      string `mapstructure:"level"`       // 日志级别：debug/info/warn/error
	Filepath   string `mapstructure:"filepath"`    // 日志文件路径（prod 环境必填）
	MaxSize    int    `mapstructure:"max_size"`    // 单个日志文件最大大小（MB）
	MaxBackups int    `mapstructure:"max_backups"` // 最大备份文件数
	MaxAge     int    `mapstructure:"max_age"`     // 日志保留天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩备份文件
	ShowCaller bool   `mapstructure:"show_caller"` // 是否显示调用者（文件:行号）
	Stacktrace string `mapstructure:"stacktrace"`  // 开启堆栈跟踪的级别（如 error）
}

// NewZapLogger 初始化 Zap Logger（核心函数）
func NewZapLogger(cfg ZapConfig) (*zap.Logger, error) {
	// 1. 解析日志级别
	infoLevel, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	//// 2. 配置 Encoder（日志格式）
	//encoder := getEncoderByEnv(cfg.Env)

	// 2. 配置控制台输出 (ConsoleEncoder)
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 彩色级别输出
		EncodeTime:     customTimeEncoder,                // 自定义时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径调用者 (e.g., pkg/file.go:23)
	})
	// 3. 配置文件输出 (JSONEncoder)
	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder, // 文件输出不需要颜色
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 文件中记录完整路径
	})
	// 4. 配置日志文件写入器 (带轮转功能)
	fileWriter := &lumberjack.Logger{
		Filename:   "./logs/app.log", // 日志文件路径
		MaxSize:    10,               // 每个日志文件最大 10 MB
		MaxBackups: 5,                // 最多保留 5 个备份文件
		MaxAge:     30,               // 最多保留 30 天
		Compress:   true,             // 是否压缩备份文件
	}
	// 控制台输出
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), infoLevel)
	// 文件输出
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), infoLevel)
	// 6. 将多个 Core 组合
	core := zapcore.NewTee(consoleCore, fileCore)
	//// 3. 配置 WriteSyncer（日志输出目标）
	//writeSyncers, err := getWriteSyncersByEnv(cfg)
	//if err != nil {
	//	return nil, err
	//}
	// 控制台输出
	//consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), infoLevel)
	// 4. 配置 Core（Zap 核心：Encoder + WriteSyncer + Level）
	//core := zapcore.NewTee(consoleCore)

	// 5. 配置 Logger 选项（调用者、堆栈跟踪等）
	opts := getLoggerOptions(cfg)
	// 6. 创建并返回 Logger
	logger := zap.New(core, opts...)
	return logger, nil
}

// 1. 解析日志级别（string → zapcore.Level）
func parseLevel(levelStr string) (zapcore.Level, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		return level, err
	}
	return level, nil
}

// 2. 根据环境获取 Encoder（文本/JSON）
func getEncoderByEnv(env string) zapcore.Encoder {
	// 基础配置：时间格式、字段顺序
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customTimeEncoder, // 自定义时间格式（统一为 "2006-01-02 15:04:05.000"）
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 调用者格式：pkg/file.go:123
	}

	// 开发环境：彩色文本格式
	if env == "dev" {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // 彩色级别（ERROR → 红色）
		return zapcore.NewConsoleEncoder(encoderCfg)
	}

	// 生产环境：JSON 格式（便于日志系统解析）
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder // 无颜色，纯文本级别（ERROR/WARN）
	return zapcore.NewJSONEncoder(encoderCfg)
}

// 3. 根据环境获取 WriteSyncer（控制台/文件）
func getWriteSyncersByEnv(cfg ZapConfig) ([]zapcore.WriteSyncer, error) {
	var writeSyncers []zapcore.WriteSyncer

	// 开发环境：仅输出到控制台
	if cfg.Env == "dev" {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
		return writeSyncers, nil
	}

	// 生产环境：输出到文件（带轮转）+ 可选控制台（调试时）
	if cfg.Filepath == "" {
		//return nil, zap.Error(zapcore.ErrInvalidWriteSyncer)
	}

	// 配置日志轮转（lumberjack）
	lumberjackLogger := &lumberjack.Logger{
		Filename:   cfg.Filepath,   // 日志文件路径（如 ./logs/app.log）
		MaxSize:    cfg.MaxSize,    // 单个文件最大 MB 数
		MaxBackups: cfg.MaxBackups, // 最大备份数（超过自动删除）
		MaxAge:     cfg.MaxAge,     // 保留天数（超过自动删除）
		Compress:   cfg.Compress,   // 压缩备份文件（节省磁盘）
	}
	writeSyncers = append(writeSyncers, zapcore.AddSync(lumberjackLogger))

	// （可选）生产环境也输出到控制台（便于紧急调试，建议线上关闭）
	// writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))

	return writeSyncers, nil
}

// 4. 获取 Logger 选项（调用者、堆栈跟踪等）
func getLoggerOptions(cfg ZapConfig) []zap.Option {
	var opts []zap.Option

	// 显示调用者（文件:行号）
	if cfg.ShowCaller {
		opts = append(opts, zap.AddCaller())
	}

	// 开启堆栈跟踪（仅指定级别及以上）
	if cfg.Stacktrace != "" {
		stackLevel, err := parseLevel(cfg.Stacktrace)
		if err == nil {
			opts = append(opts, zap.AddStacktrace(stackLevel))
		}
	}

	// （可选）添加日志实例名称（多模块区分，如 "user-service"）
	// opts = append(opts, zap.AddLoggerName("app"))

	return opts
}

// customTimeEncoder 自定义时间格式（统一格式，便于跨环境解析）
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// 自定义级别编码器（用于彩色输出）
func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}
