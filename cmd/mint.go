package cmd

import (
	"fmt"
	"os"

	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/web3"
	"github.com/spf13/cobra"
)

var (
	imageCid  string
	tokenName string
)

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mints a Proof of Science NFT",
	Long:  `Mints a Proof of Science NFT`,
	Run: func(cmd *cobra.Command, args []string) {
		dry := true
		upgradePlexVersion(dry)

		filePath, err := ipfs.DownloadFileToTemp(ioJsonCid, "io.json")
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		if tokenName == "" {
			tokenName = web3.GenerateTokenName()
		}
		if imageCid == "" {
			// default NFT image is glitchy labdao logo gif
			imageCid = "bafybeiba666bzbff5vu6rayvp5st2tk7tdltqnwjppzyvpljcycfhshdhq"
		}
		web3.MintNFT(filePath, imageCid, tokenName)
	},
}

func init() {
	mintCmd.Flags().StringVarP(&ioJsonCid, "ioJsonCid", "i", "", "IPFS CID of IO JSON")
	mintCmd.Flags().StringVarP(&imageCid, "imageCid", "", "", "IPFS CID of image")
	mintCmd.Flags().StringVarP(&tokenName, "tokenName", "", "", "Name of the NFT")

	rootCmd.AddCommand(mintCmd)
}
