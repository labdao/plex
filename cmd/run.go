package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/spf13/cobra"
)

var (
	ioJsonCid     string
	outputDir     string
	verbose       bool
	showAnimation bool
	concurrency   int
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the Run function",
	Long:  `Runs the Run function`,
	Run: func(cmd *cobra.Command, args []string) {
		PlexRun(ioJsonCid, outputDir, verbose, showAnimation, concurrency)
	},
}

func PlexRun(ioJsonCid, outputDir string, verbose, showAnimation bool, concurrency int) {
	// Create plex working directory
	id := uuid.New()
	var cwd string
	var err error
	if outputDir != "" {
		absPath, err := filepath.Abs(outputDir)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		cwd = absPath
	} else {
		cwd, err = os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		cwd = path.Join(cwd, "jobs")
	}
	workDirPath := path.Join(cwd, id.String())
	err = os.MkdirAll(workDirPath, 0755)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Created working directory: ", workDirPath)

	ioJsonPath := path.Join(workDirPath, "io.json")
	ipfsNodeUrl, err := ipfs.DeriveIpfsNodeUrl()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("iojsoncid")
	fmt.Println(ioJsonCid)
	err = ipfs.DownloadFromIPFS(ipfsNodeUrl, ioJsonCid, ioJsonPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Initialized IO file at: ", ioJsonPath)

	retry := false
	fmt.Println("Processing IO Entries")
	ipwl.ProcessIOList(workDirPath, ioJsonPath, retry, verbose, showAnimation, concurrency)
	fmt.Printf("Finished processing, results written to %s\n", ioJsonPath)
}

func init() {
	runCmd.Flags().StringVarP(&inputDir, "inputDir", "i", "", "Directory containing input files")
	runCmd.Flags().IntVarP(&layers, "layers", "l", 1, "Number of layers to search input directory")
	runCmd.Flags().StringVarP(&ioJsonCid, "ioJsonPath", "j", "", "Path to IO JSON")
	runCmd.Flags().StringVarP(&outputDir, "outputDir", "o", "", "Output directory")
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	runCmd.Flags().BoolVarP(&showAnimation, "showAnimation", "s", false, "Show animation")
	runCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent operations")

	rootCmd.AddCommand(runCmd)
}
