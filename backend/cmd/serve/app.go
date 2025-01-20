package serve

import (
	"app/backend/core"
	"app/backend/pkg/logger"
	"app/backend/pkg/util"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
)

var (
	SkipRoutePath = []string{}
)

type App struct {
	core  *core.App
	serve *echo.Echo
}

func NewServeApp(core *core.App) (*App, error) {
	serve := echo.New()
	serve.HideBanner = true
	serve.HidePort = true

	serve.Use(middleware.RequestID())
	serve.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:           true,
		LogMethod:        true,
		LogStatus:        true,
		LogLatency:       true,
		LogRemoteIP:      true,
		LogUserAgent:     true,
		LogRequestID:     true,
		LogReferer:       true,
		LogHost:          true,
		LogProtocol:      true,
		LogRoutePath:     true,
		LogURIPath:       true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogError:         true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			path := fmt.Sprintf("%s %s %s", v.RequestID, v.Method, v.URI)
			if v.Error == nil {
				logger.Info(path,
					zap.Int("status", v.Status),
					zap.String("latency", v.Latency.String()),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("user_agent", v.UserAgent),
					zap.String("referer", v.Referer),
					zap.String("host", v.Host),
					zap.String("protocol", v.Protocol),
					zap.String("route_path", v.RoutePath),
					zap.String("uri_path", v.URIPath),
					zap.String("content_length", v.ContentLength),
					zap.Int64("response_size", v.ResponseSize),
				)
			} else {
				logger.Error(path,
					zap.Int("status", v.Status),
					zap.String("latency", v.Latency.String()),
					zap.String("remote_ip", v.RemoteIP),
					zap.String("user_agent", v.UserAgent),
					zap.String("referer", v.Referer),
					zap.String("host", v.Host),
					zap.String("protocol", v.Protocol),
					zap.String("route_path", v.RoutePath),
					zap.String("uri_path", v.URIPath),
					zap.String("content_length", v.ContentLength),
					zap.Int64("response_size", v.ResponseSize),
					zap.Error(v.Error),
				)
			}
			return nil

		},
	}))

	serve.Use(middleware.Recover())
	serve.Use(middleware.CORS())
	serve.Use(middleware.Decompress())
	serve.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	return &App{
		core:  core,
		serve: serve,
	}, nil
}

func (app *App) SetupRpc(ctx context.Context) error {
	transcoder, err := app.core.GetContainer().BuildTranscoder()
	if err != nil {
		return err
	}

	app.serve.Any("/rpc/*", echo.WrapHandler(transcoderPrefixHandler("/rpc", transcoder)))
	SkipRoutePath = append(SkipRoutePath, "/rpc/*")
	return nil
}

func (app *App) SetupFrontendDev(ctx context.Context) error {
	config := app.core.GetConfig()
	frontendPort := config.GetInt("frontend.port")
	target := fmt.Sprintf("http://localhost:%d", frontendPort)

	app.serve.Use(proxyFrontend(target))
	return nil
}

func (app *App) SetupFrontend(ctx context.Context) error {
	config := app.core.GetConfig()
	distPath := config.GetString("frontend.dist")
	isSpa := config.GetBool("frontend.isSpa")

	// 生产环境：服务静态文件
	staticConfig := middleware.StaticConfig{
		Root:       distPath,
		Index:      "index.html",
		HTML5:      isSpa,
		Browse:     false,
		IgnoreBase: true,
	}
	app.serve.Use(middleware.StaticWithConfig(staticConfig))

	if isSpa {
		// SPA 模式：所有未匹配的路由都返回 index.html
		app.serve.GET("*", func(c echo.Context) error {
			// 跳过已经由静态文件中间件处理的请求
			for _, path := range SkipRoutePath {
				if strings.HasPrefix(c.Path(), path) {
					return nil
				}
			}
			return c.File(filepath.Join(distPath, "index.html"))
		})
	}
	return nil
}

func (app *App) Listen(ctx context.Context) error {
	port := app.core.GetConfig().GetInt("server.port")
	addr := fmt.Sprintf(":%d", port)

	h2s := &http2.Server{
		MaxConcurrentStreams: 250,
	}

	// 启动服务器
	go func() {
		logger.Info(fmt.Sprintf("Local: http://localhost:%d/", port))
		addrs, err := util.GetAllIPAddresses()
		if err == nil {
			for _, addr := range addrs {
				logger.Info(fmt.Sprintf("Network: http://%s:%d/", addr, port))
			}
		}

		if err := app.serve.StartH2CServer(addr, h2s); err != nil && err != http.ErrServerClosed {
			logger.Error("Server start failed", zap.Error(err))
		}
	}()

	// 等待上下文取消
	<-ctx.Done()
	logger.Info("Current context is done, shutting down server...")

	// 优雅停机
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown failed", zap.Error(err))
		return err
	}

	logger.Info("Server has been successfully shut down")
	return nil
}

func (app *App) Start(ctx context.Context) error {
	isDev := util.IsDev()

	app.SetupRpc(ctx)

	if isDev {
		app.SetupFrontendDev(ctx)
	} else {
		app.SetupFrontend(ctx)
	}

	return app.Listen(ctx)
}

func (app *App) Shutdown(ctx context.Context) error {
	return app.serve.Shutdown(ctx)
}
