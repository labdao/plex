package cmd

import (
	"fmt"
	"os"

	"github.com/labdao/plex/internal/ipfs"
	"github.com/spf13/cobra"
)

var (
	localPath string
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file or directory to IPFS",
	Long:  `Upload and pins a file or directory to IPFS. Will wrap single files in a directory before uploading. (50MB max). A`,
	Run: func(cmd *cobra.Command, args []string) {
		info, err := os.Stat(localPath)

		if os.IsNotExist(err) {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		var cid string
		if info.IsDir() {
			cid, err = ipfs.PinDir(localPath)
		} else {
			cid, err = ipfs.WrapAndPinFile(localPath)
		}

		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Uploaded CID: ", cid)
	},
}

func init() {
	uploadCmd.Flags().StringVarP(&localPath, "path", "p", "", "Local file path to the file or dir to upload")

	rootCmd.AddCommand(uploadCmd)
}
