package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/spf13/cobra"
)

var (
	inputs           string
	scatteringMethod string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initilizes an IO JSON from Tool config and inputs",
	Long:  `This command initlizes an IO JSON based on the provided Tool and inputs.`,
	Run: func(cmd *cobra.Command, args []string) {
		dry := true
		upgradePlexVersion(dry)

		var kwargs map[string][]string
		err := json.Unmarshal([]byte(inputs), &kwargs)
		if err != nil {
			log.Fatal("Invalid inputs JSON:", err)
		}

		ioJson, err := ipwl.InitializeIo(toolPath, scatteringMethod, kwargs)
		if err != nil {
			log.Fatal(err)
		}

		// Convert the ioJson to bytes
		data, err := json.MarshalIndent(ioJson, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal the ioJson: %v", err)
		}

		// Create a temp file
		tempFile, err := ioutil.TempFile("", "io.json")
		if err != nil {
			log.Fatalf("Failed to create temporary file: %v", err)
		}

		// Write the data to the temp file
		_, err = tempFile.Write(data)
		if err != nil {
			log.Fatalf("Failed to write to temporary file: %v", err)
		}

		cid, err := ipfs.PinFile(tempFile.Name())
		if err != nil {
			log.Fatalf("Failed to pin the file: %v", err)
		}

		// Used by Python SDK do not change
		fmt.Println("Pinned IO JSON CID:", cid)

		// Remember to close the file
		if err := tempFile.Close(); err != nil {
			log.Fatalf("Failed to close the temporary file: %v", err)
		}
	},
}

func init() {
	initCmd.Flags().StringVarP(&toolPath, "toolPath", "t", "", "Path of the Tool config (can be a local path or an IPFS CID)")
	initCmd.Flags().StringVarP(&inputs, "inputs", "i", "{}", "Inputs in JSON format")
	initCmd.Flags().StringVarP(&scatteringMethod, "scatteringMethod", "", "{}", "Inputs in JSON format")

	rootCmd.AddCommand(initCmd)
}
