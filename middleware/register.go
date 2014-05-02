package middleware

import (
	"errors"
	"sync"

	"fourth.com/apibutler/jsonconfig"

	"github.com/codegangsta/martini"
)

type Constructor func(jsonconfig.Obj) (martini.Handler, error)

type ConfigItem struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Definition struct {
	ConfigItems  []*ConfigItem `json:"configItems"`
	Constructor  Constructor   `json:"-"`
	FriendlyName string        `json:"friendlyName"`
	Name         string        `json:"name"`
	Id           int           `json:"id"`
}

func NewDefinition(
	friendlyName string,
	ctor Constructor,
	configItems ...*ConfigItem) *Definition {

	return &Definition{
		ConfigItems:  configItems,
		Constructor:  ctor,
		FriendlyName: friendlyName,
	}
}

type Table map[string]*Definition

var (
	mu  sync.Mutex
	mws = make(Table)
)

func Register(name string, def *Definition) error {
	mu.Lock()
	defer mu.Unlock()
	if _, dup := mws[name]; dup {
		return errors.New("Duplicate registration of middleware constructor " + name)
	}

	def.Name = name
	mws[name] = def

	return nil
}

func GetMiddlewares() Table {
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
