package lru

import (
	"container/list"
	"errors"
)

type CacheNode struct {
	Key, Value interface{}
}

func (c *CacheNode) NewCacheNode(k, v interface{}) *CacheNode {
	return &CacheNode{k, v}
}

type LRUCache struct {
	Capacity int
	dlist    *list.List
	cacheMap map[interface{}]*list.Element
}

func NewLRUCache(cap int) *LRUCache {
	return &LRUCache{
		Capacity: cap,
		dlist:    list.New(),
		cacheMap: make(map[interface{}]*list.Element)}
}

func (lru *LRUCache) Size() int {
	return lru.dlist.Len()
}

func (lru *LRUCache) Set(k, v interface{}) {

	if lru.dlist == nil {
		panic(errors.New("cacheMap=nil"))
	}

	if pElement, ok := lru.cacheMap[k]; ok {
		lru.dlist.MoveToFront(pElement)
		pElement.Value.(*CacheNode).Value = v
		return
	}

	newElement := lru.dlist.PushFront(&CacheNode{k, v})
	lru.cacheMap[k] = newElement

	if lru.dlist.Len() > lru.Capacity {
		lastElement := lru.dlist.Back()
		if lastElement == nil {
			return
		}
		cacheNode := lastElement.Value.(*CacheNode)
		delete(lru.cacheMap, cacheNode.Key)
		lru.dlist.Remove(lastElement)
	}
}

func (lru *LRUCache) Get(k interface{}) (v interface{}, ret bool) {

	if lru.cacheMap == nil {
		panic(errors.New("cacheMap=nil"))
	}

	if pElement, ok := lru.cacheMap[k]; ok {
		lru.dlist.MoveToFront(pElement)
		return pElement.Value.(*CacheNode).Value, true
	}
	return v, false
}

func (lru *LRUCache) Remove(k interface{}) bool {

	if lru.cacheMap == nil {
		return false
	}

	if pElement, ok := lru.cacheMap[k]; ok {
		cacheNode := pElement.Value.(*CacheNode)
		delete(lru.cacheMap, cacheNode.Key)
		lru.dlist.Remove(pElement)
		return true
	}
	return false
}
