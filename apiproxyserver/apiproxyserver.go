package apiproxyserver

import (
	"fourth.com/ratelimit/limiter"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

func NewProxyServer(r *limiter.RateLimit) http.Handler {
	m := createMartini(r)
	url, _ := url.Parse("http://localhost:3000")
	proxy := httputil.NewSingleHostReverseProxy(url)
	m.Action(proxy.ServeHTTP)
	m.Use(rateLimitHandler)

	return m
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
