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
