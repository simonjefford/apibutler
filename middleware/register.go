package middleware

import (
	"errors"
	"sort"
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
	Schema       []*ConfigItem `json:"schema"`
	Constructor  Constructor   `json:"-"`
	FriendlyName string        `json:"friendlyName"`
	Name         string        `json:"name"`
	Id           string        `json:"id"`
}

func NewDefinition(
	friendlyName string,
	ctor Constructor,
	schema ...*ConfigItem) *Definition {

	return &Definition{
		Schema:       schema,
		Constructor:  ctor,
		FriendlyName: friendlyName,
	}
}

type Table map[string]*Definition

type List []*Definition

func (l List) Len() int {
	return len(l)
}

func (l List) Less(i, j int) bool {
	return l[i].Id < l[j].Id
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

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
	def.Id = name
	mws[name] = def

	return nil
}

// For unit testing only
func clearTable() {
	mu.Lock()
	defer mu.Unlock()

	mws = make(Table)
}

func Definitions() List {
	mu.Lock()
	defer mu.Unlock()

	t := List(make([]*Definition, 0, len(mws)))
	for _, val := range mws {
		t = append(t, val)
	}

	sort.Sort(t)

	return t
}

func GetTable() Table {
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
