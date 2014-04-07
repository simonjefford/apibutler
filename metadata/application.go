package metadata

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
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
	apps["1"] = NewApplication("Test node backend", "http://localhost:3000", 1)
	apps["2"] = NewApplication("Another test node backend", "http://localhost:3001", 2)
}

func GetApplicationsTable() ApplicationTable {
	return apps
}

func ChangeApplication(a *Application) {
	idstring := strconv.Itoa(a.ID)
	apps[idstring] = a
}

func GetSingleApplication(id int) *Application {
	switch id {
	case 1:
		return apps["1"]
	case 2:
		return apps["2"]
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
