package utils

import (
	"sync"
)

type UpdateCtx struct {
	v   map[string]bool
	mux sync.Mutex
}

func NewUpdateCtx() *UpdateCtx {
	return &UpdateCtx{
		v: make(map[string]bool),
	}
}

// return true if udpate can be started
func (c *UpdateCtx) Add(ref string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	if _, ok := c.v[ref]; ok {
		c.v[ref] = true
		return false
	}
	c.v[ref] = false
	return true
}

// return true if another update run must be started
func (c *UpdateCtx) Remove(ref string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	if val, ok := c.v[ref]; ok && val {
		c.v[ref] = false
		return true
	}
	delete(c.v, ref)
	return false
}
