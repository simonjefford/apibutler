package applications

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Application struct {
	Name       string
	BackendURL *url.URL
	handler    http.Handler
}

type ApplicationTable map[string]*Application

func (a *Application) ProxyServer() http.Handler {
	if a.handler == nil {
		a.handler = httputil.NewSingleHostReverseProxy(a.BackendURL)
	}

	return a.handler
}

func Get() ApplicationTable {
	m := make(map[string]*Application)
	backend, _ := url.Parse("http://localhost:3000")
	m["Test node backend"] = &Application{
		Name:       "Test node backend",
		BackendURL: backend,
	}

	backend, _ = url.Parse("http://localhost:3001")
	m["Another test node backend"] = &Application{
		Name:       "Another test node backend",
		BackendURL: backend,
	}

	return ApplicationTable(m)
}
