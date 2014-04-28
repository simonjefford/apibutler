package metadata

import "testing"

func TestGetList(t *testing.T) {
	l := GetApplicationsList()
	if len(l) != 2 {
		t.Fatalf("Unexpected number of items in the list")
	}
}
