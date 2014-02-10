package limiter

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"log"
	"sync"
	"time"
)

var (
	PathNotKnown           = errors.New("Path not known")
	RateLimitExceededError = errors.New("Rate limit exceeded")
)

type RateLimit struct {
	rw    sync.RWMutex
	calls map[string]*CallInfo
	rdb   redis.Conn
}

type Path struct {
	Fragment string `json:"fragment"`
	Limit    int    `json:"limit"`
	Seconds  int    `json:"seconds"`
}

func (r *RateLimit) AddPath(p Path) {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.calls[p.Fragment] = NewCallInfo(p.Limit, p.Seconds)
}

func (r *RateLimit) Paths() []Path {
	r.rw.RLock()
	defer r.rw.RUnlock()

	ps := make([]Path, 0, len(r.calls))
	for path, c := range r.calls {
		ps = append(ps, Path{
			Fragment: path,
			Limit:    c.Limit,
			Seconds:  int(c.Seconds / time.Second),
		})
	}

	return ps
}

func (r *RateLimit) IncrementCount(path string) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	call := r.calls[path]

	if call == nil {
		return nil
	}

	call.ResetIfNeccesary()

	if call.IsLimitExceeded() {
		return RateLimitExceededError
	}

	call.Count++
	return nil
}

func (r *RateLimit) Forget(path string) {
	log.Println("Now forgetting")
	r.rw.Lock()
	defer r.rw.Unlock()
	delete(r.calls, path)
}

func (r *RateLimit) GetCount(path string) (int, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	if r.calls[path] == nil {
		return 0, PathNotKnown
	}

	return r.calls[path].Count, nil
}

func (r *RateLimit) GetRemaining(path string) (int, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	if r.calls[path] == nil {
		return 0, PathNotKnown
	}

	return r.calls[path].Remaining(), nil
}

func NewRateLimit() (*RateLimit, error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}
	return &RateLimit{
		calls: make(map[string]*CallInfo),
		rdb:   conn,
	}, nil
}
