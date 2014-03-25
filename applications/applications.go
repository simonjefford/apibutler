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
	apps = make(ApplicationTable)
	apps["Test node backend"] = newApp("Test node backend", "http://localhost:3000")
	apps["Another test node backend"] = newApp("Another test node backend", "http://localhost:3001")
}

func Get() ApplicationTable {
	return apps
}

func GetList() []*application {
	t := make([]*application, 0, len(apps))
	for _, val := range apps {
		t = append(t, val.(*application))
	}
	return t
}
