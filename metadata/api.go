package metadata

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

type Path struct {
	Fragment string `json:"fragment"`
	Limit    int    `json:"limit"`
	Seconds  int    `json:"seconds"`
	ID       int64  `json:"id"`
}

type PathStorage interface {
	AddPath(p Path)
	Paths() []Path
	Forget(path string)
}

type redisPathStore struct {
	rdb redis.Conn
}

func redisConfigKeyForPath(p string) string {
	return fmt.Sprintf("%s:config", p)
}

func (r *redisPathStore) AddPath(p Path) {
	ret, err := r.rdb.Do("RPUSH", "knownPaths", p.Fragment)
	p.ID = ret.(int64)

	enc, _ := json.Marshal(p)

	ret, err = r.rdb.Do("SET", redisConfigKeyForPath(p.Fragment), string(enc))
	fmt.Println(err, ret)
}

func (r *redisPathStore) Forget(path string) {
}

func GetPathStore() (PathStorage, error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}

	return &redisPathStore{conn}, nil
}

func (r *redisPathStore) Paths() []Path {
	n, _ := redis.Int(r.rdb.Do("LLEN", "knownPaths"))

	// TODO revisit me
	// if err != nil {
	// 	return err
	// }

	log.Println(n, "known paths")

	if n == 0 {
		return nil
	}

	retPaths := make([]Path, 0, n)

	paths, _ := redis.Strings(r.rdb.Do("LRANGE", "knownPaths", 0, n))

	// if err != nil {
	// 	return err
	// }

	for idx := range paths {
		r.rdb.Send("GET", redisConfigKeyForPath(paths[idx]))
		// if err != nil {
		// 	return err
		// }
	}

	r.rdb.Flush()

	for _ = range paths {
		config, _ := redis.String(r.rdb.Receive())
		// if err != nil {
		// 	return err
		// }
		var p Path
		json.Unmarshal([]byte(config), &p)
		retPaths = append(retPaths, p)
	}

	return retPaths
}
