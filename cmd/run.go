package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labdao/plex/internal/ipwl"
	web3pkg "github.com/labdao/plex/internal/web3"
	"github.com/spf13/cobra"
)

var (
	ioJsonPath    string
	workDir       string
	outputDir     string
	verbose       bool
	retry         bool
	local         bool
	showAnimation bool
	concurrency   int
	web3          bool
	imageCID      string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the Run function",
	Long:  `Runs the Run function`,
	Run: func(cmd *cobra.Command, args []string) {
		PlexRun(ioJsonPath, workDir, outputDir, verbose, retry, local, showAnimation, concurrency, web3, imageCID)
	},
}

func PlexRun(ioJsonPath, workDir, outputDir string, verbose, retry, local, showAnimation bool, concurrency int, web3 bool, imageCID string) {
	// mint an NFT if web3 flag is set
	if web3 {
		web3pkg.MintNFT(ioJsonPath, imageCID)
		return
	}

	var workDirPath string
	var err error
	if workDir != "" && outputDir != "" {
		fmt.Println("Error: workDir and outputDir cannot be used at the same time")
		os.Exit(1)
	} else if workDir != "" {
		workDirPath = workDir
		fmt.Println("Resumed working directory: ", workDirPath)
	} else {
		// Create plex working directory
		id := uuid.New()
		var cwd string
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
		workDirPath = path.Join(cwd, id.String())
		err = os.MkdirAll(workDirPath, 0755)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Created working directory: ", workDirPath)
	}

	// first thing to generate io json and save to plex work dir
	var ioEntries []ipwl.IO
	if toolPath != "" {
		fmt.Println("Reading tool config: ", toolPath)
		toolConfig, err := ipwl.ReadToolConfig(toolPath)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		fmt.Println("Creating IO Entries from input directory: ", inputDir)
		ioEntries, err = ipwl.CreateIOJson(inputDir, toolConfig, toolPath, layers)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	} else if ioJsonPath != "" {
		fmt.Println("Reading IO Entries from: ", ioJsonPath)
		ioEntries, err = ipwl.ReadIOList(ioJsonPath)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	} else if workDir != "" {
		fmt.Println("Reading IO Entries from: ", path.Join(workDirPath, "io.json"))
		ioEntries, err = ipwl.ReadIOList(path.Join(workDirPath, "io.json"))
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		ipwl.PrintIOGraphStatus(ioEntries)
	} else {
		fmt.Println("Error: either -input-dir or -input-io is required")
		os.Exit(1)
	}

	ioJsonPath = path.Join(workDirPath, "io.json")
	err = ipwl.WriteIOList(ioJsonPath, ioEntries)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Initialized IO file at: ", ioJsonPath)

	fmt.Println("Processing IO Entries")
	fmt.Println(workDirPath)
	fmt.Println(ioJsonPath)
	ipwl.ProcessIOList(workDirPath, ioJsonPath, retry, verbose, local, showAnimation, concurrency)
	fmt.Printf("Finished processing, results written to %s\n", ioJsonPath)

}

func init() {
	runCmd.Flags().StringVarP(&toolPath, "toolPath", "t", "", "Path to the tool JSON file")
	runCmd.Flags().StringVarP(&inputDir, "inputDir", "i", "", "Directory containing input files")
	runCmd.Flags().IntVarP(&layers, "layers", "l", 1, "Number of layers to search input directory")
	runCmd.Flags().StringVarP(&ioJsonPath, "ioJsonPath", "j", "", "Path to IO JSON")
	runCmd.Flags().StringVarP(&workDir, "workDir", "w", "", "Working directory")
	runCmd.Flags().StringVarP(&outputDir, "outputDir", "o", "", "Output directory")
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	runCmd.Flags().BoolVarP(&retry, "retry", "r", false, "Retry on failure")
	runCmd.Flags().BoolVarP(&local, "local", "", false, "Use local resources")
	runCmd.Flags().BoolVarP(&showAnimation, "showAnimation", "s", false, "Show animation")
	runCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent operations")
	runCmd.Flags().BoolVarP(&web3, "web3", "W", false, "Use Web3")
	runCmd.Flags().StringVarP(&imageCID, "imageCID", "C", "", "Image CID")

	rootCmd.AddCommand(runCmd)
}
