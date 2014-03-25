package applications

import "testing"

func TestGetList(t *testing.T) {
	l := GetList()
	if len(l) != 2 {
		t.Fatalf("Unexpected number of items in the list")
	}
	if l[0].Name != "Test node backend" {
		t.Fatalf("Unexpected item - %s", l[0].Name)
	}
}
