package ipfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/labdao/plex/internal/bacalhau"
)

func DeriveIpfsNodeUrl() (string, error) {
	bacalApiHost := bacalhau.GetBacalhauApiHost()

	// If bacalApiHost already starts with http:// or https://, don't prepend it.
	if strings.HasPrefix(bacalApiHost, "http://") || strings.HasPrefix(bacalApiHost, "https://") {
		ipfsUrl := fmt.Sprintf("%s:5001", bacalApiHost)
		return ipfsUrl, nil
	}

	ipfsUrl := fmt.Sprintf("http://%s:5001", bacalApiHost)
	return ipfsUrl, nil
}

func PinDir(dirPath string) (cid string, err error) {
	ipfsNodeUrl, err := DeriveIpfsNodeUrl()
	if err != nil {
		return "", err
	}
	sh := shell.NewShell(ipfsNodeUrl)
	cid, err = sh.AddDir(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return cid, err
	}
	return cid, err
}

func PinFile(filePath string) (cid string, err error) {
	ipfsNodeUrl, err := DeriveIpfsNodeUrl()
	if err != nil {
		return "", err
	}
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
	ipfsNodeUrl, err := DeriveIpfsNodeUrl()
	if err != nil {
		return "", err
	}
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
	ipfsNodeUrl, err := DeriveIpfsNodeUrl()
	if err != nil {
		return err
	}
	sh := shell.NewShell(ipfsNodeUrl)

	// Construct the full directory path where the CID content will be downloaded
	downloadPath := path.Join(directory, cid)

	// Use the Get method to download the file or directory with the specified CID
	err = sh.Get(cid, downloadPath)
	if err != nil {
		return err
	}

	return nil
}

func DownloadToTempDir(cid string) (string, error) {
	ipfsNodeUrl, err := DeriveIpfsNodeUrl()
	if err != nil {
		return "", err
	}
	sh := shell.NewShell(ipfsNodeUrl)

	// Get the system's temporary directory
	tempDir := os.TempDir()

	// Construct the full directory path where the CID content will be downloaded
	downloadPath := path.Join(tempDir, cid)

	// Use the Get method to download the file or directory with the specified CID
	err = sh.Get(cid, downloadPath)
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
	ipfsNodeUrl, err := DeriveIpfsNodeUrl()
	if err != nil {
		return err
	}
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
	// Create a temporary directory
	tempDir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		return "", err
	}

	// Generate a random file name
	tempFilePath := filepath.Join(tempDir, fileName)

	// Download the file from IPFS to the temporary file
	err = DownloadFileContents(cid, tempFilePath)
	if err != nil {
		return "", err
	}

	// Return the path to the temporary file
	return tempFilePath, nil
}

func IsValidCID(cidStr string) bool {
	// implement
	return true
}

func IsDirectory(cidStr string) (bool, error) {
	filePath, err := DownloadToTempDir(cidStr)
	if err != nil {
		return false, err
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	if fileInfo.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

func IsImage(filepath string) bool {
	lowercaseFilePath := strings.ToLower(filepath)
	return strings.HasSuffix(lowercaseFilePath, ".jpg") ||
		strings.HasSuffix(lowercaseFilePath, ".jpeg") ||
		strings.HasSuffix(lowercaseFilePath, ".png") ||
		strings.HasSuffix(lowercaseFilePath, ".gif")
}
