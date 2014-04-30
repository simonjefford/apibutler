package middleware

import (
	"errors"
	"log"
	"sync"

	"fourth.com/apibutler/jsonconfig"

	"github.com/codegangsta/martini"
)

type MiddlewareConstructor func(jsonconfig.Obj) (martini.Handler, error)

var (
	mu    sync.Mutex
	ctors = make(map[string]MiddlewareConstructor)
)

func Register(name string, fn MiddlewareConstructor) error {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := ctors[name]; dup {
		return errors.New("Duplicate registration of middleware constructor " + name)
	}

	ctors[name] = fn
	log.Println(ctors)
	return nil
}

// TODO - should be able to return singletons if so configured
func Create(name string, cfg jsonconfig.Obj) (martini.Handler, error) {
	mu.Lock()
	defer mu.Unlock()
	fn, ok := ctors[name]
	if !ok {
		return nil, errors.New("Unknown middleware: " + name)
	}

	return fn(cfg)
}
