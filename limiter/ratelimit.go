package limiter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	PathNotKnown           = errors.New("Path not known")
	RateLimitExceededError = errors.New("Rate limit exceeded")
)

type RateLimit interface {
	AddPath(p Path)
	Paths() []Path
	IncrementCount(path string) error
	Forget(path string)
	GetCount(path string) (int, error)
	GetRemaining(path string) (int, error)
}

type rateLimit struct {
	sync.RWMutex
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

func (r *rateLimit) AddPath(p Path) {
	r.Lock()
	defer r.Unlock()

	r.calls[p.Fragment] = NewCallInfo(p.Limit, p.Seconds)

	r.rdb.Do("RPUSH", "knownPaths", p.Fragment)

	enc, _ := json.Marshal(p)

	err, ret := r.rdb.Do("SET", redisConfigKeyForPath(p.Fragment), string(enc))
	fmt.Println(err, ret)
}

func (r *rateLimit) Paths() []Path {
	r.RLock()
	defer r.RUnlock()

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

func (r *rateLimit) IncrementCount(path string) error {
	r.Lock()
	defer r.Unlock()
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

func (r *rateLimit) Forget(path string) {
	log.Println("Now forgetting")
	r.Lock()
	defer r.Unlock()
	delete(r.calls, path)
}

func (r *rateLimit) GetCount(path string) (int, error) {
	r.RLock()
	defer r.RUnlock()

	if r.calls[path] == nil {
		return 0, PathNotKnown
	}

	return r.calls[path].Count, nil
}

func (r *rateLimit) GetRemaining(path string) (int, error) {
	r.RLock()
	defer r.RUnlock()

	if r.calls[path] == nil {
		return 0, PathNotKnown
	}

	return r.calls[path].Remaining(), nil
}

func (r *rateLimit) loadPaths() error {
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

func (r *rateLimit) Handler(res http.ResponseWriter, req *http.Request, ctx martini.Context) {
	path := req.URL.Path
	err := r.IncrementCount(path)
	if err == RateLimitExceededError {
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

func NewRateLimit() (RateLimit, error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}

	r := &rateLimit{
		calls: make(map[string]*CallInfo),
		rdb:   conn,
	}

	err = r.loadPaths()

	if err != nil {
		return nil, err
	}

	return r, nil
}
