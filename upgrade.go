package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Masterminds/semver"
)

var CURRENT_PLEX_VERSION = "v0.4.1"

func downloadAndDecompressBinary(url, destination string) error {
	fmt.Printf("Downloading binary from: %s\n", url)
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

	tarReader := tar.NewReader(gzipReader)

	var fileContent []byte

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if filepath.Base(header.Name) == "plex" {
			fileContent, err = io.ReadAll(tarReader)
			if err != nil {
				return err
			}
			break
		}
	}

	err = ioutil.WriteFile(destination, fileContent, 0755)
	if err != nil {
		return err
	}

	return nil
}

func upgradePlexVersion() error {
	// auto update to latest plex version
	localReleaseVersion, err := semver.NewVersion(CURRENT_PLEX_VERSION)
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
		fmt.Printf("OS detected: %s, Arch detected: %s\n", userOS, userArch)

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

		fmt.Printf("Plex updated successfully to (v%s). Please re-run your desired command.\n", latestReleaseVersion)
		os.Exit(1)
	} else {
		fmt.Printf("Plex version (v%s) up to date.\n", localReleaseVersion)
	}

	return nil
}
