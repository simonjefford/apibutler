package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"fourth.com/ratelimit/limiter"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"strconv"
)

type options struct {
	proxyPort     int
	dashboardPort int
	publicPath    string
}

var (
	opts options
)

func (o options) proxyPortString() string {
	return fmt.Sprintf(":%d", o.proxyPort)
}

func (o options) dashboardPortString() string {
	return fmt.Sprintf(":%d", o.dashboardPort)
}

func init() {
	flag.IntVar(&opts.proxyPort, "proxyPort", 4000, "Port on which to run the rate limiting proxy")
	flag.IntVar(&opts.dashboardPort, "dashboardPort", 8080, "Port on which to run the dashboard webapp")
	flag.StringVar(&opts.publicPath, "publicPath", "public", "Folder containing the webapp static assets")
	flag.Parse()
}

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

	log.Println("Running proxy on", opts.proxyPortString())
	log.Fatalln(http.ListenAndServe(opts.proxyPortString(), martini))
}

type StatusResponse struct {
	Message string `json:message`
}

func newDashboardServer() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(martini.Static(opts.publicPath))
	m.Action(r.Handle)

	return &martini.ClassicMartini{m, r}
}

func startDashboardServer(r *limiter.RateLimit) {
	srv := newDashboardServer()

	srv.Use(render.Renderer())

	srv.Post("/paths", func(res http.ResponseWriter, req *http.Request, rdr render.Render) {
		decoder := json.NewDecoder(req.Body)
		var p limiter.Path
		err := decoder.Decode(&p)
		if err != nil {
			rdr.JSON(http.StatusBadRequest, StatusResponse{err.Error()})
			return
		}
		log.Println(p)
		r.AddPath(p)
		rdr.JSON(http.StatusCreated, StatusResponse{"Created"})
	})

	srv.Get("/paths", func(rdr render.Render) {
		rdr.JSON(200, r.Paths())
	})

	log.Println("Running dashboard on", opts.dashboardPortString())
	log.Fatalln(http.ListenAndServe(opts.dashboardPortString(), srv))
}

func main() {
	// use all available cores
	runtime.GOMAXPROCS(runtime.NumCPU())
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
