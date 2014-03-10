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

func NewProxyServer(r *limiter.RateLimit) http.Handler {
	m := createMartini(r)
	apps := applications.Get()
	mux := triemux.NewMux()

	for _, route := range routes.Get() {
		app, ok := apps[route.ApplicationName]
		if ok {
			mux.Handle(route.Path, route.IsPrefix, app)
		}
	}

	m.Action(mux.ServeHTTP)
	m.Use(oauth.GetIdFromRequest)
	m.Use(logToken)
	m.Use(rateLimitHandler)
	return m
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
