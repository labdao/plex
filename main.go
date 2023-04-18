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
		fmt.Println("Fill out this form to have an access token sent to your email: https://whe68a12b61.typeform.com/to/PpbO2HYf")
		os.Exit(1)
	} else if expectedToken != accessToken {
		fmt.Println("PLEX_ACCESS_TOKEN is incorrect")
		os.Exit(1)
	}

	// Env settings
	bacalApiHost, exists := os.LookupEnv("BACALHAU_API_HOST")
	if exists {
		fmt.Println("Using BACALHAU_API_HOST:", bacalApiHost)
	} else {
		fmt.Println("BACALHAU_API_HOST not set, using default host")
	}

	toolPath := flag.String("tool", "", "tool path")
	inputDir := flag.String("input-dir", "", "input directory path")
	ioJsonPath := flag.String("input-io", "", "IO JSON path")
	verbose := flag.Bool("verbose", false, "show verbose debugging logs")

	// required flags
	app := flag.String("app", "", "Application name")

	// optional flags
	appConfigsFilePath := flag.String("app-configs", "config/app.jsonl", "App Configurations file")
	layers := flag.Int("layers", 2, "Number of layers to search in the directory path")
	concurrency := flag.Int("concurrency", 1, "How many IO entries to run at once")
	memory := flag.Int("memory", 0, "Memory for job in GB, 0 autopicks a value")
	local := flag.Bool("local", false, "Use Docker on local machine to run job instead of Bacalhau")
	dry := flag.Bool("dry", false, "Do not send request and just print Bacalhau cmd")
	gpu := flag.Bool("gpu", false, "Use GPU")
	network := flag.Bool("network", false, "All http requests during job runtime")

	flag.Parse()

	// process tool input to be relative to tools directory
	if *toolPath != "" && !strings.Contains(*toolPath, "/") && !strings.HasSuffix(*toolPath, ".json") {
		*toolPath = filepath.Join("tools", *toolPath+".json")
	}
	fmt.Println("toolPath", *toolPath)

	if *toolPath != "" {
		fmt.Println("Running IPWL tool path")
		plex.Run(*toolPath, *inputDir, *ioJsonPath, *verbose, *local, *concurrency, *layers)
	} else {
		// Env settings
		bacalApiHost, exists := os.LookupEnv("BACALHAU_API_HOST")
		if exists {
			fmt.Println("Using BACALHAU_API_HOST:", bacalApiHost)
		} else {
			fmt.Println("BACALHAU_API_HOST not set, using default host")
		}

		if *app != "" {
			fmt.Println("WARNING: The -app flag is being deprecated and will be removed by v0.7.0. Please use -tool instead.")
		}

		// print the values of the flags
		fmt.Println("## User input ##")
		fmt.Println("Provided application name:", *app)
		fmt.Println("Provided directory path:", *inputDir)
		fmt.Println("Using GPU:", *gpu)
		fmt.Println("Using Network:", *network)

		fmt.Println("## Default parameters ##")
		fmt.Println("Using app configs:", *appConfigsFilePath)
		fmt.Println("Setting layers to:", *layers)

		plex.Execute(*app, *inputDir, *appConfigsFilePath, *layers, *memory, *local, *gpu, *network, *dry)
	}
}
