package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
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

const (
	CurrentPlexVersion = "v0.7.1"
	ReleaseURL         = "https://api.github.com/repos/labdao/plex/releases/latest"
	ToolsURL           = "https://api.github.com/repos/labdao/plex/contents/tools?ref=main"
)

func getLatestReleaseVersionStr() (string, error) {
	resp, err := http.Get(ReleaseURL)
	if err != nil {
		return "", fmt.Errorf("Error getting latest release: %v", err)
	}
	defer resp.Body.Close()

	var responseMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		return "", fmt.Errorf("Error decoding latest release: %v", err)
	}

	htmlURL, ok := responseMap["html_url"].(string)
	if !ok {
		return "", fmt.Errorf("Error getting latest release html_url")
	}

	urlPartition := strings.Split(htmlURL, "/")
	latestReleaseVersionStr := urlPartition[len(urlPartition)-1]

	return latestReleaseVersionStr, nil
}

func getLocalFilesSHA(toolsFolderPath string) (map[string]string, error) {
	localFilesSHA := make(map[string]string)

	err := filepath.Walk(toolsFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, _ := filepath.Rel(toolsFolderPath, path)
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			h := sha1.New()
			if _, err := io.Copy(h, file); err != nil {
				return err
			}
			sha := fmt.Sprintf("%x", h.Sum(nil))
			localFilesSHA[relPath] = sha
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return localFilesSHA, nil
}

func downloadFile(url, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

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

func readManifestFile(manifestPath string) (map[string]string, error) {
	manifest := make(map[string]string)

	fileBytes, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return manifest, nil
		}
		return nil, err
	}

	err = json.Unmarshal(fileBytes, &manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func writeManifestFile(manifestPath string, manifest map[string]string) error {
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(manifestPath, manifestBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func upgradeToolsFolder(latestReleaseVersionStr string) error {
	toolsFolderPath := filepath.Join(".", "tools")
	manifestPath := filepath.Join(toolsFolderPath, ".manifest.json")

	// Check if the tools folder exists, and create it if it doesn't
	if _, err := os.Stat(toolsFolderPath); os.IsNotExist(err) {
		fmt.Println("Creating the tools folder...")
		err := os.Mkdir(toolsFolderPath, 0755)
		if err != nil {
			return err
		}
	}

	localFilesSHA, err := readManifestFile(manifestPath)
	if err != nil {
		return err
	}

	updatedManifest := make(map[string]string)

	resp, err := http.Get(ToolsURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var files []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return err
	}

	for _, file := range files {
		if fileType, ok := file["type"].(string); ok && fileType == "file" {
			fileName, ok := file["name"].(string)
			if !ok {
				continue
			}
			fileSHA, ok := file["sha"].(string)
			if !ok {
				continue
			}
			downloadURL, ok := file["download_url"].(string)
			if !ok {
				continue
			}

			localFilePath := filepath.Join(toolsFolderPath, fileName)
			_, err := os.Stat(localFilePath)
			fileExists := !os.IsNotExist(err)

			localSHA, exists := localFilesSHA[fileName]
			if !fileExists || !exists || localSHA != fileSHA {
				fmt.Printf("Downloading %s...\n", fileName)
				destination := filepath.Join(toolsFolderPath, fileName)
				if err := downloadFile(downloadURL, destination); err != nil {
					return err
				}
				updatedManifest[fileName] = fileSHA
			} else {
				updatedManifest[fileName] = localSHA
			}
		}
	}

	return writeManifestFile(manifestPath, updatedManifest)
}

func upgradePlexVersion() error {
	// Auto update to latest plex version
	localReleaseVersion, err := semver.NewVersion(CurrentPlexVersion)
	if err != nil {
		fmt.Printf("Error parsing local release version: %v\n", err)
		os.Exit(1)
	}

	latestReleaseVersionStr, err := getLatestReleaseVersionStr()
	if err != nil {
		fmt.Println("Error getting latest release version:", err)
		os.Exit(1)
	}

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
		fmt.Printf("- Operating System detected: %s\n- Chip Architecture detected: %s\n", userOS, userArch)

		// Upgrade tools folder
		err := upgradeToolsFolder(latestReleaseVersionStr)
		if err != nil {
			fmt.Println("Error upgrading tools folder:", err)
			os.Exit(1)
		}

		// Upgrade plex binary
		oldBinaryPath := os.Args[0]
		backupBinaryPath := fmt.Sprintf("%s.bak", oldBinaryPath)
		err = os.Rename(oldBinaryPath, backupBinaryPath)
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
			// Rollback the rename operation by restoring the backup binary
			if err := os.Rename(backupBinaryPath, oldBinaryPath); err != nil {
				fmt.Printf("Error restoring the backup binary: %v\n", err)
			}
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
