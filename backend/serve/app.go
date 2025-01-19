package serve

import (
	"app/backend/core"
	"app/backend/pkg/util"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

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
	return app.serve.StartH2CServer(addr, h2s)
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
				"前端开发服务未启动，请先在前端目录运行开发服务器 (端口: %d)\n"+
					"Frontend development server is not running, please start it first (port: %d)",
				frontendPort, frontendPort,
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
			ModifyResponse: func(res *http.Response) error {
				// 允许跨域
				res.Header.Set("Access-Control-Allow-Origin", "*")
				return nil
			},
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
