package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "plex",
	Short: "Plex is a CLI application for running scientific workflows on peer to peer compute",
	Long:  `Plex is a CLI application for running scientific workflows on peer to peer compute. Complete documentation is available at https://docs.labdao.xyz/`,
	Run: func(cmd *cobra.Command, args []string) {
		dry := true
		upgradePlexVersion(dry)

		fmt.Println("Type ./plex --help to see commands")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
