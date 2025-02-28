package cache

import (
	"errors"
)

type GetterFunc func(string) ([]byte, error)
type Getter interface {
	Get(string) ([]byte, error)
}

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	ErrKeyNotFound = errors.New("KeyNotFound")
	ErrKeyIsNil    = errors.New("KeyIsNil")
)

// 命名空间
type Group struct {
	cache  *cache
	getter Getter
	name   string
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, ErrKeyIsNil
	}
	byteView, ok := g.cache.get(key)
	if !ok {
		// 回源
		if g.getter != nil {
			g.getter.Get(key)
		} else {
			return ByteView{}, ErrKeyNotFound
		}
	}
	return byteView, nil
}
func (g *Group) Add(key string, value ByteView) error {
	if key == "" {
		return ErrKeyIsNil
	}
	g.cache.add(key, value)
	return nil
}
