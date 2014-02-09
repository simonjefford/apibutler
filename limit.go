package main

import (
	"encoding/json"
	"errors"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
)

var (
	limiter                = NewRateLimit()
	RateLimitExceededError = errors.New("Rate limit exceeded")
)

func createMartini() *martini.Martini {
	m := martini.New()
	m.Use(martini.Logger())
	l := log.New(os.Stdout, "[martini ratelimiter] ", 0)
	m.Map(l)
	return m
}

func rateLimitHandler(res http.ResponseWriter, req *http.Request, ctx martini.Context) {
	path := req.URL.Path
	err := limiter.IncrementCount(path)
	if err == RateLimitExceededError {
		http.Error(res, err.Error(), http.StatusForbidden)
		return
	}

	rw := res.(martini.ResponseWriter)
	rw.Before(func(martini.ResponseWriter) {
		h := rw.Header()
		count, err := limiter.GetCount(path)
		if err == nil {
			h.Add("X-Call-Count", strconv.Itoa(count))
		}

		remaining, err := limiter.GetRemaining(path)
		if err == nil {
			h.Add("X-Call-Remaining", strconv.Itoa(remaining))
		}
	})

	ctx.Next()
}

func statusCodeIsSuccessful(status int) bool {
	return status >= 200 && status <= 299
}

func startLimitServer() {
	martini := createMartini()
	url, _ := url.Parse("http://localhost:3000")
	proxy := httputil.NewSingleHostReverseProxy(url)
	martini.Action(proxy.ServeHTTP)
	martini.Use(rateLimitHandler)

	martini.Run()
}

type Path struct {
	Fragment string
}

func startDashboardServer() {
	m := martini.Classic()

	m.Post("/paths", func(res http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var p Path
		err := decoder.Decode(&p)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			log.Println(p)
			limiter.AddPath(p.Fragment)
			res.WriteHeader(http.StatusCreated)
		}

	})

	log.Fatalln(http.ListenAndServe(":8080", m))
}

func main() {
	go startLimitServer()
	go startDashboardServer()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
