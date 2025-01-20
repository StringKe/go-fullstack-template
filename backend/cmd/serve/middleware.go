package serve

import (
	"app/backend/pkg/util"
	"net/http"
	"net/url"
	"strings"

	"connectrpc.com/vanguard"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func transcoderPrefixHandler(prefix string, transcoder *vanguard.Transcoder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		transcoder.ServeHTTP(w, r)
	})
}

func proxyFrontend(target string) echo.MiddlewareFunc {
	return middleware.ProxyWithConfig(middleware.ProxyConfig{
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
}
