package apiproxyserver

import (
	"fourth.com/ratelimit/applications"
	"fourth.com/ratelimit/limiter"
	"fourth.com/ratelimit/oauth"
	"fourth.com/ratelimit/routes"
	"github.com/codegangsta/martini"
	"github.com/nickstenning/router/triemux"
	"log"
	"net/http"
	"os"
	"strconv"
)

type proxyserver struct {
	apps   applications.ApplicationTable
	routes []routes.Route
	http.Handler
}

func NewProxyServer(r *limiter.RateLimit) http.Handler {
	s := proxyserver{
		apps:   applications.Get(),
		routes: routes.Get(),
	}

	s.configure(r)

	return s
}

func (s *proxyserver) configure(r *limiter.RateLimit) {
	m := createMartini(r)

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

	m.Action(mux.ServeHTTP)
	m.Use(oauth.GetIdFromRequest)
	m.Use(logToken)
	s.Handler = m
}

func logToken(t oauth.AccessToken, l *log.Logger) {
	l.Println(t)
}

func createMartini(r *limiter.RateLimit) *martini.Martini {
	m := martini.New()
	m.Use(martini.Logger())
	l := log.New(os.Stdout, "[proxy server] ", 0)
	m.Map(l)
	m.Map(r)
	return m
}

func rateLimitHandler(res http.ResponseWriter, req *http.Request, ctx martini.Context, r *limiter.RateLimit) {
	path := req.URL.Path
	err := r.IncrementCount(path)
	if err == limiter.RateLimitExceededError {
		http.Error(res, err.Error(), http.StatusForbidden)
		return
	}

	rw := res.(martini.ResponseWriter)
	rw.Before(func(martini.ResponseWriter) {
		h := rw.Header()
		count, err := r.GetCount(path)
		if err == nil {
			h.Add("X-Call-Count", strconv.Itoa(count))
		}

		remaining, err := r.GetRemaining(path)
		if err == nil {
			h.Add("X-Call-Remaining", strconv.Itoa(remaining))
		}
	})

	ctx.Next()
}
