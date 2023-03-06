package ipfs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/web3-storage/go-w3s-client"
)

func DeriveIpfsNodeUrl() (string, error) {
	apiHost, exists := os.LookupEnv("BACALHAU_API_HOST")
	if !exists {
		return apiHost, errors.New("can not derive IPFS node url, BACALHAU_API_HOST not set")
	}
	ipfsUrl := fmt.Sprintf("http://%s:5001", apiHost)
	return ipfsUrl, nil
}

func PutFile(client w3s.Client, file fs.File, opts ...w3s.PutOption) (cid.Cid, error) {
	fmt.Printf("Uploading to IPFS via web3.storage... \n")
	cid, err := client.Put(context.Background(), file, opts...)
	if err != nil {
		return cid, err
	}
	fmt.Printf("CID: %s\n", cid)
	return cid, nil
}

func PutDirectory(client w3s.Client, directoryPath string) (cid.Cid, error) {
	directory, err := os.Open(directoryPath)
	if err != nil {
		return cid.Cid{}, err
	}
	defer directory.Close()
	return PutFile(client, directory)
}

func GetFiles(client w3s.Client, cidStr string) error {
	fmt.Printf("Retrieving files from IPFS... \n")

	cid, err := cid.Parse(cidStr)
	if err != nil {
		return err
	}

	res, err := client.Get(context.Background(), cid)
	if err != nil {
		return err
	}

	f, fsys, err := res.Files()
	if err != nil {
		return err
	}

	info, err := f.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		err = fs.WalkDir(fsys, "/", func(path string, d fs.DirEntry, err error) error {
			info, _ := d.Info()
			fmt.Printf("%s (%d bytes)\n", path, info.Size())
			return err
		})
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s (%d bytes)\n", cid.String(), info.Size())

	return nil
}

func createInputCID(inputDirPath string, cmd string) (string, error) {
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
