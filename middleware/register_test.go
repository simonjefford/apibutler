package middleware

import (
	"testing"

	"github.com/codegangsta/inject"
	"github.com/codegangsta/martini"
)

func ctor(cfg MiddlewareConfig) (martini.Handler, error) {
	return func() string {
		return "martini.Handler"
	}, nil
}

func Test_DuplicateRegistration(t *testing.T) {
	Register("mw", ctor)
	err := Register("mw", ctor)
	if err == nil {
		t.Fatal("Expected error was not returned on duplicate middleware registration")
	}
}

func Test_RegisterAndCreate(t *testing.T) {
	err := Register("mw2", ctor)

	if err != nil {
		t.Fatal(err)
	}

	h, err := Create("mw2", nil)

	i := inject.New()
	v, err := i.Invoke(h)
	if v[0].String() != "martini.Handler" {
		t.Fatal("Unexpected handler returned")
	}
}

func Test_UnknownMiddleware(t *testing.T) {
	h, err := Create("missing", nil)
	if h != nil || err == nil {
		t.Fatal("No error was raised, or a handler was returned.")
	}

	if err.Error() != "Unknown middleware: missing" {
		t.Fatalf("An unexpected error was returned: %v", err)
	}
}
