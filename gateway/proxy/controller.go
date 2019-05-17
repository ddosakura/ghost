package gateproxy

import (
	"sync"

	"github.com/ddosakura/ghost"
)

// Controller of Gateway-Proxy
type Controller struct {
	alive    bool
	mutex    *sync.RWMutex
	run      func()
	shutdown func()
}

func newController() *Controller {
	return &Controller{
		alive: false,
		mutex: &sync.RWMutex{},
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
	defer c.mutex.Unlock()
	c.mutex.Lock()
	if c.alive {
		return
	}
	c.alive = true
	if c.run != nil {
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
	defer c.mutex.Unlock()
	c.mutex.Lock()
	if !c.alive {
		return
	}
	c.alive = false
	if c.shutdown != nil {
		c.shutdown()
	}
}
