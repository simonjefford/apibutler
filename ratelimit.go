package main

import (
	"log"
	"sync"
)

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

func (r *RateLimit) Forget(path string) {
	log.Println("Now forgetting")
	r.rw.Lock()
	defer r.rw.Unlock()
	delete(r.calls, path)
}

func (r *RateLimit) GetCount(path string) int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.calls[path].Count
}

func (r *RateLimit) GetRemaining(path string) int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.calls[path].Remaining()
}

func NewRateLimit() *RateLimit {
	return &RateLimit{
		calls: make(map[string]*CallInfo),
	}
}
