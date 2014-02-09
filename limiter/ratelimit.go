package limiter

import (
	"errors"
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
}

type Path struct {
	Fragment string
	Limit    int
	Seconds  int
}

func (r *RateLimit) AddPath(path string, limit int, seconds int) {
	r.rw.Lock()
	defer r.rw.Unlock()

	if r.calls[path] == nil {
		r.calls[path] = NewCallInfo(limit, seconds)
	}
}

func (r *RateLimit) Paths() []Path {
	r.rw.RLock()
	defer r.rw.RUnlock()

	ps := make([]Path, len(r.calls))
	idx := 0
	for path, c := range r.calls {
		ps[idx] = Path{
			Fragment: path,
			Limit:    c.Limit,
			Seconds:  int(c.Seconds / time.Second),
		}
		idx++
	}

	return ps
}

func (r *RateLimit) IncrementCount(path string) error {
	r.rw.Lock()
	defer r.rw.Unlock()
	call := r.calls[path]

	if r.calls[path] == nil {
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

func NewRateLimit() *RateLimit {
	return &RateLimit{
		calls: make(map[string]*CallInfo),
	}
}
