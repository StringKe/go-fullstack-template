package main

import (
	"app/backend/core"
	"app/backend/service"

	"github.com/spf13/viper"
)

func InitApp() (*core.App, error) {
	config := viper.GetViper()
	core := core.NewCoreApp(config)

	testService := service.NewTestService(core)
	core.RegisterService(testService)

	return core, nil
}
