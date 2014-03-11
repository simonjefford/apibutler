package apiproxyserver

import (
	"fmt"

	"fourth.com/ratelimit/applications"
	"fourth.com/ratelimit/routes"

	"log"
	"net/http"
	"net/http/httptest"
	"os"
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
	res, err := makeRequest("GET", "/endpoint1", "Bearer some.bearer.token", srv)

	if err != nil {
		t.Fatalf(err.Error())
	}

	body := res.Body.String()

	checkResponse(http.StatusOK, res, t)

	if body != "endpoint1" {
		t.Fatalf("Unexpected response, \"%s\". Was the wrong endpoint hit?", body)
	}
}

func TestUnknownEndpoint(t *testing.T) {
	srv := configureProxyServer()
	res, err := makeRequest("GET", "/not.present", "Bearer some.bearer.token", srv)

	if err != nil {
		t.Fatalf(err.Error())
	}

	checkResponse(http.StatusNotFound, res, t)
}

func TestNoAuthHeader(t *testing.T) {
	srv := configureProxyServer()
	res, err := makeRequest("GET", "/not.present", "", srv)

	if err != nil {
		t.Fatalf(err.Error())
	}

	checkResponse(http.StatusUnauthorized, res, t)
}

func BenchmarkProxyRequests(b *testing.B) {
	srv := configureProxyServer()
	for i := 0; i < b.N; i++ {
		makeRequest("GET", "/endpoint1", "Bearer some.bearer.token", srv)
	}
}

func checkResponse(expected int, res *httptest.ResponseRecorder, t *testing.T) {
	if expected != res.Code {
		t.Fatalf("Response was not %d, was %d", expected, res.Code)
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
		logger: log.New(os.Stderr, "[TESTS] ", 0),
	}

	s.configure()
	return s
}

func makeRequest(method, path, auth string, handler http.Handler) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return nil, err
	}
	if auth != "" {
		req.Header.Add("Authorization", auth)
	}
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	return res, nil
}
