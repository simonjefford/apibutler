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

func Test_Validate(t *testing.T) {
	json := `{"intval": 1, "stringval": "foo", "objval": {"foo": 42}}`
	o, err := Create(strings.NewReader(json))
	if err != nil {
		t.Fatal(err)
	}

	o.RequiredInt("intval")
	o.RequiredString("stringval")
	o.RequiredObject("objval")
	o.RequiredObject("missingobj")

	err = o.Validate()

	if err == nil {
		t.Fatal("No error was thrown")
	}

	if err.Error() != "Missing required config key \"missingobj\" (object)" {
		t.Errorf("Unexpected error %v", err)
	}
}
