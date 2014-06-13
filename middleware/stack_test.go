package middleware

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"fmt"

	"fourth.com/apibutler/jsonconfig"
	"fourth.com/apibutler/testhelpers"
	"github.com/codegangsta/martini"
)

func setup() {
	clearTable()
	Register("stack1", NewDefinition("stack1", func(obj jsonconfig.Obj) (martini.Handler, error) {
		return func(r http.ResponseWriter) {
			r.Header().Add("X-Stack1", "stack1")
		}, nil
	}))

	Register("stack2", NewDefinition("stack2", func(obj jsonconfig.Obj) (martini.Handler, error) {
		return func(r http.ResponseWriter) {
			r.Header().Add("X-Stack2", obj.RequiredString("header"))
		}, nil
	}))

	Register("errors", NewDefinition("errors", func(_ jsonconfig.Obj) (martini.Handler, error) {
		return nil, errors.New("failed to create")
	}))
}

func Test_AddMiddleware(t *testing.T) {
	setup()
	s := NewStack()
	s.AddMiddleware("stack1", nil)

	mw := len(s.Middlewares)

	if mw != 1 {
		t.Errorf("Unexpected middleware count: %d", mw)
	}
}

func Test_AddToServer(t *testing.T) {
	setup()
	m := martini.Classic()
	m.Get("/", func(r http.ResponseWriter) {
		fmt.Fprintf(r, "Hello")
	})
	s := NewStack()
	s.AddMiddleware("stack1", nil)
	s.AddMiddleware("stack2", jsonconfig.Obj{
		"header": "foo",
	})

	err := s.AddToServer(m.Martini)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := testhelpers.MakeTestableRequest(m, req)
	res.CheckHeader("X-Stack1", "stack1", t)
	res.CheckHeader("X-Stack2", "foo", t)
}

func checkForError(message, expected string, t *testing.T) {
	if !strings.Contains(message, expected) {
		t.Errorf("Expected %s to contain error %s", message, expected)
	}
}

func Test_MiddlewareErrors(t *testing.T) {
	setup()
	s := NewStack()
	s.AddMiddleware("missing1", nil)
	s.AddMiddleware("missing2", nil)
	s.AddMiddleware("errors", nil)
	err := s.reify()

	if err == nil {
		t.Fatal("No error was thrown")
	}

	msg := err.Error()

	checkForError(msg, "Unknown middleware: missing1", t)
	checkForError(msg, "Unknown middleware: missing2", t)
	checkForError(msg, "failed to create", t)
}
