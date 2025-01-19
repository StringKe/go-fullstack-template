package util

import "github.com/spf13/viper"

// 判断是否是开发环境
func IsDev() bool {
	return viper.GetString("env") == "development" || viper.GetString("env") == "dev"
}

// 辅助函数
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
