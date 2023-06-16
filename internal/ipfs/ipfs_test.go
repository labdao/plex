package ipfs

import (
	"testing"
)

func TestPinDir(t *testing.T) {
	expectedCid := "QmWVKoVYBWHWdRLrL8Td5kUpqN2qH6zQ5piwtdCE1fjSYt"
	actualCid, err := PinDir("../../testdata/ipfs_test")
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
