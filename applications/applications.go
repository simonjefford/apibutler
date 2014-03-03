package applications

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	apps ApplicationTable
)

type application struct {
	Name             string `json:"name"`
	backendURL       *url.URL
	BackendURLString string `json:"backendURL"`
	http.Handler     `json:"-"`
}

type ApplicationTable map[string]http.Handler

func newApp(name, urlString string) *application {
	parsed, _ := url.Parse(urlString)
	return &application{
		Name:             name,
		BackendURLString: urlString,
		backendURL:       parsed,
		Handler:          httputil.NewSingleHostReverseProxy(parsed),
	}
}

func init() {
	m := make(map[string]http.Handler)
	m["Test node backend"] = newApp("Test node backend", "http://localhost:3000")
	m["Another test node backend"] = newApp("Another test node backend", "http://localhost:3001")

	apps = ApplicationTable(m)
}

func Get() ApplicationTable {
	return apps
}
