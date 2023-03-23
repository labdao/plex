package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/labdao/plex/cmd/plex"
)

func main() {
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

	// required flags
	tool := flag.String("tool", "", "Tool name")
	inputDir := flag.String("input-dir", "", "Input directory path")

	// ensuring backward compatability to v.0.4 and earlier
	tool := flag.String("app", "", "Tool name")

	// optional flags
	toolConfigsFilePath := flag.String("tool-config", "config", "Tool Configurations path")
	layers := flag.Int("layers", 2, "Number of layers to search in the directory path")
	memory := flag.Int("memory", 0, "Memory for job in GB, 0 autopicks a value")
	local := flag.Bool("local", false, "Use Docker on local machine to run job instead of Bacalhau")
	dry := flag.Bool("dry", false, "Do not send request and just print Bacalhau cmd")
	gpu := flag.Bool("gpu", false, "Use GPU")
	network := flag.Bool("network", false, "All http requests during job runtime")
	flag.Parse()

	// print the values of the flags
	fmt.Println("## User input ##")
	fmt.Println("Provided tool name:", *tool)
	fmt.Println("Provided directory path:", *inputDir)
	fmt.Println("Using GPU:", *gpu)
	fmt.Println("Using Network:", *network)

	fmt.Println("## Default parameters ##")
	fmt.Println("Using tool config path:", *toolConfigsFilePath)
	fmt.Println("Setting layers to:", *layers)

	plex.Execute(*tool, *inputDir, *toolConfigsFilePath, *layers, *memory, *local, *gpu, *network, *dry)
}
