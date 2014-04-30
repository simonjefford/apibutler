package middleware

import (
	"net/http"
	"testing"

	"fmt"

	"fourth.com/apibutler/jsonconfig"
	"fourth.com/apibutler/testhelpers"
	"github.com/codegangsta/martini"
)

func init() {
	Register("stack1", func(obj jsonconfig.Obj) (martini.Handler, error) {
		return func(r http.ResponseWriter) {
			r.Header().Add("x-stack1", "stack1")
		}, nil
	})

	Register("stack2", func(obj jsonconfig.Obj) (martini.Handler, error) {
		return func(r http.ResponseWriter) {
			r.Header().Add("x-stack2", obj.RequiredString("header"))
		}, nil
	})
}

func Test_AddMiddleware(t *testing.T) {
	s := NewStack()
	s.AddMiddleware("stack1", nil)

	mw := len(s.middlewares)
	cfg := len(s.configs)

	if mw != 1 {
		t.Errorf("Unexpected middleware count: %d", mw)
	}

	if cfg != 1 {
		t.Errorf("Unexpected config count: %d", cfg)
	}
}

func Test_AddToServer(t *testing.T) {
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
