package main

import (
	"sync"
	"time"
)

/*
* This file contains all oven definitions.
*
 */

type timedOven int

func (oven timedOven) Bake(unbakedPizza *Pizza) *Pizza {
	time.Sleep(time.Duration(oven) * time.Millisecond)
	unbakedPizza.isBaked = true
	return unbakedPizza
}

type ConcurrentOven struct {
	time int
	Lock *sync.Mutex
}

//up to caller to syncronize
func (oven ConcurrentOven) Bake(unbakedPizza *Pizza) *Pizza {
	time.Sleep(time.Duration(oven.time) * time.Millisecond)
	unbakedPizza.isBaked = true
	return unbakedPizza
}
func (oven ConcurrentOven) TryLock() bool {
	return oven.Lock.TryLock()
}
func (oven ConcurrentOven) Unlock() {
	oven.Lock.Unlock()
}
func newOven(time int) ConcurrentOven {
	return ConcurrentOven{time: time, Lock: new(sync.Mutex)}
}
