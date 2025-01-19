package serve

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"app/backend/core"
	"app/backend/serve"

	"github.com/spf13/cobra"
)

var (
	Command = &cobra.Command{
		Use:   "serve",
		Short: "Start the app backend server",
		Long: `Start the app backend server that provides API and services.
The server can be configured using flags or config file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			app, ok := cmd.Context().Value(core.CoreAppKey).(*core.App)
			if !ok {
				return fmt.Errorf("failed to get core app from context")
			}
			return runCommand(app)
		},
	}
)

func init() {

}

func runCommand(app *core.App) error {
	serveApp, err := serve.NewServeApp(app)
	if err != nil {
		return err
	}

	// 创建一个可以被取消的上下文
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// 启动服务器并等待它完成或出错
	return serveApp.Start(ctx)
}
