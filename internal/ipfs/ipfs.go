package ipfs

import (
	"fmt"
	"io"
	"os"

	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/labdao/plex/internal/bacalhau"
)

func DownloadFromIPFS(ipfsNodeUrl, cid, filepath string) error {
	// Connect to the local IPFS node
	sh := shell.NewShell(ipfsNodeUrl)

	// Use the cat method to get the file with the specified CID
	fmt.Println("Cid")
	fmt.Println(cid)
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

func DeriveIpfsNodeUrl() (string, error) {
	bacalApiHost := bacalhau.GetBacalhauApiHost()
	ipfsUrl := fmt.Sprintf("http://%s:5001", bacalApiHost)
	return ipfsUrl, nil
}

func AddDirHttp(ipfsNodeUrl, dirPath string) (cid string, err error) {
	sh := shell.NewShell(ipfsNodeUrl)
	cid, err = sh.AddDir(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return cid, err
	}
	return cid, err
}

// Used to generate the IPFS cid for each file; each cid added to IO struct
func addFileHttp(ipfsNodeUrl, filePath string) (cid string, err error) {
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

// returns the cid of a directory
func GetDirCid(dirPath string) (string, error) {
	ipfsNodeUrl, err := DeriveIpfsNodeUrl()
	if err != nil {
		return "", fmt.Errorf("error deriving IPFS node URL: %w", err)
	}
	cid, err := AddDirHttp(ipfsNodeUrl, dirPath)
	if err != nil {
		return "", fmt.Errorf("error adding directory to IPFS: %w", err)
	}
	return cid, nil
}

// returns the CID of a file
func GetFileCid(filePath string) (string, error) {
	ipfsNodeUrl, err := DeriveIpfsNodeUrl()
	if err != nil {
		return "", fmt.Errorf("error deriving IPFS node URL: %w", err)
	}
	cid, err := addFileHttp(ipfsNodeUrl, filePath)
	if err != nil {
		return "", fmt.Errorf("error adding file to IPFS: %w", err)
	}
	return cid, nil
}

func IsValidCID(cidStr string) bool {
	_, err := cid.Decode(cidStr)
	if err != nil {
		// Invalid CID
		return false
	}
	return true
}
