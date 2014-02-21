package limiter

import (
	"encoding/json"
	"errors"
	"fmt"
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

func redisConfigKeyForPath(p string) string {
	return fmt.Sprintf("%s:config", p)
}

func (r *RateLimit) AddPath(p Path) {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.calls[p.Fragment] = NewCallInfo(p.Limit, p.Seconds)

	r.rdb.Do("RPUSH", "knownPaths", p.Fragment)

	enc, _ := json.Marshal(p)

	err, ret := r.rdb.Do("SET", redisConfigKeyForPath(p.Fragment), string(enc))
	fmt.Println(err, ret)
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

func (r *RateLimit) loadPaths() error {
	n, err := redis.Int(r.rdb.Do("LLEN", "knownPaths"))
	if err != nil {
		return err
	}

	log.Println(n, "known paths")

	if n == 0 {
		return nil
	}

	paths, err := redis.Strings(r.rdb.Do("LRANGE", "knownPaths", 0, n))

	if err != nil {
		return err
	}

	for idx := range paths {
		err = r.rdb.Send("GET", redisConfigKeyForPath(paths[idx]))
		if err != nil {
			return err
		}
	}

	r.rdb.Flush()

	for _ = range paths {
		config, err := redis.String(r.rdb.Receive())
		if err != nil {
			return err
		}
		var p Path
		err = json.Unmarshal([]byte(config), &p)
		if err != nil {
			return err
		}
		r.calls[p.Fragment] = NewCallInfo(p.Limit, p.Seconds)
	}

	return nil
}

func NewRateLimit() (*RateLimit, error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}

	r := &RateLimit{
		calls: make(map[string]*CallInfo),
		rdb:   conn,
	}

	err = r.loadPaths()

	if err != nil {
		return nil, err
	}

	return r, nil
}
