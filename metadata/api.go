package metadata

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

type Api struct {
	Fragment string `json:"fragment"`
	Limit    int    `json:"limit"`
	Seconds  int    `json:"seconds"`
	ID       int64  `json:"id"`
}

type ApiStorage interface {
	AddApi(a *Api)
	Apis() []Api
	Forget(path string)
}

type redisApiStore struct {
	rdb redis.Conn
}

func redisConfigKeyForApi(p string) string {
	return fmt.Sprintf("%s:config", p)
}

func (r *redisApiStore) AddApi(a *Api) {
	ret, err := r.rdb.Do("RPUSH", "knownPaths", a.Fragment)
	a.ID = ret.(int64)

	enc, _ := json.Marshal(a)

	ret, err = r.rdb.Do("SET", redisConfigKeyForApi(a.Fragment), string(enc))
	fmt.Println(err, ret)
}

func (r *redisApiStore) Forget(path string) {
}

func GetApiStore() (ApiStorage, error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}

	return &redisApiStore{conn}, nil
}

func (r *redisApiStore) Apis() []Api {
	n, _ := redis.Int(r.rdb.Do("LLEN", "knownPaths"))

	// TODO revisit me
	// if err != nil {
	// 	return err
	// }

	log.Println(n, "known paths")

	if n == 0 {
		return nil
	}

	retApis := make([]Api, 0, n)

	paths, _ := redis.Strings(r.rdb.Do("LRANGE", "knownPaths", 0, n))

	// if err != nil {
	// 	return err
	// }

	for idx := range paths {
		r.rdb.Send("GET", redisConfigKeyForApi(paths[idx]))
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
		var p Api
		json.Unmarshal([]byte(config), &p)
		retApis = append(retApis, p)
	}

	return retApis
}
