package plex

import (
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/labdao/plex/internal/ipwl"
)

func Run(toolPath, inputDir, ioJsonPath string, verbose, local bool, concurrency, layers int) {
	// Create plex working directory
	id := uuid.New()
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	workDirPath := path.Join(cwd, id.String())
	err = os.Mkdir(workDirPath, 0755)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Created job directory: ", workDirPath)

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
	ipwl.ProcessIOList(ioEntries, workDirPath, ioJsonPath, verbose, local, concurrency)
	fmt.Printf("Finished processing, results written to %s", ioJsonPath)
}
