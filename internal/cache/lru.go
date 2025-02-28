package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	onEvicted func(key string, value Value)
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

type Value interface {
	Len() int
}
type entry struct {
	key   string
	value Value
}

func (c *Cache) Get(key string) (Value, bool) {
	element, ok := c.cache[key]
	if !ok {
		return nil, ok
	}
	c.ll.MoveToFront(element)
	return element.Value.(*entry).value, ok
}

func (c *Cache) Add(key string, value Value) {
	element, ok := c.cache[key]
	if ok {
		c.nBytes += (int64(value.Len()) - int64(element.Value.(*entry).value.Len()))
		element.Value = &entry{
			key:   key,
			value: value,
		}
		c.ll.MoveToFront(element)
		return
	}
	element = c.ll.PushFront(&entry{
		key:   key,
		value: value,
	})
	c.nBytes += (int64(len(key) + value.Len()))
	c.cache[key] = element
	for c.Len() > 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}
func (c *Cache) RemoveOldest() {
	element := c.ll.Back()
	c.nBytes -= int64(len(element.Value.(*entry).key) + element.Value.(*entry).value.Len())
	c.ll.Remove(element)
	delete(c.cache, element.Value.(*entry).key)
	if c.onEvicted != nil {
		c.onEvicted(element.Value.(*entry).key, element.Value.(*entry).value)
	}
}
func (c *Cache) Len() int {
	return c.ll.Len()
}
