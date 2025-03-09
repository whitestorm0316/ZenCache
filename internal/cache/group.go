package cache

import (
	"errors"
	"zencache/internal/peers"
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
	cache       *cache // 从内存获取
	getter      Getter // 从本地获取
	name        string
	peersPicker peers.PeersPicker
}

func (g *Group) RegisterPicker(picker peers.PeersPicker) {
	g.peersPicker = picker
}
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, ErrKeyIsNil
	}
	peer, ok := g.peersPicker.PickPeer(key)
	if ok {
		bs, err := peer.Get(g.name, key)
		return NewByteView(bs), err
	} else {
		return g.getLocally(key)
	}
}

// 从本地获取
func (g *Group) getLocally(key string) (ByteView, error) {
	byteView, ok := g.cache.get(key)
	if !ok {
		// 回源
		if g.getter != nil {
			bs, err := g.getter.Get(key)
			if err != nil {
				return ByteView{}, err
			}
			byteView = NewByteView(bs)
			g.cache.add(key, byteView)
			return byteView, nil
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
