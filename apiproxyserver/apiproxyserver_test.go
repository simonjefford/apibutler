package apiproxyserver

import (
	"fmt"
	"fourth.com/ratelimit/applications"
	"fourth.com/ratelimit/routes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testendpoint struct {
	outputString string
}

func (t *testendpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, t.outputString)
}

func TestEndpointRouting(t *testing.T) {
	srv := configureProxyServer()
	res, err := makeRequest("GET", "/endpoint1", srv)

	if err != nil {
		t.Fatalf(err.Error())
	}

	body := res.Body.String()

	if res.Code != http.StatusOK {
		t.Fatalf("Response was not %d, was %d", http.StatusOK, res.Code)
	}

	if body != "endpoint1" {
		t.Fatalf("Unexpected response, \"%s\". Was the wrong endpoint hit?", body)
	}
}

func TestUnknownEndpoint(t *testing.T) {
	srv := configureProxyServer()
	res, err := makeRequest("GET", "/not.present", srv)

	if err != nil {
		t.Fatalf(err.Error())
	}

	if res.Code != http.StatusNotFound {
		t.Fatalf("Response was not %d, was %d", http.StatusNotFound, res.Code)
	}
}

func configureProxyServer() *proxyserver {
	m := applications.ApplicationTable(make(map[string]http.Handler))
	m["endpoint1"] = &testendpoint{"endpoint1"}
	r := []routes.Route{
		routes.Route{
			Path:            "/endpoint1",
			ApplicationName: "endpoint1",
			IsPrefix:        true,
		},
	}
	s := &proxyserver{
		apps:   m,
		routes: r,
	}

	s.configure(nil)
	return s
}

func makeRequest(method, path string, handler http.Handler) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer some.bearer.token")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	return res, nil
}
