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
	"fourth.com/apibutler/oauth"
	"github.com/codegangsta/martini"
	"github.com/nickstenning/router/triemux"
)

type proxyserver struct {
	apps    metadata.ApplicationTable
	routes  []metadata.Route
	logger  *log.Logger
	handler http.Handler
	sync.RWMutex
}

// APIProxyServer represents an HTTP server that acts as a transparent
// routing proxy server.
type APIProxyServer interface {
	// Update updates the application and routing tables used by
	// the APIProxyServer
	Update(metadata.ApplicationTable, []metadata.Route)

	// ServeHTTP is the method needed to implement http.Handler
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// NewAPIProxyServer returns a type implementing APIProxyServer. This type
// implements http.Handler so the return from this function can be given to
// http.ListenAndServe and friends.
func NewAPIProxyServer() APIProxyServer {
	s := &proxyserver{
		apps:   metadata.GetApplicationsTable(),
		routes: metadata.GetRoutes(),
		logger: log.New(os.Stdout, "[proxy server] ", 0),
	}

	s.configure()

	return s
}

func wrapApp(app http.Handler, route metadata.Route) http.Handler {
	m := martini.New()
	m.Action(app.ServeHTTP)
	l := log.New(os.Stdout, fmt.Sprintf("[%s (%s)] ", route.Path, route.ApplicationName), 0)
	m.Map(l)
	if route.NeedsAuth {
		m.Use(oauth.GetIdFromRequest)
		m.Use(logToken)
	}
	return m
}

func (s *proxyserver) ServeHTTP(res http.ResponseWriter, r *http.Request) {
	s.RLock()
	defer s.RUnlock()
	s.handler.ServeHTTP(res, r)
}

func (s *proxyserver) Update(apps metadata.ApplicationTable, routes []metadata.Route) {
	s.Lock()
	defer s.Unlock()
	s.apps = apps
	s.routes = routes
	s.configure()
}

func (s *proxyserver) configure() {
	mux := triemux.NewMux()

	for _, route := range s.routes {
		app, ok := s.apps[route.ApplicationName]
		if ok {
			log.Printf("Handling %s with %v", route.Path, app)
			wrapped := wrapApp(app, route)
			mux.Handle(route.Path, route.IsPrefix, wrapped)
		} else {
			log.Printf("app not found")
		}
	}

	s.handler = createHost(s.logger, mux)
}

func logToken(t oauth.AccessToken, l *log.Logger) {
	l.Println(t)
}

func createHost(l *log.Logger, mux *triemux.Mux) http.Handler {
	m := martini.New()
	m.Map(l)
	m.Action(mux.ServeHTTP)
	m.Use(martini.Logger())
	return m
}
