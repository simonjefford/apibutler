// Package apiproxyserver provides the core transparent HTTP routing and management
// proxy capability.
package apiproxyserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"fourth.com/apibutler/metadata"
	"fourth.com/apibutler/middleware"
	_ "fourth.com/apibutler/middleware/oauth"

	"github.com/codegangsta/martini"
	"github.com/nickstenning/router/triemux"
)

type proxyserver struct {
	apps    metadata.ApplicationTable
	apis    []*metadata.Api
	logger  *log.Logger
	handler http.Handler
	sync.RWMutex
}

// APIProxyServer represents an HTTP server that acts as a transparent
// routing proxy server.
type APIProxyServer interface {
	// Update updates the routing tables used by the APIProxyServer
	UpdateApis(apis []*metadata.Api)

	// UpdateApps updates the list of backend applications used by
	// the APIProxyServer
	UpdateApps(apps metadata.ApplicationTable)

	// ServeHTTP is the method needed to implement http.Handler
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// NewAPIProxyServer returns a type implementing APIProxyServer. This type
// implements http.Handler so the return from this function can be given to
// http.ListenAndServe and friends.
func NewAPIProxyServer(apps metadata.ApplicationTable, apis []*metadata.Api) APIProxyServer {
	s := &proxyserver{
		apps:   apps,
		apis:   apis,
		logger: log.New(os.Stdout, "[proxy server] ", 0),
	}

	s.configure()

	return s
}

func wrapApp(app http.Handler, api *metadata.Api) http.Handler {
	m := martini.New()
	m.Action(app.ServeHTTP)
	l := log.New(os.Stdout, fmt.Sprintf("[%s (%s)] ", api.Fragment, api.App), 0)
	m.Map(l)
	if api.NeedsAuth {
		a, _ := middleware.Create("auth", nil)
		m.Use(a)
	}
	return m
}

func (s *proxyserver) ServeHTTP(res http.ResponseWriter, r *http.Request) {
	s.RLock()
	defer s.RUnlock()
	s.handler.ServeHTTP(res, r)
}

func (s *proxyserver) UpdateApis(apis []*metadata.Api) {
	s.Lock()
	defer s.Unlock()
	s.apis = apis
	s.configure()
}

func (s *proxyserver) UpdateApps(apps metadata.ApplicationTable) {
	s.Lock()
	defer s.Unlock()
	s.apps = apps
	s.configure()
}

func (s *proxyserver) configure() {
	mux := triemux.NewMux()

	for _, api := range s.apis {
		app, ok := s.apps[api.App]
		if ok {
			log.Printf("Handling %s with %v", api.Fragment, app)
			wrapped := wrapApp(app, api)
			mux.Handle(api.Fragment, api.IsPrefix, wrapped)
		} else {
			log.Printf("app not found")
		}
	}

	s.handler = createHost(s.logger, mux)
}

func createHost(l *log.Logger, mux *triemux.Mux) http.Handler {
	m := martini.New()
	m.Map(l)
	m.Action(mux.ServeHTTP)
	m.Use(martini.Logger())
	return m
}
