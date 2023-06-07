package plex

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labdao/plex/internal/ipwl"
	web3pkg "github.com/labdao/plex/internal/web3"
)

func Run(toolPath, inputDir, ioJsonPath, workDir, outputDir string, verbose, retry, local, showAnimation bool, concurrency, layers int, web3 bool) {
	// mint an NFT if web3 flag is set
	if web3 {
		fmt.Println("Minting NFT...")
		web3pkg.MintNFT(ioJsonPath)
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
