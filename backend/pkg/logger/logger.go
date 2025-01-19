package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level string `yaml:"level"` // 日志级别
	File  string `yaml:"file"`  // 日志文件路径
	// 日志滚动配置
	MaxSize    int  `yaml:"max_size"`    // 每个日志文件的最大大小（MB）
	MaxBackups int  `yaml:"max_backups"` // 保留的旧日志文件的最大数量
	MaxAge     int  `yaml:"max_age"`     // 保留旧日志文件的最大天数
	Compress   bool `yaml:"compress"`    // 是否压缩旧日志文件
}

var globalLogger *zap.Logger

// getConsoleEncoderConfig 获取控制台编码器配置
func getConsoleEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 使用带颜色的日志级别
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339Nano),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// getJSONEncoderConfig 获取JSON编码器配置
func getJSONEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// InitLogger 初始化日志
func InitLogger(cfg *Config) error {
	// 设置默认值
	if cfg.Level == "" {
		cfg.Level = "info" // 默认使用 info 级别
	}
	if cfg.MaxSize == 0 {
		cfg.MaxSize = 100 // 默认100MB
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 30 // 默认保留30个备份
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 30 // 默认保留30天
	}

	// 解析日志级别
	level := zap.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return fmt.Errorf("parse log level error: %w", err)
	}

	var cores []zapcore.Core

	// 添加控制台输出（带颜色）
	cores = append(cores, zapcore.NewCore(
		zapcore.NewConsoleEncoder(getConsoleEncoderConfig()),
		zapcore.Lock(os.Stdout),
		level,
	))

	// 如果配置了文件输出，添加文件输出
	if cfg.File != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(cfg.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("create log directory error: %w", err)
		}

		// 配置 lumberjack 进行日志滚动
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.File,
			MaxSize:    cfg.MaxSize,    // 每个日志文件的最大大小（MB）
			MaxBackups: cfg.MaxBackups, // 保留的旧日志文件的最大数量
			MaxAge:     cfg.MaxAge,     // 保留旧日志文件的最大天数
			Compress:   cfg.Compress,   // 是否压缩旧日志文件
			LocalTime:  true,           // 使用本地时间
		})

		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(getJSONEncoderConfig()),
			fileWriter,
			level,
		))
	}

	// 创建logger
	core := zapcore.NewTee(cores...)
	globalLogger = zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel), // Error 及以上级别日志添加堆栈信息
	)

	return nil
}

// GetLogger 获取全局logger
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		globalLogger = zap.NewExample()
	}
	return globalLogger
}

// Debug logs a message at DebugLevel
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info logs a message at InfoLevel
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn logs a message at WarnLevel
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error logs a message at ErrorLevel
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal logs a message at FatalLevel
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}
