package serve

import (
	"fmt"

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

	err = serveApp.Start()
	if err != nil {
		return err
	}
	return nil
}
