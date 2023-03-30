package ipwl

import (
	"fmt"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/web3-storage/go-w3s-client"
)

func ProcessIOList(ioList []) (string, error) {
	client, err := w3s.NewClient(
		w3s.WithEndpoint("https://api.web3.storage"),
		w3s.WithToken(os.Getenv("WEB3STORAGE_TOKEN")),
	)
	if err != nil {
		return "", err
	}
	inputDir, err := os.Open(inputDirPath)
	if err != nil {
		return "", err
	}
	cid, err := PutFile(client, inputDir)
	if err != nil {
		return cid.String(), err
	}
	return cid.String(), nil
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
