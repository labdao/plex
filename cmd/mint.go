package cmd

import (
	"fmt"
	"os"

	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/web3"
	"github.com/spf13/cobra"
)

var (
	imageCid string
)

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mints a Proof of Science NFT",
	Long:  `Mints a Proof of Science NFT`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, err := ipfs.DownloadFileToTemp(ioJsonCid, "io.json")
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		web3.MintNFT(filePath, imageCid)
	},
}

func init() {
	mintCmd.Flags().StringVarP(&ioJsonCid, "ioJsonCid", "i", "", "IPFS CID of IO JSON")
	mintCmd.Flags().StringVarP(&imageCid, "imageCid", "", "", "IPFS CID of image")

	rootCmd.AddCommand(mintCmd)
}
