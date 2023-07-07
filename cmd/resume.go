package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/spf13/cobra"
)

var (
	ioJsonPath string
	retry      bool
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resumes processing an IO JSON",
	Long:  `Resumes processing a local IO JSON in a plex working directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		dry := true
		upgradePlexVersion(dry)

		_, err := Resume(ioJsonPath, outputDir, verbose, showAnimation, retry, concurrency, annotations)
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func Resume(ioJsonFilePath, outputDir string, verbose, showAnimation, retry bool, concurrency int, annotations []string) (completedIoJsonCid string, err error) {
	fmt.Println("Continuing to process IO JSON file at: ", ioJsonPath)
	fmt.Println("Processing IO Entries")
	workDirPath := filepath.Dir(ioJsonFilePath)
	ipwl.ProcessIOList(workDirPath, ioJsonPath, retry, verbose, showAnimation, concurrency, annotations)
	fmt.Printf("Finished processing, results written to %s\n", ioJsonPath)
	completedIoJsonCid, err = ipfs.PinFile(ioJsonPath)
	if err != nil {
		return completedIoJsonCid, err
	}

	// The Python SDK string matches here so make sure to change in both places
	fmt.Println("Completed IO JSON CID:", completedIoJsonCid)
	return completedIoJsonCid, err
}

func init() {
	resumeCmd.Flags().StringVarP(&ioJsonPath, "ioJsonPath", "i", "", "Local file path to IO JSON")
	resumeCmd.Flags().StringVarP(&outputDir, "outputDir", "o", "", "Output directory")
	resumeCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	resumeCmd.Flags().BoolVarP(&showAnimation, "showAnimation", "", true, "Show job processing animation")
	resumeCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent operations")
	resumeCmd.Flags().StringArrayP("annotations", "a", []string{}, "Annotations to add to Bacalhau job")
	resumeCmd.Flags().BoolVarP(&retry, "retry", "", true, "Retry failed jobs")

	rootCmd.AddCommand(resumeCmd)
}
