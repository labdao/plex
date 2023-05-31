package plex

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/labdao/plex/internal/ipfs"
	"github.com/labdao/plex/internal/ipwl"
	"github.com/labdao/plex/internal/web3"
)

func Run(toolPath, inputDir, ioJsonPath, workDir string, verbose, retry, local, showAnimation bool, concurrency, layers int, web3 bool) {
	// mint an NFT if web3 flag is set
	// ./plex -tool equibind -input-io /some/path/to/io.json -web3=true
	if web3 {
		fmt.Println("Minting NFT...")
		mintNFT(toolPath, ioJsonPath)
		return
	}

	var workDirPath string
	var err error
	if workDir != "" {
		workDirPath = workDir
		fmt.Println("Resumed working directory: ", workDirPath)
	} else {
		// Create plex working directory
		id := uuid.New()
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		workDirPath = path.Join(cwd, id.String())
		err = os.Mkdir(workDirPath, 0755)
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
	// if web3 {
	// 	mintNFT(toolPath, ioJsonPath)
	// }
}

func mintNFT(toolPath, ioJsonPath string) {
	// Build NFT metadata
	fmt.Println("Preparing NFT metadata...")
	metadata, err := web3.BuildTokenMetadata(toolPath, ioJsonPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Create a temporary file
	tempFile, err := ioutil.TempFile("", "metadata-*.json")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer os.Remove(tempFile.Name()) // clean up

	// Write the metadata to the temporary file
	_, err = tempFile.WriteString(metadata)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Close the file
	err = tempFile.Close()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Upload the metadata to IPFS and return the CID
	fmt.Println("Uploading NFT metadata to IPFS...")
	cid, err := ipfs.GetFileCid(tempFile.Name())
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Printf("NFT metadata uploaded to IPFS: ipfs://%s\n", cid)
}
