package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/labdao/plex/cmd/plex"
)

func downloadAndDecompressBinary(url, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, gzipReader)
	return err
}

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

		userOS := runtime.GOOS
		userArch := runtime.GOARCH

		oldBinaryPath := os.Args[0]
		backupBinaryPath := fmt.Sprintf("%s.bak", oldBinaryPath)
		err := os.Rename(oldBinaryPath, backupBinaryPath)
		if err != nil {
			fmt.Printf("Error renaming the current binary: %v\n", err)
			os.Exit(1)
		}

		binaryURL := fmt.Sprintf("https://github.com/labdao/plex/releases/download/%s/plex_%s_%s_%s.tar.gz", latestReleaseVersionStr, strings.TrimPrefix(latestReleaseVersionStr, "v"), userOS, userArch)

		tempFile, err := ioutil.TempFile("", "plex_*")
		if err != nil {
			fmt.Printf("Error creating temporary file: %v\n", err)
			os.Exit(1)
		}
		defer os.Remove(tempFile.Name())

		err = downloadAndDecompressBinary(binaryURL, tempFile.Name())
		if err != nil {
			fmt.Printf("Error downloading and decompressing latest binary: %v\n", err)
			os.Exit(1)
		}

		out, err := os.Create("plex_new")
		if err != nil {
			fmt.Printf("Error creating new binary: %v\n", err)
			os.Exit(1)
		}
		defer out.Close()

		_, err = io.Copy(out, tempFile)
		if err != nil {
			fmt.Printf("Error writing to new binary: %v\n", err)
			os.Exit(1)
		}

		err = os.Chmod("plex_new", 0755)
		if err != nil {
			fmt.Printf("Error setting permissions on new binary: %v\n", err)
			os.Exit(1)
		}

		err = os.Rename("plex_new", "plex")
		if err != nil {
			fmt.Printf("Error renaming new binary: %v\n", err)
			os.Exit(1)
		}

		err = os.Remove(backupBinaryPath)
		if err != nil {
			fmt.Printf("Error removing backup binary: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Plex updated successfully to (v%s). Running new binary...\n", latestReleaseVersion)
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
	toolPath := flag.String("tool-path", "", "tool path")
	inputDir := flag.String("input-dir", "", "input directory path")

	// required flags
	app := flag.String("app", "", "Application name")

	// optional flags
	appConfigsFilePath := flag.String("app-configs", "config/app.jsonl", "App Configurations file")
	layers := flag.Int("layers", 2, "Number of layers to search in the directory path")
	memory := flag.Int("memory", 0, "Memory for job in GB, 0 autopicks a value")
	local := flag.Bool("local", false, "Use Docker on local machine to run job instead of Bacalhau")
	dry := flag.Bool("dry", false, "Do not send request and just print Bacalhau cmd")
	gpu := flag.Bool("gpu", false, "Use GPU")
	network := flag.Bool("network", false, "All http requests during job runtime")
	flag.Parse()

	fmt.Println("toolPath", *toolPath)

	if *toolPath != "" {
		fmt.Println("Running IPWL tool path")
		plex.Run(*toolPath, *inputDir)
	} else {
		// Env settings
		bacalApiHost, exists := os.LookupEnv("BACALHAU_API_HOST")
		if exists {
			fmt.Println("Using BACALHAU_API_HOST:", bacalApiHost)
		} else {
			fmt.Println("BACALHAU_API_HOST not set, using default host")
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
