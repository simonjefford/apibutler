package metadata

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Application struct {
	Name             string `json:"name"`
	backendURL       *url.URL
	BackendURLString string `json:"backendURL"`
	http.Handler     `json:"-"`
	ID               int `json:"id"`
}

func NewApplication(name, urlString string) *Application {
	parsed, _ := url.Parse(urlString)
	return &Application{
		Name:             name,
		BackendURLString: urlString,
		backendURL:       parsed,
		Handler:          httputil.NewSingleHostReverseProxy(parsed),
	}
}

var (
	apps ApplicationTable
)

type ApplicationTable map[string]*Application

func init() {
	apps = make(ApplicationTable)
	apps["Test node backend"] = NewApplication("Test node backend", "http://localhost:3000")
	apps["Another test node backend"] = NewApplication("Another test node backend", "http://localhost:3001")
}

func GetApplicationsTable() ApplicationTable {
	return apps
}

func GetApplicationsList() []*Application {
	t := make([]*Application, 0, len(apps))
	i := 0
	for _, val := range apps {
		i = i + 1
		val.ID = i
		t = append(t, val)
	}
	return t
}
