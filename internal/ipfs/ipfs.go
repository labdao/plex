package ipfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
)

func GetBacalhauApiHost() string {
	bacalApiHost, exists := os.LookupEnv("BACALHAU_API_HOST")
	plexEnv, _ := os.LookupEnv("PLEX_ENV")
	if exists {
		return bacalApiHost
	} else if plexEnv == "stage" {
		return "bacalhau.staging.labdao.xyz"
	} else {
		return "bacalhau.labdao.xyz"
	}
}

func DeriveIpfsNodeUrl() string {
	ipfsApiHost, exists := os.LookupEnv("IPFS_API_HOST")

	// If IPFS_API_HOST is not set, use BACALHAU_API_HOST
	if !exists {
		ipfsApiHost = GetBacalhauApiHost()
	}

	// If ipfsApiHost already starts with http:// or https://, don't prepend it.
	if strings.HasPrefix(ipfsApiHost, "http://") || strings.HasPrefix(ipfsApiHost, "https://") {
		ipfsUrl := fmt.Sprintf("%s:5001", ipfsApiHost)
		return ipfsUrl
	}
	ipfsUrl := fmt.Sprintf("http://%s:5001", ipfsApiHost)
	return ipfsUrl
}

func PinDir(dirPath string) (cid string, err error) {
	ipfsNodeUrl := DeriveIpfsNodeUrl()
	sh := shell.NewShell(ipfsNodeUrl)
	cid, err = sh.AddDir(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return cid, err
	}
	return cid, err
}

func PinFile(filePath string) (cid string, err error) {
	ipfsNodeUrl := DeriveIpfsNodeUrl()
	sh := shell.NewShell(ipfsNodeUrl)

	// Open the file and ensure it gets closed
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return cid, err
	}
	defer file.Close()

	// Pass the file reader to sh.Add()
	cid, err = sh.Add(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return cid, err
	}

	return cid, err
}

// wraps a file in a directory and adds it to IPFS
func WrapAndPinFile(filePath string) (cid string, err error) {
	ipfsNodeUrl := DeriveIpfsNodeUrl()
	sh := shell.NewShell(ipfsNodeUrl)

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return cid, err
	}

	tempDir, err := ioutil.TempDir("", "inputFile")
	if err != nil {
		return cid, err
	}

	defer os.RemoveAll(tempDir)

	_, fileName := filepath.Split(absPath)
	tempFilePath := filepath.Join(tempDir, fileName)

	srcFile, err := os.Open(filePath)
	if err != nil {
		return cid, err
	}
	defer srcFile.Close()

	destFile, err := os.Create(tempFilePath)
	if err != nil {
		return cid, err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return cid, err
	}

	cid, err = sh.AddDir(tempDir)
	if err != nil {
		return cid, err
	}

	return cid, err
}

func DownloadToDirectory(cid, directory string) error {
	ipfsNodeUrl := DeriveIpfsNodeUrl()
	sh := shell.NewShell(ipfsNodeUrl)

	// Use the Get method to download the file or directory with the specified CID
	err := sh.Get(cid, directory)
	if err != nil {
		return err
	}

	return nil
}

func DownloadToTempDir(cid string) (string, error) {
	ipfsNodeUrl := DeriveIpfsNodeUrl()
	sh := shell.NewShell(ipfsNodeUrl)

	// Get the system's temporary directory
	tempDir := os.TempDir()

	// Construct the full directory path where the CID content will be downloaded
	downloadPath := path.Join(tempDir, cid)

	// Use the Get method to download the file or directory with the specified CID
	err := sh.Get(cid, downloadPath)
	if err != nil {
		return "", err
	}

	return downloadPath, nil
}

func UnwrapAndDownloadFileContents(cid, outputFilePath string) error {
	// First download the CID content to a temporary file
	tempDirPath, err := DownloadToTempDir(cid)
	if err != nil {
		return err
	}

	// Ensure that the temporary directory is deleted after we are done
	defer os.RemoveAll(tempDirPath)

	onlyOneFile, tempFilePath, err := onlyOneFile(tempDirPath)
	if err != nil {
		return err
	}

	if !onlyOneFile {
		return fmt.Errorf("more than one file in the CID %s content", cid)
	}

	// Now copy the downloaded content to the output file path
	inputFile, err := os.Open(tempFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Ensure the directory exists
	outputDir := filepath.Dir(outputFilePath)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			return err
		}
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return err
	}

	return nil
}

