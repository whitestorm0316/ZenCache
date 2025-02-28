package cache

import (
	"sync"
	"zencache/internal/lru"
)

// 外部交互使用
type Engine struct {
	groups map[string]*Group
	mutex  sync.RWMutex
}

func NewEngine() *Engine {
	return &Engine{
		groups: make(map[string]*Group),
	}
}

func (e *Engine) GetGroup(name string) *Group {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.groups[name]
}
func (e *Engine) AddGroup(name string, getter Getter, maxBytes int64) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	g := &Group{
		cache: &cache{
			lru:      lru.New(maxBytes, nil),
			mu:       sync.RWMutex{},
			maxBytes: maxBytes,
		},
		getter: getter,
		name:   name,
	}
	e.groups[name] = g
}
