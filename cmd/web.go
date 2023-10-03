package cmd

import (
	"github.com/labdao/plex/gateway"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Runs the Gateway web app",
	Long:  `Runs the Gateway web app`,
	Run: func(cmd *cobra.Command, args []string) {
		dry := true
		upgradePlexVersion(dry)
		gateway.ServeWebApp()
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}
