package ipfs

import (
	"fmt"
	"os"
	"strings"

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

type FileEntry map[string]string
