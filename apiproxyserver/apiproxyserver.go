package apiproxyserver

import (
	"fmt"
	"fourth.com/ratelimit/applications"
	"fourth.com/ratelimit/oauth"
	"fourth.com/ratelimit/routes"
	"github.com/codegangsta/martini"
	"github.com/nickstenning/router/triemux"
	"log"
	"net/http"
	"os"
)

type proxyserver struct {
	apps   applications.ApplicationTable
	routes []routes.Route
	logger *log.Logger
	http.Handler
}

type destinationApp struct {
	original http.Handler
	*martini.Martini
}

func NewProxyServer() (http.Handler, error) {
	s := proxyserver{
		apps:   applications.Get(),
		routes: routes.Get(),
		logger: log.New(os.Stdout, "[proxy server] ", 0),
	}

	s.configure()

	return s, nil
}

func wrapApp(app http.Handler, route routes.Route) *destinationApp {
	m := martini.New()
	m.Action(app.ServeHTTP)
	l := log.New(os.Stdout, fmt.Sprintf("[%s (%s)] ", route.Path, route.ApplicationName), 0)
	m.Map(l)
	if route.NeedsAuth {
		m.Use(oauth.GetIdFromRequest)
		m.Use(logToken)
	}
	return &destinationApp{app, m}
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

	s.Handler = createHost(s.logger, mux)
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
