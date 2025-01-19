package main

import (
	"app/backend/core"
	"app/backend/service"

	"github.com/spf13/viper"
)

func InitApp() *core.App {
	coreApp := core.NewCoreApp(viper.GetViper())

	testService := service.NewTestService(coreApp)
	coreApp.RegisterService(testService)

	return coreApp
}
