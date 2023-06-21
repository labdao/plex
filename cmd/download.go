package cmd

import (
	"fmt"
	"os"

	"github.com/labdao/plex/internal/ipfs"
	"github.com/spf13/cobra"
)

var (
	localOutputPath string
	cid             string
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download content from IPFS",
	Long:  `Downloads content from IPFS by providing a CID.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Downloading IPFS cid", cid, "to", localOutputPath)
		err := ipfs.DownloadToDirectory(cid, localOutputPath)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	downloadCmd.Flags().StringVarP(&cid, "cid", "c", "", "CID of the content to download")
	downloadCmd.Flags().StringVarP(&localOutputPath, "path", "p", "", "Local path at which to download the IPFS content")

	rootCmd.AddCommand(downloadCmd)
}
