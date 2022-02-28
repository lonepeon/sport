package job

import (
	"context"
	"sync"
)

type Handler interface {
	Name() string
	Handle(context.Context, []byte) error
}

type HandlerFunc func(context.Context, []byte) error

type Registry struct {
	registry map[string]HandlerFunc
	l        *sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		l:        &sync.RWMutex{},
		registry: make(map[string]HandlerFunc),
	}
}

func (r *Registry) Register(job Handler) {
	r.RegisterFunc(job.Name(), job.Handle)
}

func (r *Registry) RegisterFunc(name string, handler HandlerFunc) {
	r.l.Lock()
	defer r.l.Unlock()
	r.registry[name] = handler
}

func (r *Registry) Handler(name string) (HandlerFunc, bool) {
	r.l.RLock()
	defer r.l.RUnlock()
	h, ok := r.registry[name]
	return h, ok
}
