package main

import (
	"time"
)

/*
* This file contains all oven definitions.
*
 */

type timedOven int

func (oven timedOven) Bake(unbakedPizza *Pizza) *Pizza {
	time.Sleep(time.Duration(oven) * time.Millisecond)
	return unbakedPizza
}
