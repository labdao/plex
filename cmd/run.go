package cmd

import (
	"fmt"
	"os"

	"github.com/labdao/plex/internal/ipwl"
	"github.com/spf13/cobra"
)

var (
	ioJsonCid     string
	outputDir     string
	selector      string
	verbose       bool
	showAnimation bool
	maxTime       int
	annotations   *[]string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Processes an IO JSON",
	Long:  `Processes an IO JSON on Bacalhau and IPFS`,
	Run: func(cmd *cobra.Command, args []string) {
		dry := true
		upgradePlexVersion(dry)

		_, _, err := ipwl.RunIO(ioJsonCid, outputDir, selector, verbose, showAnimation, maxTime, *annotations)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	runCmd.Flags().StringVarP(&ioJsonCid, "ioJsonCid", "i", "", "IPFS CID of IO JSON")
	runCmd.Flags().StringVarP(&outputDir, "outputDir", "o", "", "Output directory")
	runCmd.Flags().StringVarP(&selector, "selector", "s", "", "Bacalhau Selector (label query) to filter nodes on which this job can be executed, supports '=', '==', and '!='.(e.g. -s key1=value1,key2=value2). Matching objects must satisfy all of the specified label constraints")
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	runCmd.Flags().BoolVarP(&showAnimation, "showAnimation", "", true, "Show job processing animation")
	runCmd.Flags().IntVarP(&maxTime, "maxTime", "m", 60, "Maximum time (min) to run a job")
	annotations = runCmd.Flags().StringArrayP("annotations", "a", []string{}, "Annotations to add to Bacalhau job")

	rootCmd.AddCommand(runCmd)
}
