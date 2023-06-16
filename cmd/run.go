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
	annotations   []string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the Run function",
	Long:  `Runs the Run function`,
	Run: func(cmd *cobra.Command, args []string) {
		_, _, err := PlexRun(ioJsonCid, outputDir, verbose, showAnimation, concurrency, annotations)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func PlexRun(ioJsonCid, outputDir string, verbose, showAnimation bool, concurrency int, annotations []string) (completedIoJsonCid, ioJsonPath string, err error) {
	// Create plex working directory
	id := uuid.New()
	var cwd string
	if outputDir != "" {
		absPath, err := filepath.Abs(outputDir)
		if err != nil {
			return completedIoJsonCid, ioJsonPath, err
		}
		cwd = absPath
	} else {
		cwd, err = os.Getwd()
		if err != nil {
			return completedIoJsonCid, ioJsonPath, err
		}
		cwd = path.Join(cwd, "jobs")
	}
	workDirPath := path.Join(cwd, id.String())
	err = os.MkdirAll(workDirPath, 0755)
	if err != nil {
		return completedIoJsonCid, ioJsonPath, err
	}
	fmt.Println("Created working directory: ", workDirPath)

	ioJsonPath = path.Join(workDirPath, "io.json")
	err = ipfs.DownloadFile(ioJsonCid, ioJsonPath)
	if err != nil {
		return completedIoJsonCid, ioJsonPath, err
	}
	fmt.Println("Initialized IO file at: ", ioJsonPath)

	retry := false
	fmt.Println("Processing IO Entries")
	ipwl.ProcessIOList(workDirPath, ioJsonPath, retry, verbose, showAnimation, concurrency, annotations)
	fmt.Printf("Finished processing, results written to %s\n", ioJsonPath)
	completedIoJsonCid, err = ipfs.PinFile(ioJsonPath)
	if err != nil {
		return completedIoJsonCid, ioJsonPath, err
	}

	// The Python SDK string matches here so make sure to change in both places
	fmt.Println("Completed IO JSON CID:", completedIoJsonCid)
	return completedIoJsonCid, ioJsonPath, err
}

func init() {
	runCmd.Flags().StringVarP(&ioJsonCid, "ioJsonCid", "i", "", "IPFS CID of IO JSON")
	runCmd.Flags().StringVarP(&outputDir, "outputDir", "o", "", "Output directory")
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	runCmd.Flags().BoolVarP(&showAnimation, "showAnimation", "", true, "Show job processing animation")
	runCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent operations")
	runCmd.Flags().StringArrayP("annotations", "a", []string{}, "Annotations to add to Bacalhau job")

	rootCmd.AddCommand(runCmd)
}
