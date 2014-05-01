package jsonconfig

import (
	"strings"
	"testing"
)

func Test_SimpleCreate(t *testing.T) {
	json := `{"val":1}`
	o, err := Create(strings.NewReader(json))
	if err != nil {
		t.Fatal(err)
	}

	val := o.RequiredInt("val")
	if val != 1 {
		t.Errorf("Expected 1, got %d", val)
	}
}
