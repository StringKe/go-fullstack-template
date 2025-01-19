package main

import (
	"context"

	"app/backend/cmd/config"
	"app/backend/cmd/serve"
	"app/backend/core"

	"github.com/spf13/cobra"
)

var (
	CfgFile string
	RootCmd = &cobra.Command{
		Use:   "app",
		Short: "a full stack web framework",
		Long:  `Your app description`,
	}
)

func init() {
	// 提前解析配置文件路径参数
	RootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is ./config.yaml)")

	// 读取并初始化配置
	InitConfig()
	// 初始化核心
	coreApp := InitApp()
	// 将核心应用添加到根命令的上下文中
	RootCmd.SetContext(context.WithValue(context.Background(), core.CoreAppKey, coreApp))

	// 添加子命令，如果子命令还存在子命令，需要在子命令的 init 函数中添加
	RootCmd.AddCommand(config.Command)
	RootCmd.AddCommand(serve.Command)
}
