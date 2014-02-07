package main

import (
	"errors"
	"github.com/codegangsta/martini"
	"strconv"

	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	limiter                = NewRateLimit()
	RateLimitExceededError = errors.New("Rate limit exceeded")
)

func createMartini() *martini.Martini {
	m := martini.New()
	m.Use(martini.Logger())
	return m
}

func rateLimitHandler(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	err := limiter.IncrementCount(path)
	if err == RateLimitExceededError {
		http.Error(res, err.Error(), http.StatusForbidden)
		return
	}
	h := res.Header()
	h.Add("X-Call-Count", strconv.Itoa(limiter.GetCount(path)))
	h.Add("X-Endpoint", path)
}

func main() {
	martini := createMartini()
	url, _ := url.Parse("http://localhost:3000")
	proxy := httputil.NewSingleHostReverseProxy(url)
	martini.Action(proxy.ServeHTTP)
	martini.Use(rateLimitHandler)
	martini.Run()
}
