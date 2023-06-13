package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/labdao/plex/cmd/plex"
)

func main() {
	// check for new plex version
	upgradePlexVersion()

	// token access
	accessToken, exists := os.LookupEnv("PLEX_ACCESS_TOKEN")
	expectedToken := "mellon" // speak friend and enter
	if !exists {
		fmt.Println("PLEX_ACCESS_TOKEN is required")
		fmt.Println("Run export PLEX_ACCESS_TOKEN=<token>")
		fmt.Println("Fill out this form to have an access token sent to your email: https://airtable.com/shrfEDQj2fPffUge8")
		os.Exit(1)
	} else if expectedToken != accessToken {
		fmt.Println("PLEX_ACCESS_TOKEN is incorrect")
		os.Exit(1)
	}

	toolPath := flag.String("tool", "", "tool path")
	inputDir := flag.String("input-dir", "", "input directory path")
	ioJsonPath := flag.String("input-io", "", "IO JSON path")
	workDir := flag.String("work-dir", "", "PLEx working directory path")
	outputDir := flag.String("output-dir", "", "directory to store job results")
	verbose := flag.Bool("verbose", false, "show verbose debugging logs")
	layers := flag.Int("layers", 2, "Number of layers to search in the directory path")
	concurrency := flag.Int("concurrency", 1, "How many IO entries to run at once")
	local := flag.Bool("local", false, "Use Docker on local machine to run job instead of Bacalhau")
	showAnimation := flag.Bool("show-animation", true, "Show animation while Bacalhau job is running")
	retry := flag.Bool("retry", false, "Retry any io subgraphs that failed")
	proto := flag.Bool("proto", false, "Option to run in prototype mode")

	web3 := flag.Bool("web3", false, "Option to mint an NFT")
	image := flag.String("image", "", "Image to add to NFT metadata")

	flag.Parse()

	// process tool input to be relative to tools directory
	if *toolPath != "" && !strings.Contains(*toolPath, "/") && !strings.HasSuffix(*toolPath, ".json") {
		*toolPath = filepath.Join("tools", *toolPath+".json")
	}
	fmt.Println("toolPath", *toolPath)

	if *toolPath != "" {
		fmt.Println("Running IPWL tool path")
		fmt.Println("Warning: tool path support will be removed and moved to the Python SDK in the future")
		if *inputDir == "" && *web3 == false {
			fmt.Println("Input dir or web3 flag set to true is required when using the -tool option")
			os.Exit(1)
		}
		*retry = false // can only retry from an PLEx work dir not input directory input
		if *proto {
			plex.ProtoRun(*toolPath, *inputDir, *layers)
		} else {
			plex.Run(*toolPath, *inputDir, *ioJsonPath, *workDir, *outputDir, *verbose, *retry, *local, *showAnimation, *concurrency, *layers, *web3, *image)
		}
	} else if *ioJsonPath != "" {
		fmt.Println("Running IPWL io path")
		*retry = false // can only retry from an PLEx work dir not io json path input
		if *proto {
			plex.ProtoRun(*toolPath, *inputDir, *layers)
		} else {
			plex.Run(*toolPath, *inputDir, *ioJsonPath, *workDir, *outputDir, *verbose, *retry, *local, *showAnimation, *concurrency, *layers, *web3, *image)
		}
	} else if *workDir != "" {
		if *proto {
			plex.ProtoRun(*toolPath, *inputDir, *layers)
		} else {
			plex.Run(*toolPath, *inputDir, *ioJsonPath, *workDir, *outputDir, *verbose, *retry, *local, *showAnimation, *concurrency, *layers, *web3, *image)
		}
	} else {
		fmt.Println("Requirements invalid. Please run './plex -h' for help.")
	}
}
