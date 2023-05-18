package ipfs

import (
	"fmt"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/labdao/plex/internal/bacalhau"
)

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

// returns the CID of the file
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
