package main

import (
	"app/backend/pkg/logger"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {
	// 解析配置文件
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	} else {
		// 如果没有指定配置文件，则默认使用当前目录下的 config.yaml
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return
	}

	// 环境变量覆盖配置文件中的变量
	viper.AutomaticEnv()

	//---- 默认配置
	viper.SetDefault("env", "development")
	viper.SetDefault("server.port", 21421)
	viper.SetDefault("server.host", "0.0.0.0")

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.file", "")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 30)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.compress", true)
	viper.SetDefault("log.console", true)

	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.user", "postgres")
	viper.SetDefault("db.password", "postgres")
	viper.SetDefault("db.name", "app")
	viper.SetDefault("db.sslmode", "disable")
	viper.SetDefault("db.timezone", "Asia/Shanghai")
	viper.SetDefault("db.pool_max_conns", 10)
	viper.SetDefault("db.pool_max_idle_conns", 5)
	viper.SetDefault("db.pool_max_lifetime", 10*time.Minute)
	viper.SetDefault("db.pool_max_idle_time", 5*time.Minute)

	viper.SetDefault("frontend.port", 21422)
	viper.SetDefault("frontend.host", "0.0.0.0")
	viper.SetDefault("frontend.dist", "./dist")
	viper.SetDefault("frontend.isSpa", true) // 是否是 spa 应用，也就是前端路由是单页应用

	//----

	// 更新并存储最新的配置，犹豫存在后续代码更新但是目前的配置文件不是最新的情况
	if err := viper.WriteConfig(); err != nil {
		panic(err)
	}

	// 初始化日志
	logConfig := &logger.Config{
		Level:      viper.GetString("log.level"),
		File:       viper.GetString("log.file"),
		MaxSize:    viper.GetInt("log.max_size"),
		MaxBackups: viper.GetInt("log.max_backups"),
		MaxAge:     viper.GetInt("log.max_age"),
		Compress:   viper.GetBool("log.compress"),
	}
	if err := logger.InitLogger(logConfig); err != nil {
		panic(err)
	}

	logger.Info("Current config: ", zap.Any("config", viper.AllSettings()))
}
