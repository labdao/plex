package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/spf13/cobra"
)

var (
	toolPath string
	inputDir string
	layers   int
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates and pins an IO JSON",
	Long:  `Creates and pins an IO JSON`,
	Run: func(cmd *cobra.Command, args []string) {
		CreateIO(toolPath, inputDir, layers)
	},
}

func CreateIO(toolpath, inputDir string, layers int) {
	tempDirPath, err := ioutil.TempDir("", uuid.New().String())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Temporary directory created:", tempDirPath)
	defer os.RemoveAll(tempDirPath)

	var ioEntries []ipwl.IO
	fmt.Println("Reading tool config: ", toolPath)
	toolConfig, err := ipwl.ReadToolConfig(toolPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Creating IO entries from input directory: ", inputDir)
	ioEntries, err = ipwl.CreateIOJson(inputDir, toolConfig, toolPath, layers)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	ioJsonPath = path.Join(tempDirPath, "io.json")
	err = ipwl.WriteIOList(ioJsonPath, ioEntries)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Initialized IO file at: ", ioJsonPath)

	cid, err := ipfs.GetFileCid(ioJsonPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Initial IO file CID: ", cid)
}

func init() {
	createCmd.Flags().StringVarP(&toolPath, "toolPath", "t", "", "Path to the tool JSON file")
	createCmd.Flags().StringVarP(&inputDir, "inputDir", "i", "", "Directory containing input files")
	createCmd.Flags().IntVarP(&layers, "layers", "l", 1, "Number of layers to search input directory")

	rootCmd.AddCommand(createCmd)
}
