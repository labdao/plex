package ipfs

import (
	"testing"
)

func TestAddDirHttp(t *testing.T) {
	expectedCid := "QmWVKoVYBWHWdRLrL8Td5kUpqN2qH6zQ5piwtdCE1fjSYt"
	ipfsNodeUrl, err := DeriveIpfsNodeUrl()
	if err != nil {
		t.Fatalf("error deriving IPFS node Url: %v", err)
	}
	actualCid, err := AddDirHttp(ipfsNodeUrl, "../../testdata/ipfs_test")
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	if expectedCid != actualCid {
		t.Errorf(`unmatching cids
			expected CID: %s
			actual CID: %s`, expectedCid, actualCid,
		)
	}
}

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
