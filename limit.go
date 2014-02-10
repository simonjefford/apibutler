package main

import (
	"encoding/json"
	"fourth.com/ratelimit/limiter"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
)

var (
	rateLimiter = limiter.NewRateLimit()
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
	err := rateLimiter.IncrementCount(path)
	if err == limiter.RateLimitExceededError {
		http.Error(res, err.Error(), http.StatusForbidden)
		return
	}

	rw := res.(martini.ResponseWriter)
	rw.Before(func(martini.ResponseWriter) {
		h := rw.Header()
		count, err := rateLimiter.GetCount(path)
		if err == nil {
			h.Add("X-Call-Count", strconv.Itoa(count))
		}

		remaining, err := rateLimiter.GetRemaining(path)
		if err == nil {
			h.Add("X-Call-Remaining", strconv.Itoa(remaining))
		}
	})

	ctx.Next()
}

func startLimitServer() {
	martini := createMartini()
	url, _ := url.Parse("http://localhost:3000")
	proxy := httputil.NewSingleHostReverseProxy(url)
	martini.Action(proxy.ServeHTTP)
	martini.Use(rateLimitHandler)

	martini.Run()
}

type StatusResponse struct {
	Message string `json:message`
}

func startDashboardServer() {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Post("/paths", func(res http.ResponseWriter, req *http.Request, r render.Render) {
		decoder := json.NewDecoder(req.Body)
		var p limiter.Path
		err := decoder.Decode(&p)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			log.Println(p)
			rateLimiter.AddPath(p)
			r.JSON(http.StatusCreated, StatusResponse{"Created"})
		}
	})

	m.Get("/paths", func(r render.Render) {
		r.JSON(200, rateLimiter.Paths())
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
