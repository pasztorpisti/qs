// +build !go1.9

package qs

import "sync"

func newSyncMap() syncMap {
	return syncMap{
		m: make(map[interface{}]interface{}),
	}
}

type syncMap struct {
	mu sync.RWMutex
	m  map[interface{}]interface{}
}

func (o *syncMap) Load(key interface{}) (value interface{}, ok bool) {
	o.mu.RLock()
	value, ok = o.m[key]
	o.mu.RUnlock()
	return
}

func (o *syncMap) Store(key, value interface{}) {
	o.mu.Lock()
	o.m[key] = value
	o.mu.Unlock()
}
