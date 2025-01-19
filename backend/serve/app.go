package serve

import (
	"app/backend/core"
	"app/backend/pkg/logger"
	"app/backend/pkg/util"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
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

func (app *App) Start(ctx context.Context) error {
	// 获取 Vanguard transcoder 也就是 http 处理器
	transcoder, err := app.core.GetContainer().GetVanguardTranscoder()
	if err != nil {
		return fmt.Errorf("failed to get vanguard transcoder: %w", err)
	}

	stripPrefixHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/rpc")
		transcoder.ServeHTTP(w, r)
	})

	// 先注册 RPC 路由处理器
	app.serve.Any("/rpc/*", echo.WrapHandler(stripPrefixHandler))

	// 再设置前端处理
	app.frontend()

	port := app.core.GetConfig().GetInt("server.port")
	addr := fmt.Sprintf(":%d", port)
	// 使用 H2C (HTTP/2 Cleartext) 模式启动服务器
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

// Shutdown 优雅地关闭服务器
func (app *App) Shutdown(ctx context.Context) error {
	return app.serve.Shutdown(ctx)
}

// 进行前端处理
func (app *App) frontend() {
	config := app.core.GetConfig()
	isDev := util.IsDev()
	frontendPort := config.GetInt("frontend.port")
	distPath := config.GetString("frontend.dist")
	isSpa := config.GetBool("frontend.isSpa")

	if isDev {
		// 开发环境：先检查前端服务是否可用
		target := fmt.Sprintf("http://localhost:%d", frontendPort)

		// 检查前端服务是否可用
		client := http.Client{
			Timeout: 100 * time.Millisecond,
		}
		_, err := client.Get(target)
		if err != nil {
			app.serve.Logger.Warn(fmt.Sprintf(
				"Frontend development server is not running, please start it first (port: %d)",
				frontendPort,
			))
			// 返回提示信息的中间件
			app.serve.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					if strings.HasPrefix(c.Path(), "/rpc") {
						return next(c)
					}
					return c.JSON(http.StatusServiceUnavailable, map[string]string{
						"error": fmt.Sprintf("Frontend development server is not running on port %d", frontendPort),
						"hint":  "Please start the frontend development server first",
					})
				}
			})
			return
		}

		// 前端服务可用，设置代理
		proxy := middleware.ProxyWithConfig(middleware.ProxyConfig{
			Skipper: func(c echo.Context) bool {
				// 跳过 /rpc 路径的请求
				return strings.HasPrefix(c.Path(), "/rpc")
			},
			Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
				{
					URL: util.Must(url.Parse(target)),
				},
			}),
		})
		app.serve.Use(proxy)
	} else {
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
				if strings.HasPrefix(c.Path(), "/rpc") {
					return nil
				}
				return c.File(filepath.Join(distPath, "index.html"))
			})
		}
	}
}
