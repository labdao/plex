package ipfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/labdao/plex/internal/bacalhau"
)

func DeriveIpfsNodeUrl() (string, error) {
	bacalApiHost := bacalhau.GetBacalhauApiHost()
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

func DownloadFile(cid, filepath string) error {
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

func IsValidCID(cidStr string) bool {
	_, err := cid.Decode(cidStr)
	return err == nil
}
