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

func NewApplication(name, urlString string, id int) *Application {
	parsed, _ := url.Parse(urlString)
	return &Application{
		Name:             name,
		BackendURLString: urlString,
		backendURL:       parsed,
		Handler:          httputil.NewSingleHostReverseProxy(parsed),
		ID:               id,
	}
}

var (
	apps ApplicationTable
)

type ApplicationTable map[string]*Application

func init() {
	apps = make(ApplicationTable)
	apps["Test node backend"] = NewApplication("Test node backend", "http://localhost:3000", 1)
	apps["Another test node backend"] = NewApplication("Another test node backend", "http://localhost:3001", 2)
}

func GetApplicationsTable() ApplicationTable {
	return apps
}

func GetSingleApplication(id int) *Application {
	switch id {
	case 1:
		return apps["Test node backend"]
	case 2:
		return apps["Another test node backend"]
	}

	return nil
}

func GetApplicationsList() []*Application {
	t := make([]*Application, 0, len(apps))
	for _, val := range apps {
		t = append(t, val)
	}
	return t
}
