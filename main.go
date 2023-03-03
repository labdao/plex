package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/labdao/plex/cmd/plex"
)

func main() {
	// Required env settings
	bacalApiHost, exists := os.LookupEnv("BACALHAU_API_HOST")
	if exists {
		fmt.Printf("Using BACALHAU_API_HOST: %s", bacalApiHost)
	} else {
		fmt.Println("Setting BACALHAU_API_HOST is required")
		os.Exit(1)
	}

	// required flags
	app := flag.String("app", "", "Application name")
	inputDir := flag.String("input-dir", "", "Input directory path")

	// optional flags
	appConfigsFilePath := flag.String("app-configs", "config/app.jsonl", "App Configurations file")
	layers := flag.Int("layers", 2, "Number of layers to search in the directory path")
	memory := flag.Int("memory", 0, "Memory for job in GB, 0 autopicks a value")
	dry := flag.Bool("dry", false, "Do not send request and just print Bacalhau cmd")
	gpu := flag.Bool("gpu", false, "Use GPU")
	network := flag.Bool("network", false, "All http requests during job runtime")
	flag.Parse()

	// print the values of the flags
	fmt.Println("## User input ##")
	fmt.Println("Provided application name:", *app)
	fmt.Println("Provided directory path:", *inputDir)
	fmt.Println("Using GPU:", *gpu)
	fmt.Println("Using Network:", *network)

	fmt.Println("## Default parameters ##")
	fmt.Println("Using app configs:", *appConfigsFilePath)
	fmt.Println("Setting layers to:", *layers)

	plex.Execute(*app, *inputDir, *appConfigsFilePath, *layers, *memory, *gpu, *network, *dry)
}
