package main

import (
	"errors"
	"time"
)

/**
* Define all bakers here
*
 */

type BasicBaker struct {
	orderTime   int
	prepareTime int
	checkTime   int
}

func (baker BasicBaker) ProcessOrder() Order {
	time.Sleep(time.Duration(baker.orderTime) * time.Millisecond)
	return Order{}
}

func (baker BasicBaker) Prepare(order Order) *Pizza {
	time.Sleep(time.Duration(baker.prepareTime) * time.Millisecond)
	return new(Pizza)
}

func (baker BasicBaker) QualityCheck(pizza *Pizza) (*Pizza, error) {
	time.Sleep(time.Duration(baker.prepareTime) * time.Millisecond)
	if !(pizza.isBaked) {
		return pizza, errors.New("not baked yet")
	}
	return pizza, nil
}
