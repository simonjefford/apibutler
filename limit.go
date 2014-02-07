package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"strconv"
	"sync"
	"time"

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

type CallInfo struct {
	Count int
	Start time.Time
}

func (c *CallInfo) timeSinceLastReset() time.Duration {
	return time.Since(c.Start)
}

func (c *CallInfo) LimitExceeded() bool {
	dur := c.timeSinceLastReset()
	return dur < time.Second*20 && c.Count >= 5
}

func (c *CallInfo) ResetIfNeccesary() {
	dur := c.timeSinceLastReset()
	log.Println(dur)
	if dur > time.Second*20 {
		log.Println("Resetting")
		c.Start = time.Now()
		c.Count = 0
	} else {
		log.Println("Not resetting")
	}
}

func NewCallInfo() *CallInfo {
	return &CallInfo{Start: time.Now()}
}

type RateLimit struct {
	rw    sync.RWMutex
	calls map[string]*CallInfo
}

func (r *RateLimit) IncrementCount(path string) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	call := r.calls[path]
	if r.calls[path] == nil {
		call = NewCallInfo()
		r.calls[path] = call
	}

	call.ResetIfNeccesary()

	if call.LimitExceeded() {
		return RateLimitExceededError
	}

	call.Count++
	return nil
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
