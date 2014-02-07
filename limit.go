package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"strconv"
	"sync"
	"time"

	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	limiter = NewRateLimit()
)

func createMartini() *martini.Martini {
	m := martini.New()
	m.Use(martini.Logger())
	return m
}

type CallInfo struct {
	Count int
	Time  time.Time
	Start time.Time
}

type RateLimit struct {
	rw    sync.RWMutex
	calls map[string]*CallInfo
}

func (r *RateLimit) IncrementCount(path string) {
	r.rw.Lock()
	defer r.rw.Unlock()
	if r.calls[path] == nil {
		r.calls[path] = &CallInfo{Start: time.Now()}
	}
	r.calls[path].Count++
	r.calls[path].Time = time.Now()
}

func (r *RateLimit) GetCount(path string) int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	fmt.Println(r.calls[path])
	return r.calls[path].Count
}

func NewRateLimit() *RateLimit {
	return &RateLimit{
		calls: make(map[string]*CallInfo),
	}
}

func rateLimitHandler(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	limiter.IncrementCount(path)
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
