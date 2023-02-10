package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
)

func main() {
	client, err := w3s.NewClient(
		w3s.WithEndpoint("https://api.web3.storage"),
		w3s.WithToken(os.Getenv("WEB3STORAGE_TOKEN")),
	)
	errorCheck(err)

	if len(os.Args) < 2 {
		fmt.Println("Error: Please specify a command (putFile, putDirectory, getFiles)")
		os.Exit(1)
	}

	command := strings.ToLower(os.Args[1])

	switch command {
	case "putfile":
		if len(os.Args) != 3 {
			fmt.Println("Error: Please specify a file path")
			os.Exit(1)
		}
		filePath := os.Args[2]
		file, err := os.Open(filePath)
		errorCheck(err)
		defer file.Close()
		putFile(client, file)
	case "putdirectory":
		if len(os.Args) != 3 {
			fmt.Println("Error: Please specify a directory path")
			os.Exit(1)
		}
		directoryPath := os.Args[2]
		putDirectory(client, directoryPath)
	case "getfiles":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please specify a CID")
			os.Exit(1)
		}
		cidString := os.Args[2]
		getFiles(client, cidString)
	}
}

func putFile(client w3s.Client, file fs.File, opts ...w3s.PutOption) (cid.Cid, error) {
	fmt.Printf("Uploading to IPFS via web3.storage... \n")
	cid, err := client.Put(context.Background(), file, opts...)
	if err != nil {
		return cid, err
	}
	fmt.Printf("CID: %s\n", cid)
	return cid, nil
}

func putDirectory(client w3s.Client, directoryPath string) (cid.Cid, error) {
	directory, err := os.Open(directoryPath)
	if err != nil {
		return cid.Cid{}, err
	}
	defer directory.Close()
	return putFile(client, directory)
}

func getFiles(client w3s.Client, cidStr string) error {
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

func errorCheck(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
