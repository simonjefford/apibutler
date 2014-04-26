package middleware

import "testing"

func Test_MissingKeys(t *testing.T) {
	cfg := make(MiddlewareConfig)
	cfg["key1"] = "foo"
	err := cfg.CheckForMandatoryKeys("key1", "key2")
	if err == nil {
		t.Fatal("No errors")
	}
	if err.Error() != "The following keys are missing: key2" {
		t.Fatalf("Incorrect message: %s", err.Error())
	}
}

func Test_NoMissingKeys(t *testing.T) {
	cfg := make(MiddlewareConfig)
	cfg["key1"] = "foo"
	cfg["key2"] = "foo"
	err := cfg.CheckForMandatoryKeys("key1", "key2")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
