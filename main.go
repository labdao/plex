package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Masterminds/semver"
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

	// auto update to latest plex version
	localReleaseVersion, err := semver.NewVersion("v0.4.1")
	if err != nil {
		fmt.Printf("Error parsing local release version: %v\n", err)
		os.Exit(1)
	}
	releaseURL := "https://api.github.com/repos/labdao/plex/releases/latest"
	resp, err := http.Get(releaseURL)
	if err != nil {
		fmt.Println("Error getting latest release:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	var responseMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		fmt.Println("Error decoding latest release:", err)
		os.Exit(1)
	}
	htmlURL, ok := responseMap["html_url"].(string)
	if !ok {
		fmt.Println("Error getting latest release html_url")
		os.Exit(1)
	}
	urlPartition := strings.Split(htmlURL, "/")
	latestReleaseVersionStr := urlPartition[len(urlPartition)-1]
	latestReleaseVersion, err := semver.NewVersion(latestReleaseVersionStr)
	if err != nil {
		fmt.Printf("Error parsing latest release version: %v\n", err)
		os.Exit(1)
	}
	if localReleaseVersion.LessThan(latestReleaseVersion) {
		fmt.Printf("The version of plex you are running (v%s) is outdated.\n", localReleaseVersion)
		fmt.Printf("Updating to latest plex version (v%s)...\n", latestReleaseVersion)
		// TODO: update plex
		// need to redownload plex binary based on user's OS and architecture
		// and replace the current binary
	} else {
		fmt.Printf("Plex version (v%s) up to date.\n", localReleaseVersion)
	}

	// Env settings
	bacalApiHost, exists := os.LookupEnv("BACALHAU_API_HOST")
	if exists {
		fmt.Println("Using BACALHAU_API_HOST:", bacalApiHost)
	} else {
		fmt.Println("BACALHAU_API_HOST not set, using default host")
	}

	// required flags
	app := flag.String("app", "", "Application name")
	inputDir := flag.String("input-dir", "", "Input directory path")

	// optional flags
	appConfigsFilePath := flag.String("app-configs", "config/app.jsonl", "App Configurations file")
	layers := flag.Int("layers", 2, "Number of layers to search in the directory path")
	memory := flag.Int("memory", 0, "Memory for job in GB, 0 autopicks a value")
	local := flag.Bool("local", false, "Use Docker on local machine to run job instead of Bacalhau")
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

	plex.Execute(*app, *inputDir, *appConfigsFilePath, *layers, *memory, *local, *gpu, *network, *dry)
}
