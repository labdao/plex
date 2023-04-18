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
	fmt.Printf("added %s", cid)
	return cid, err
}
