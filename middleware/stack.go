package middleware

import (
	"fmt"
	"strings"

	"fourth.com/apibutler/jsonconfig"
	"github.com/codegangsta/martini"
)

type Stack struct {
	middlewares []string
	configs     []jsonconfig.Obj
	reified     []martini.Handler
}

const (
	initialStackCapacity = 10
)

func NewStack() *Stack {
	return NewStackWithCapacity(initialStackCapacity)
}

func NewStackWithCapacity(capacity int) *Stack {
	return &Stack{
		middlewares: make([]string, 0, capacity),
		configs:     make([]jsonconfig.Obj, 0, capacity),
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
	for i, name := range s.middlewares {
		mw, err := Create(name, s.configs[i])
		if err != nil {
			mwerr.AddError(name, err)
		}
		s.reified = append(s.reified, mw)
	}

	if !mwerr.IsEmpty() {
		return mwerr
	}

	return nil
}

func (s *Stack) AddMiddleware(name string, cfg jsonconfig.Obj) {
	s.middlewares = append(s.middlewares, name)
	s.configs = append(s.configs, cfg)
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
