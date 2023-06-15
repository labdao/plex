package ipfs

import (
	"fmt"
	"io"
	"os"

	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/labdao/plex/internal/bacalhau"
)

func DeriveIpfsNodeUrl() (string, error) {
	bacalApiHost := bacalhau.GetBacalhauApiHost()
	ipfsUrl := fmt.Sprintf("http://%s:5001", bacalApiHost)
	return ipfsUrl, nil
}

func AddDir(dirPath string) (cid string, err error) {
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

func AddFile(filePath string) (cid string, err error) {
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
