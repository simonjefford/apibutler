package middleware

import (
	"errors"
	"sync"

	"fourth.com/apibutler/jsonconfig"

	"github.com/codegangsta/martini"
)

type MiddlewareConstructor func(jsonconfig.Obj) (martini.Handler, error)

type MiddlewareConfigItem struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type MiddlewareDefinition struct {
	Schema       []*MiddlewareConfigItem `json:"schema"`
	Constructor  MiddlewareConstructor   `json:"-"`
	FriendlyName string                  `json:"friendly-name"`
}

func NewMiddlewareDefinition(
	friendlyName string,
	ctor MiddlewareConstructor,
	schema ...*MiddlewareConfigItem) *MiddlewareDefinition {

	return &MiddlewareDefinition{
		Schema:       schema,
		Constructor:  ctor,
		FriendlyName: friendlyName,
	}
}

type MiddlewareTable map[string]*MiddlewareDefinition

var (
	mu  sync.Mutex
	mws = make(MiddlewareTable)
)

func Register(name string, def *MiddlewareDefinition) error {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := mws[name]; dup {
		return errors.New("Duplicate registration of middleware constructor " + name)
	}

	mws[name] = def

	return nil
}

func GetMiddlewares() MiddlewareTable {
	mu.Lock()
	defer mu.Unlock()

	return mws
}

// TODO - should be able to return singletons if so configured
func Create(name string, cfg jsonconfig.Obj) (martini.Handler, error) {
	mu.Lock()
	defer mu.Unlock()
	def, ok := mws[name]
	if !ok {
		return nil, errors.New("Unknown middleware: " + name)
	}

	return def.Constructor(cfg)
}