func onlyOneFile(dirPath string) (bool, string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return false, "", err
	}

	var filePath string
	fileCount := 0
	for _, file := range files {
		if !file.IsDir() {
			fileCount++
			filePath = filepath.Join(dirPath, file.Name())
		}
	}

	if fileCount == 1 {
		return true, filePath, nil
	} else {
		return false, "", nil
	}
}

func DownloadFileContents(cid, filepath string) error {
	ipfsNodeUrl := DeriveIpfsNodeUrl()
	sh := shell.NewShell(ipfsNodeUrl)

	// Use the cat method to get the file with the specified CID
	reader, err := sh.Cat(cid)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Create the destination file
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the data from the IPFS file to the local file
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFileToTemp(cid, fileName string) (string, error) {
	tempDir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		log.Printf("Failed to create temporary directory: %v", err)
		return "", err
	}

	isDir, err := IsDirectory(cid)
	if err != nil {
		log.Printf("Failed to determine if CID %s is a directory: %v", cid, err)
		return "", err
	}

	if isDir {
		log.Printf("CID %s is a directory, downloading contents", cid)
		err = DownloadToDirectory(cid, tempDir)
		if err != nil {
			log.Printf("Failed to download directory with CID %s: %v", cid, err)
			return "", err
		}

		files, err := ioutil.ReadDir(tempDir)
		if err != nil {
			log.Printf("Failed to read temporary directory: %v", err)
			return "", err
		}

		for _, file := range files {
			if file.Name() == fileName {
				tempFilePath := filepath.Join(tempDir, fileName)
				log.Printf("Found file %s in downloaded directory", fileName)
				return tempFilePath, nil
			}
		}

		return "", fmt.Errorf("file %s not found in downloaded directory", fileName)
	} else {
		tempFilePath := filepath.Join(tempDir, fileName)
		log.Printf("Downloading file with CID %s to %s", cid, tempFilePath)

		err = DownloadFileContents(cid, tempFilePath)
		if err != nil {
			log.Printf("Failed to download file with CID %s: %v", cid, err)
			return "", err
		}

		log.Printf("Downloaded file with CID %s to %s", cid, tempFilePath)
		return tempFilePath, nil
	}
}

func IsValidCID(cidStr string) bool {
	_, err := cid.Decode(cidStr)
	return err == nil
}

func IsDirectory(cidStr string) (bool, error) {
	ipfsNodeUrl := DeriveIpfsNodeUrl()
	sh := shell.NewShell(ipfsNodeUrl)

	// Use the FilesList function to get information about the CID
	entries, err := sh.List(cidStr)

	// If an error is returned, it might be a file or an invalid CID
	if err != nil {
		// You can further check the error message if needed to differentiate between a file and an invalid CID
		return false, fmt.Errorf("error getting CID info from IPFS: %v", err)
	}

	// If there are multiple entries, it's definitely a directory
	if len(entries) >= 1 {
		return true, nil
	}

	// If there is one entry and it's a different CID then it is a direcotry with a single file
	if len(entries) == 1 && entries[0].Hash != cidStr {
		return true, nil
	}

	// Otherwise it is a file
	return false, nil
}

func IsImage(filepath string) bool {
	lowercaseFilePath := strings.ToLower(filepath)
	return strings.HasSuffix(lowercaseFilePath, ".jpg") ||
		strings.HasSuffix(lowercaseFilePath, ".jpeg") ||
		strings.HasSuffix(lowercaseFilePath, ".png") ||
		strings.HasSuffix(lowercaseFilePath, ".gif")
}

type FileEntry map[string]string

func getContents(sh *shell.Shell, cid string) ([]FileEntry, error) {
	links, err := sh.List(cid)
	if err != nil {
		return nil, err
	}

	var files []FileEntry
	for _, link := range links {
		if link.Type == shell.TDirectory {
			subFiles, err := getContents(sh, link.Hash)
			if err != nil {
				fmt.Println("Error getting subdirectory links:", err)
				continue
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, FileEntry{"filename": link.Name, "CID": link.Hash})
		}
	}

	return files, nil
}

func ListFilesInDirectory(cid string) ([]FileEntry, error) {
	ipfsNodeUrl := DeriveIpfsNodeUrl()
	sh := shell.NewShell(ipfsNodeUrl)
	return getContents(sh, cid)
}
