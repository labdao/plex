package main

import (
	"os"
	"testing"

	cid "github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
)

func TestPutFile(t *testing.T) {
	client, err := w3s.NewClient(
		w3s.WithEndpoint("https://api.web3.storage"),
		// set your web3 storage token here
		w3s.WithToken(os.Getenv("WEB3STORAGE_TOKEN")),
	)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}

	expectedCidStr := "bafybeibht52e6gjbbt2qrrdsje527ijo64k2i5wp5wxzmj63ai6hsw7lpu"
	expectedCid, err := cid.Decode(expectedCidStr)
	if err != nil {
		t.Fatalf("error decoding expected CID: %v", err)
	}

	file, err := os.Open("test-directory/haiku2.txt")

	actualCid, err := putFile(client, file)

	if !expectedCid.Equals(actualCid) {
		t.Errorf(`unmatching cids
			expected CID: %s
			actual CID: %s`, expectedCid, actualCid,
		)
	}
}
