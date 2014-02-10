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

func createMartini(r *limiter.RateLimit) *martini.Martini {
	m := martini.New()
	m.Use(martini.Logger())
	l := log.New(os.Stdout, "[martini ratelimiter] ", 0)
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

func startLimitServer(r *limiter.RateLimit) {
	martini := createMartini(r)
	url, _ := url.Parse("http://localhost:3000")
	proxy := httputil.NewSingleHostReverseProxy(url)
	martini.Action(proxy.ServeHTTP)
	martini.Use(rateLimitHandler)

	martini.Run()
}

type StatusResponse struct {
	Message string `json:message`
}

func startDashboardServer(r *limiter.RateLimit) {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Post("/paths", func(res http.ResponseWriter, req *http.Request, rdr render.Render) {
		decoder := json.NewDecoder(req.Body)
		var p limiter.Path
		err := decoder.Decode(&p)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
		} else {
			log.Println(p)
			r.AddPath(p)
			rdr.JSON(http.StatusCreated, StatusResponse{"Created"})
		}
	})

	m.Get("/paths", func(rdr render.Render) {
		rdr.JSON(200, r.Paths())
	})

	log.Fatalln(http.ListenAndServe(":8080", m))
}

func main() {
	r, err := limiter.NewRateLimit()
	if err != nil {
		log.Fatalln(err)
	}
	go startLimitServer(r)
	go startDashboardServer(r)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
