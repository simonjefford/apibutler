package apiproxyserver

import (
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

func NewProxyServer() (http.Handler, error) {
	s := proxyserver{
		apps:   applications.Get(),
		routes: routes.Get(),
		logger: log.New(os.Stdout, "[proxy server] ", 0),
	}

	s.configure()

	return s, nil
}

func (s *proxyserver) configure() {
	mux := triemux.NewMux()

	for _, route := range s.routes {
		app, ok := s.apps[route.ApplicationName]
		if ok {
			log.Printf("Handling %s with %v", route.Path, app)
			mux.Handle(route.Path, route.IsPrefix, app)
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
	m.Use(oauth.GetIdFromRequest)
	m.Use(logToken)
	return m
}
