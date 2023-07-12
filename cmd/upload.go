package cmd

import (
	"fmt"
	"os"

	"github.com/labdao/plex/internal/ipfs"
	"github.com/spf13/cobra"
)

var (
	localPath string
	wrapFile  bool
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file or directory to IPFS",
	Long:  `Upload and pins a file or directory to IPFS. Default wraps single files in a directory before uploading. (50MB max). A`,
	Run: func(cmd *cobra.Command, args []string) {
		dry := true
		upgradePlexVersion(dry)

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
			if wrapFile {
				cid, err = ipfs.WrapAndPinFile(localPath)
			} else {
				cid, err = ipfs.PinFile(localPath)
			}
		}

		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		// Do not change this output format, it is used by the Python wrapper
		fmt.Println("Uploaded CID: ", cid)
	},
}

func init() {
	uploadCmd.Flags().StringVarP(&localPath, "path", "p", "", "Local file path to the file or dir to upload")
	uploadCmd.Flags().BoolVarP(&wrapFile, "wrap", "w", true, "Wrap single files in a directory before uploading (default true)")

	rootCmd.AddCommand(uploadCmd)
}
