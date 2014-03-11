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
	apps    applications.ApplicationTable
	routes  []routes.Route
	logger  *log.Logger
	limiter limiter.RateLimit
	http.Handler
}

func NewProxyServer() (http.Handler, error) {
	l, err := limiter.NewRateLimit()

	if err != nil {
		return nil, err
	}

	s := proxyserver{
		apps:    applications.Get(),
		routes:  routes.Get(),
		logger:  log.New(os.Stdout, "[proxy server] ", 0),
		limiter: l,
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

	s.Handler = createHost(s.limiter, s.logger, mux)
}

func logToken(t oauth.AccessToken, l *log.Logger) {
	l.Println(t)
}

func createHost(r limiter.RateLimit, l *log.Logger, mux *triemux.Mux) http.Handler {
	m := martini.New()
	m.Map(l)
	m.MapTo(r, (*limiter.RateLimit)(nil))
	m.Action(mux.ServeHTTP)
	m.Use(martini.Logger())
	m.Use(oauth.GetIdFromRequest)
	m.Use(logToken)
	m.Use(rateLimitHandler)
	return m
}

func rateLimitHandler(res http.ResponseWriter, req *http.Request, ctx martini.Context, r limiter.RateLimit) {
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
