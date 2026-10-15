package main

import (
	"container/list"
	"sync"
)

type LRU[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Set(key K, value V)
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

type LRUCache[K comparable, V any] struct {
	capacity int
	mu       sync.Mutex
	ll       list.List
	items    map[K]*list.Element
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	c := &LRUCache[K, V]{capacity: capacity}
	if capacity > 0 {
		c.items = make(map[K]*list.Element, capacity)
	}
	return c
}

func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, found := c.items[key]
	if !found {
		return value, false
	}
	c.ll.MoveToFront(elem)
	return elem.Value.(*entry[K, V]).value, true
}

func (c *LRUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.capacity <= 0 {
		return
	}
	if elem, found := c.items[key]; found {
		elem.Value.(*entry[K, V]).value = value
		return
	}
	if c.ll.Len() >= c.capacity {
		oldest := c.ll.Back()
		c.ll.Remove(oldest)
		delete(c.items, oldest.Value.(*entry[K, V]).key)
	}
	elem := c.ll.PushFront(&entry[K, V]{key: key, value: value})
	c.items[key] = elem
}
