package apiproxyserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"fourth.com/apibutler/metadata"
	"fourth.com/apibutler/testhelpers"
)

type testendpoint struct {
	outputString string
}

func (t *testendpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, t.outputString)
}

func TestEndpointRouting(t *testing.T) {
	srv := configureProxyServer()
	res, err := makeRequest("GET", "/endpoint1", "Bearer some.bearer.token", srv)

	if err != nil {
		t.Fatalf(err.Error())
	}

	res.CheckStatus(http.StatusOK, t)
	res.CheckBody("endpoint1", t)
}

func TestUnknownEndpoint(t *testing.T) {
	srv := configureProxyServer()
	res, err := makeRequest("GET", "/not.present", "Bearer some.bearer.token", srv)

	if err != nil {
		t.Fatalf(err.Error())
	}

	res.CheckStatus(http.StatusNotFound, t)
}

func TestNoAuthHeader(t *testing.T) {
	srv := configureProxyServer()
	res, err := makeRequest("GET", "/endpoint1", "", srv)

	if err != nil {
		t.Fatalf(err.Error())
	}

	res.CheckStatus(http.StatusUnauthorized, t)
}

func TestNoAuthHeaderOnPublicRoute(t *testing.T) {
	srv := configureProxyServer()
	res, err := makeRequest("GET", "/public", "", srv)

	if err != nil {
		t.Fatalf(err.Error())
	}

	res.CheckStatus(http.StatusOK, t)
}

func BenchmarkProxyRequests(b *testing.B) {
	srv := configureProxyServer()
	for i := 0; i < b.N; i++ {
		makeRequest("GET", "/endpoint1", "Bearer some.bearer.token", srv)
	}
}

func TestUpdateServer(t *testing.T) {
	m := make(metadata.ApplicationTable)
	m["endpoint1"] = &metadata.Application{
		Handler: &testendpoint{"endpoint1"},
		Name:    "endpoint1",
	}
	r := []*metadata.Api{
		&metadata.Api{Path: "/endpoint1", App: "endpoint1"},
		&metadata.Api{Path: "/newendpoint", App: "endpoint1"},
	}

	srv := configureProxyServer()
	srv.UpdateApis(r)
	srv.UpdateApps(m)
	res, err := makeRequest("GET", "/newendpoint", "Bearer some.bearer.token", srv)

	if err != nil {
		t.Fatalf(err.Error())
	}

	res.CheckStatus(http.StatusOK, t)
}

func configureProxyServer() APIProxyServer {
	m := make(metadata.ApplicationTable)
	m["endpoint1"] = &metadata.Application{
		Handler: &testendpoint{"endpoint1"},
		Name:    "endpoint1",
	}
	m["public"] = &metadata.Application{
		Handler: &testendpoint{"public"},
		Name:    "public",
	}
	r := []*metadata.Api{
		&metadata.Api{
			Path:      "/endpoint1",
			App:       "endpoint1",
			NeedsAuth: true,
		},
		&metadata.Api{
			Path:      "/public",
			App:       "public",
			NeedsAuth: false,
		},
	}
	s := &proxyserver{
		apps:   m,
		apis:   r,
		logger: log.New(os.Stderr, "[TESTS] ", 0),
	}

	s.configure()
	return s
}

func makeRequest(method, path, auth string, handler http.Handler) (*testhelpers.CheckableResponse, error) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return nil, err
	}
	if auth != "" {
		req.Header.Add("Authorization", auth)
	}
	return testhelpers.MakeTestableRequest(handler, req), nil
}
