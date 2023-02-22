package main

import (
	"flag"
	"fmt"

	"github.com/labdao/plex/cmd/plex"
)

func main() {
	// required flags
	app := flag.String("app", "", "Application name")
	inputDir := flag.String("input-dir", "", "Input directory path")
	// todo: needs to be a boolean flag
	gpu := flag.String("gpu", "true", "Use GPU")

	// optional flags
	appConfigsFilePath := flag.String("app-configs", "config/app.jsonl", "App Configurations file")
	layers := flag.Int("layers", 2, "number of layers to search in the directory path")
	flag.Parse()

	// print the values of the flags
	fmt.Println("## User input ##")
	fmt.Println("Provided application name:", *app)
	fmt.Println("Provided directory path:", *inputDir)
	fmt.Println("Using GPU:", *gpu)

	fmt.Println("## Default parameters ##")
	fmt.Println("Using app configs:", *appConfigsFilePath)
	fmt.Println("Setting layers to:", *layers)

	plex.Execute(*app, *inputDir, *gpu, *appConfigsFilePath, *layers)
}
