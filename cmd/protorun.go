package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/spf13/cobra"
)

var protorunCmd = &cobra.Command{
	Use:   "protorun",
	Short: "Runs the protorun function",
	Long:  `Runs the protorun function`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running ProtoRun function...")

		// Create job working directory
		// change to temp directory, the Run function will create the local workDir
		var workDirPath string
		id := uuid.New()
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		workDirPath = path.Join(cwd, "jobs", id.String())
		err = os.MkdirAll(workDirPath, 0755)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Created working directory: ", workDirPath)

		// mcmenemy change the struct from FilePath to FileName
		// also update tool configs args to still work
		var ioEntries []ipwl.IO
		if toolPath != "" {
			// mcmenemy bonus move to CID
			fmt.Println("Reading tool config: ", toolPath)
			toolConfig, err := ipwl.ReadToolConfig(toolPath)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			fmt.Println("Creating IO entries from input directory: ", inputDir)
			ioEntries, err = ipwl.ProtoCreateIOJson(inputDir, toolConfig, toolPath, layers)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}

		// mcmenemy this will be moved to a temp directory
		var ioJsonPath string
		ioJsonPath = path.Join(workDirPath, "io.json")
		err = ipwl.WriteIOList(ioJsonPath, ioEntries)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Initialized IO file at: ", ioJsonPath)

		var cid string
		// mcmenemy change name of getfilecid since it pins
		cid, err = ipfs.GetFileCid(ioJsonPath)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		// mcmenemy return the cid don't just print it
		fmt.Println("Initial IO file CID: ", cid)

	},
}

func init() {
	protorunCmd.Flags().StringVarP(&toolPath, "toolPath", "t", "", "Path to the tool")
	protorunCmd.Flags().StringVarP(&inputDir, "inputDir", "i", "", "Directory of input")
	protorunCmd.Flags().IntVarP(&layers, "layers", "l", 1, "Number of layers")

	rootCmd.AddCommand(protorunCmd)
}
