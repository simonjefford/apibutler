package ratelimiter

import (
	"errors"
	"net/http"
	"strconv"
	"sync"

	"fourth.com/apibutler/jsonconfig"
	"fourth.com/apibutler/middleware"

	"github.com/codegangsta/martini"
)

var (
	PathNotKnown           = errors.New("Path not known")
	RateLimitExceededError = errors.New("Rate limit exceeded")
)

func init() {
	middleware.Register(
		"ratelimiter",
		middleware.NewMiddlewareDefinition("Rate Limiting", ratelimiterCtor,
			&middleware.MiddlewareConfigItem{
				Name: "limit",
				Type: "integer",
			},
		))
}

type RateLimit interface {
	IncrementCount() error
	GetCount() int
	GetRemaining() int
	Handler(http.ResponseWriter, *http.Request, martini.Context)
}

type rateLimit struct {
	sync.RWMutex
	call *CallInfo
}

func (r *rateLimit) IncrementCount() error {
	r.Lock()
	defer r.Unlock()

	r.call.ResetIfNeccesary()

	if r.call.IsLimitExceeded() {
		return RateLimitExceededError
	}

	r.call.Count++

	return nil
}

func (r *rateLimit) GetCount() int {
	r.RLock()
	defer r.RUnlock()

	return r.call.Count
}

func (r *rateLimit) GetRemaining() int {
	r.RLock()
	defer r.RUnlock()

	return r.call.Remaining()
}

func (r *rateLimit) Handler(res http.ResponseWriter, req *http.Request, ctx martini.Context) {
	err := r.IncrementCount()
	if err == RateLimitExceededError {
		http.Error(res, err.Error(), http.StatusForbidden)
		return
	}
	rw := res.(martini.ResponseWriter)
	rw.Before(func(martini.ResponseWriter) {
		h := rw.Header()
		count := r.GetCount()
		if err == nil {
			h.Add("X-Call-Count", strconv.Itoa(count))
		}

		remaining := r.GetRemaining()
		if err == nil {
			h.Add("X-Call-Remaining", strconv.Itoa(remaining))
		}
	})

	ctx.Next()
}

func NewRateLimit(limit, seconds int) RateLimit {
	return &rateLimit{
		call: NewCallInfo(limit, seconds),
	}
}

func ratelimiterCtor(cfg jsonconfig.Obj) (martini.Handler, error) {
	l := cfg.RequiredInt("limit")
	s := cfg.RequiredInt("seconds")
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	return NewRateLimit(l, s).Handler, nil
}
