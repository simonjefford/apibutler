package metadata

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

type redisApiStore struct {
	rdb redis.Conn
}

func redisConfigKeyForApi(p string) string {
	return fmt.Sprintf("%s:config", p)
}

func (r *redisApiStore) AddApi(a *Api) {
	ret, err := r.rdb.Do("RPUSH", "knownPaths", a.Path)
	a.ID = ret.(int64)

	enc, _ := json.Marshal(a)

	ret, err = r.rdb.Do("SET", redisConfigKeyForApi(a.Path), string(enc))
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

func (r *redisApiStore) Apis() ([]*Api, error) {
	n, err := redis.Int(r.rdb.Do("LLEN", "knownPaths"))

	if err != nil {
		return nil, err
	}

	log.Println(n, "known paths")

	retApis := make([]*Api, 0, n)

	if n == 0 {
		return retApis, nil
	}

	paths, _ := redis.Strings(r.rdb.Do("LRANGE", "knownPaths", 0, n))

	if err != nil {
		return nil, err
	}

	for idx := range paths {
		r.rdb.Send("GET", redisConfigKeyForApi(paths[idx]))
		if err != nil {
			return nil, err
		}
	}

	r.rdb.Flush()

	for _ = range paths {
		config, _ := redis.String(r.rdb.Receive())
		if err != nil {
			return nil, err
		}
		var a Api
		json.Unmarshal([]byte(config), &a)
		retApis = append(retApis, &a)
	}

	return retApis, nil
}
