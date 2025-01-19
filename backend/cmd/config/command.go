package config

import (
	"app/backend/core"
	"fmt"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long:  `Configure and manage application settings`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app, ok := cmd.Context().Value(core.CoreAppKey).(*core.App)
		if !ok {
			return fmt.Errorf("failed to get core app from context")
		}
		return runCommand(app)
	},
}

func init() {
	// 这里可以添加子命令
	// Command.AddCommand(getCmd)
	// Command.AddCommand(setCmd)
}

func runCommand(coreApp *core.App) error {
	settings := coreApp.GetConfig().AllSettings()
	fmt.Printf("Current Configuration:\n")
	for key, value := range settings {
		fmt.Printf("%s: %v\n", key, value)
	}
	return nil
}
