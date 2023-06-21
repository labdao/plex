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
		fmt.Println("Uploading ", localPath, " to IPFS...")
		cid, err := uploadPath(localPath)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Uploaded CID: ", cid)
	},
}

func uploadPath(path string) (cid string, err error) {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return "", err
	}

	if err != nil {
		// Handle other possible errors here
		return "", err
	}

	if info.IsDir() {
		return ipfs.PinDir(path)
	} else {
		return ipfs.WrapAndPinFile(path)
	}
}

func init() {
	uploadCmd.Flags().StringVarP(&localPath, "path", "p", "", "Local file path to the file or dir to upload")

	rootCmd.AddCommand(uploadCmd)
}
