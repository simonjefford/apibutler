package middleware

import (
	"fmt"
	"strings"

	"labix.org/v2/mgo/bson"

	"fourth.com/apibutler/jsonconfig"
	"github.com/codegangsta/martini"
)

type Stack struct {
	Middlewares []string          `json:"middlewares"`
	Configs     []jsonconfig.Obj  `json:"configs"`
	ID          bson.ObjectId     `bson:"_id" json:"id"`
	reified     []martini.Handler `json:"-"`
}

const (
	initialStackCapacity = 10
)

func NewStack() *Stack {
	return NewStackWithCapacity(initialStackCapacity)
}

func NewStackWithCapacity(capacity int) *Stack {
	return &Stack{
		Middlewares: make([]string, 0, capacity),
		Configs:     make([]jsonconfig.Obj, 0, capacity),
	}
}

func (s *Stack) AddToServer(srv *martini.Martini) error {
	if s.reified == nil {
		err := s.reify()
		if err != nil {
			return err
		}
	}

	for _, mw := range s.reified {
		srv.Use(mw)
	}

	return nil
}

func (s *Stack) reify() error {
	mwerr := &MiddlewareStackError{}
	for i, name := range s.Middlewares {
		mw, err := Create(name, s.Configs[i])
		if err != nil {
			mwerr.AddError(name, err)
		}
		s.reified = append(s.reified, mw)
	}

	if !mwerr.IsEmpty() {
		s.reified = nil
		return mwerr
	}

	return nil
}

func (s *Stack) AddMiddleware(name string, cfg jsonconfig.Obj) {
	s.Middlewares = append(s.Middlewares, name)
	s.Configs = append(s.Configs, cfg)
	s.reified = nil
}

type MiddlewareStackError struct {
	errors map[string]error
}

func (e *MiddlewareStackError) AddError(name string, err error) {
	if e.IsEmpty() {
		e.errors = make(map[string]error)
	}
	e.errors[name] = err
}

func (e *MiddlewareStackError) IsEmpty() bool {
	return e.errors == nil
}

func (e *MiddlewareStackError) Error() string {
	if e.IsEmpty() {
		return ""
	}

	l := make([]string, 0, len(e.errors))

	for name, err := range e.errors {
		l = append(l, fmt.Sprintf("Middleware error with %s - %v", name, err))
	}

	return strings.Join(l, "\n")
}
