package tsmap

import "sync"

type TSMap struct {
	store map[string]interface{}
	l     *sync.RWMutex
}

func New() *TSMap {
	tmp := new(TSMap)
	tmp.l = new(sync.RWMutex)
	return tmp
}

func (this *TSMap) Release() {
	this.l.Lock()
	defer this.l.Unlock()
	if this.store != nil {
		for k, _ := range this.store {
			delete(this.store, k)
		}
	}
}

func (this *TSMap) Pop() map[string]interface{} {
	this.l.Lock()
	defer this.l.Unlock()
	tmp := this.store
	this.store = nil
	return tmp
}

func (this *TSMap) Push(s map[string]interface{}) {
	this.l.Lock()
	defer this.l.Unlock()
	this.store = s
}

func (this *TSMap) Set(key string, v interface{}) {
	this.l.Lock()
	defer this.l.Unlock()
	if this.store == nil {
		this.store = make(map[string]interface{})
	}
	this.store[key] = v
}

func (this *TSMap) Get(key string) interface{} {
	this.l.RLock()
	defer this.l.RUnlock()
	if this.store == nil {
		return nil
	}
	return this.store[key]
}

func (this *TSMap) Has(key string) bool {
	this.l.RLock()
	defer this.l.RUnlock()
	if this.store == nil {
		return false
	}
	_, ok := this.store[key]
	return ok
}

func (this *TSMap) Del(key string) {
	this.l.Lock()
	defer this.l.Unlock()
	if this.store != nil {
		if _, ok := this.store[key]; ok {
			delete(this.store, key)
		}
	}
}
