package ipfs

import (
	"testing"
)

func TestAddFileHttp(t *testing.T) {
	// Derive the IPFS node URL
	ipfsNodeUrl, err := DeriveIpfsNodeUrl()
	if err != nil {
		t.Fatal(err)
	}

	// Add the specified file to the IPFS node
	cid, err := AddFileHttp(ipfsNodeUrl, "../../testdata/binding/abl/7n9g.pdb")
	if err != nil {
		t.Fatal(err)
	}

	// Check if the obtained CID matches the expected CID
	if cid != "QmbHbHLGXyB9FSBvroMeSVZW72fZnqXoopES82Q8qP2etw" {
		t.Fatalf("Obtained CID (%s) does not match the expected CID (%s)", cid, "QmbHbHLGXyB9FSBvroMeSVZW72fZnqXoopES82Q8qP2etw")
	}
}
