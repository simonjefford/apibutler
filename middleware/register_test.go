package middleware

import (
	"testing"

	"fourth.com/apibutler/jsonconfig"

	"github.com/codegangsta/inject"
	"github.com/codegangsta/martini"
)

func ctor(cfg jsonconfig.Obj) (martini.Handler, error) {
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

func ctorWithConfig(cfg jsonconfig.Obj) (martini.Handler, error) {
	return func() string {
		return cfg.RequiredString("foo")
	}, nil
}

func Test_Configuration(t *testing.T) {
	err := Register("mw3", ctorWithConfig)

	if err != nil {
		t.Fatal(err)
	}

	h, err := Create("mw3", jsonconfig.Obj{
		"foo": "bar",
	})

	if err != nil {
		t.Fatal(err)
	}

	i := inject.New()
	v, err := i.Invoke(h)
	if v[0].String() != "bar" {
		t.Fatalf("Unexpected value %s", v[0])
	}
}
