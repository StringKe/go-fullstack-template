package serve

import (
	"app/backend/core"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/http2"
)

type App struct {
	core  *core.App
	serve *echo.Echo
}

func NewServeApp(core *core.App) (*App, error) {
	serve := echo.New()
	// gRPC requires HTTP/2
	serve.HideBanner = true
	serve.Use(middleware.Logger())
	serve.Use(middleware.Recover())
	// Enable CORS
	serve.Use(middleware.CORS())

	return &App{
		core:  core,
		serve: serve,
	}, nil
}

func (app *App) Start() error {
	// 获取 Vanguard transcoder 也就是 http 处理器
	transcoder, err := app.core.GetContainer().GetVanguardTranscoder()
	if err != nil {
		return fmt.Errorf("failed to get vanguard transcoder: %w", err)
	}

	stripPrefixHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/rpc")
		transcoder.ServeHTTP(w, r)
	})

	// 注册路由处理器
	app.serve.Any("/rpc/*", echo.WrapHandler(stripPrefixHandler))

	port := app.core.GetConfig().GetInt("server.port")
	addr := fmt.Sprintf(":%d", port)
	// 使用 H2C (HTTP/2 Cleartext) 模式启动服务器
	h2s := &http2.Server{
		MaxConcurrentStreams: 250,
	}
	return app.serve.StartH2CServer(addr, h2s)
}
