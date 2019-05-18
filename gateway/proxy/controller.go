package gateproxy

import (
	"sync"

	"github.com/ddosakura/ghost"
)

// TODO: 提取为单独的 Lib

// Controller of Gateway-Proxy
type Controller struct {
	mutex          *sync.RWMutex
	alive          bool
	aliveListeners map[string]func(bool)
	run            func()
	shutdown       func()
}

func newController() *Controller {
	return &Controller{
		alive:          false,
		mutex:          &sync.RWMutex{},
		aliveListeners: make(map[string]func(bool)),
	}
}

// Alive State
func (c *Controller) Alive() bool {
	defer c.mutex.RUnlock()
	c.mutex.RLock()
	return c.alive
}

// Start Gateway-Proxy
func (c *Controller) Start() {
	change := c.alivePub(true)
	if change && c.run != nil {
		go func() {
			defer func() {
				e := recover()
				if e != nil {
					ghost.ErrorInDefer(e)
					c.Stop()
				}
			}()
			c.run()
		}()
	}
}

// Stop Gateway-Proxy
func (c *Controller) Stop() {
	change := c.alivePub(false)
	if change && c.shutdown != nil {
		c.shutdown()
	}
}

func (c *Controller) alivePub(b bool) (change bool) {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	if c.alive == b {
		return false
	}
	c.alive = b
	for _, fn := range c.aliveListeners {
		fn(c.alive)
	}
	return true
}

// Sub Alive
func (c *Controller) Sub(label string, fn func(bool)) {
	if fn == nil {
		c.Unsub(label)
		return
	}
	defer c.mutex.Unlock()
	c.mutex.Lock()
	c.aliveListeners[label] = fn
}

// Unsub Alive
func (c *Controller) Unsub(label string) {
	defer c.mutex.Unlock()
	c.mutex.Lock()
	delete(c.aliveListeners, label)
}
