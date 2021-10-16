package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   string
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lru *lruCache) Set(k Key, v interface{}) bool {
	if _, ok := lru.items[k]; ok {
		lru.items[k].Value = cacheItem{string(k), v}
		lru.queue.MoveToFront(lru.items[k])
		return true
	}

	lru.items[k] = lru.queue.PushFront(cacheItem{string(k), v})

	if lru.capacity < lru.queue.Len() {
		rm := lru.queue.Back().Value.(cacheItem)
		delete(lru.items, Key(rm.key))
		lru.queue.Remove(lru.queue.Back())
	}

	return false
}

func (lru *lruCache) Get(k Key) (interface{}, bool) {
	if _, ok := lru.items[k]; ok {
		lru.queue.MoveToFront(lru.items[k])
		v := lru.items[k].Value.(cacheItem)
		return v.value, true
	}

	return nil, false
}

func (lru *lruCache) Clear() {
	lru.items = nil
	lru.queue = NewList()
}
