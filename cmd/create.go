package cmd

import (
	"fmt"
	"io/ioutil"
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
	autoRun  bool
	layers   int
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates and pins an IO JSON",
	Long:  `Creates and pins an IO JSON`,
	Run: func(cmd *cobra.Command, args []string) {
		cid, err := CreateIO(toolPath, inputDir, layers)
		if err != nil {
			fmt.Println("Error:", err)
		}
		if autoRun {
			PlexRun(cid, outputDir, verbose, showAnimation, concurrency, annotations)
		}
	},
}

func CreateIO(toolpath, inputDir string, layers int) (string, error) {
	tempDirPath, err := ioutil.TempDir("", uuid.New().String())
	if err != nil {
		return "", err
	}

	fmt.Println("Temporary directory created:", tempDirPath)
	defer os.RemoveAll(tempDirPath)

	var ioEntries []ipwl.IO
	fmt.Println("Reading tool config: ", toolPath)
	toolConfig, err := ipwl.ReadToolConfig(toolPath)
	if err != nil {
		return "", err
	}

	fmt.Println("Creating IO entries from input directory: ", inputDir)
	ioEntries, err = ipwl.CreateIOJson(inputDir, toolConfig, toolPath, layers)
	if err != nil {
		return "", err
	}

	ioJsonPath := path.Join(tempDirPath, "io.json")
	err = ipwl.WriteIOList(ioJsonPath, ioEntries)
	if err != nil {
		return "", err
	}
	fmt.Println("Initialized IO file at: ", ioJsonPath)

	cid, err := ipfs.PinFile(ioJsonPath)
	if err != nil {
		return "", nil
	}

	// The Python SDK string matches here so make sure to change in both places
	fmt.Println("Initial IO JSON file CID: ", cid)
	return cid, nil
}

func init() {
	createCmd.Flags().StringVarP(&toolPath, "toolPath", "t", "", "Path to the tool JSON file")
	createCmd.Flags().StringVarP(&inputDir, "inputDir", "i", "", "Directory containing input files")
	createCmd.Flags().BoolVarP(&autoRun, "autoRun", "", false, "Auto submit the IO to plex run")
	createCmd.Flags().IntVarP(&layers, "layers", "", 2, "Number of layers to search input directory")
	createCmd.Flags().StringVarP(&outputDir, "outputDir", "o", "", "Output directory")
	createCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	createCmd.Flags().BoolVarP(&showAnimation, "showAnimation", "", true, "Show job processing animation")
	createCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent operations")
	createCmd.Flags().StringArrayP("annotations", "a", []string{}, "Annotations to add to Bacalhau job")

	rootCmd.AddCommand(createCmd)
}
