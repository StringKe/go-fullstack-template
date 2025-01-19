package core

import (
	"github.com/spf13/viper"
)

type App struct {
	cfg       *viper.Viper
	container *Container
}

func NewCoreApp(cfg *viper.Viper) *App {
	serviceContainer := NewContainer()
	return &App{
		cfg:       cfg,
		container: serviceContainer,
	}
}

func (c *App) GetContainer() *Container {
	return c.container
}

func (c *App) RegisterService(svc ServiceHandler) {
	c.container.RegisterService(svc)
}

func (c *App) GetConfig() *viper.Viper {
	return c.cfg
}
