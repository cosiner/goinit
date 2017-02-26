package goinit

import "sync"

var DefaultRegister = NewRegister()

func NewRegister() *Register {
	return &Register{
		sliceCache: make(map[Category][]interface{}),
		mapCache:   make(map[Category]map[string]interface{}),
	}
}

type Category int32

type Register struct {
	mu         sync.RWMutex
	category   int32
	sliceCache map[Category][]interface{}
	mapCache   map[Category]map[string]interface{}
}

func (r *Register) NewCategory() Category {
	r.mu.Lock()
	c := r.category + 1
	r.category = c
	r.mu.Unlock()

	return Category(c)
}

func (r *Register) Append(category Category, val ...interface{}) *Register {
	r.mu.Lock()
	r.sliceCache[category] = append(r.sliceCache[category], val...)
	r.mu.Unlock()

	return r
}

func (r *Register) Set(category Category, key string, val interface{}) *Register {
	r.mu.Lock()
	m, has := r.mapCache[category]
	if !has {
		m = make(map[string]interface{})
		r.mapCache[category] = m
	}
	m[key] = val
	r.mu.Unlock()
	return r
}

func (r *Register) Slice(category Category) []interface{} {
	r.mu.RLock()
	c := r.sliceCache[category]
	r.mu.RUnlock()
	return c
}

func (r *Register) Map(category Category) map[string]interface{} {
	r.mu.RLock()
	c := r.mapCache[category]
	r.mu.RUnlock()
	return c
}

func (r *Register) RemoveCategory(category ...Category) {
	r.mu.Lock()
	for _, c := range category {
		delete(r.sliceCache, c)
		delete(r.mapCache, c)
	}
	r.mu.Unlock()
}

func (r *Register) Destroy() {
	r.mu.Lock()
	r.sliceCache = nil
	r.mapCache = nil
	r.mu.Unlock()
}
