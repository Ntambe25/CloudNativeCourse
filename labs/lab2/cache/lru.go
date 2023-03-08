package cache

import (
	"errors"
)

// Cacher interface
type Cacher interface {
	Get(interface{}) (interface{}, error)
	Put(interface{}, interface{}) error
}

// Struct for a new lrucache
type lruCache struct {
	Size      int
	Remaining int
	Cache     map[string]string
	Queue     []string
}

// Make a new Empty Cache of size "Size"
// Makes an empty map that is the cahce
// Makes an empty slice that is the queue
func NewCache(size int) Cacher {
	return &lruCache{Size: size, Remaining: size, Cache: make(map[string]string), Queue: make([]string, 0)}
}

// Input is the memory location of lru(lrucache)
func (lru *lruCache) Get(key interface{}) (interface{}, error) {

	//If k is already present in the queue, remove the duplicate
	for _, v := range lru.Queue {

		if v == key.(string) {
			lru.qDel(key.(string))
			lru.Remaining++
		}
	}

	//If the key is found
	if val, found := lru.Cache[key.(string)]; found {

		//Append to the tail of the queue
		lru.Queue = append(lru.Queue, key.(string))

		//Remaining is Decreased since an element from the queue
		//has been deleted
		lru.Remaining--

		return val, nil
	}
	return nil, errors.New("Error: Key not found")

}

func (lru *lruCache) Put(key, val interface{}) error {

	//If k is already present in the queue, remove the duplicate
	for _, v := range lru.Queue {

		if v == key.(string) {
			lru.qDel(key.(string))
			lru.Remaining++
		}
	}

	//If Remaining is 0, then Queue is full
	if lru.Remaining == 0 {

		//Delete the last element of the queue from the cache
		delete(lru.Cache, lru.Queue[0])

		//Delete last element of the queue
		lru.qDel(lru.Queue[0])

		//Remaining is incresead since an element from the queue
		//has been deleted
		lru.Remaining++
	}

	//Append the new key and value to the tail of the queue
	lru.Queue = append(lru.Queue, key.(string))

	//Insert the new key and value to the Cache
	lru.Cache[key.(string)] = val.(string)

	//Decrement Remaining
	lru.Remaining--

	return nil
}

// Delete element from queue
func (lru *lruCache) qDel(ele string) {
	for i := 0; i < len(lru.Queue); i++ {
		if lru.Queue[i] == ele {
			oldlen := len(lru.Queue)
			copy(lru.Queue[i:], lru.Queue[i+1:])
			lru.Queue = lru.Queue[:oldlen-1]
			break
		}
	}
}
