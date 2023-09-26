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
	ret := new(Pizza)
	ret.orderId = order.id
	return ret
}

func (baker BasicBaker) QualityCheck(pizza *Pizza) (*Pizza, error) {
	if pizza == nil {
		return nil, errors.New("Nil")
	}
	time.Sleep(time.Duration(baker.prepareTime) * time.Millisecond)
	if !(pizza.isBaked) {
		return pizza, errors.New("not baked yet")
	}
	return pizza, nil
}

// make baker easier to manage
type ConcurrentBaker struct {
	orderTime   int
	prepareTime int
	checkTime   int
	isBusy      bool
	pizzeria    *ConcurrentPizzeria
}

func (baker ConcurrentBaker) ProcessOrder(id int) Order {
	baker.isBusy = true
	time.Sleep(time.Duration(baker.orderTime) * time.Millisecond)

	baker.isBusy = false
	return Order{id: id}
}

func (baker ConcurrentBaker) Prepare(order Order) *Pizza {
	baker.isBusy = true
	time.Sleep(time.Duration(baker.prepareTime) * time.Millisecond)
	ret := new(Pizza)
	ret.orderId = order.id
	baker.isBusy = false
	return ret
}

func (baker ConcurrentBaker) QualityCheck(pizza *Pizza) (*Pizza, error) {
	if pizza == nil {
		return nil, errors.New("Nil")
	}
	baker.isBusy = true
	time.Sleep(time.Duration(baker.prepareTime) * time.Millisecond)
	baker.isBusy = false
	if !(pizza.isBaked) {
		return pizza, errors.New("not baked yet")
	}
	return pizza, nil
}
func newConcurrentBaker(pizzeria *ConcurrentPizzeria) ConcurrentBaker {
	return ConcurrentBaker{orderTime: OrderTime, prepareTime: PrepareTime, checkTime: QualityCheck, pizzeria: pizzeria}
}
func (baker ConcurrentBaker) run() {
	for {
		t, err := baker.pizzeria.GetTask()
		if err != nil {
			//fmt.Printf("no tasks left\n")
			time.Sleep(time.Millisecond)
			if baker.pizzeria.isDone() {
				return
			}
			continue
		}
		switch t.Task {
		case takeOrder:
			//fmt.Printf("take order %d\n", t.id)
			order := baker.ProcessOrder(t.id)
			baker.pizzeria.AddOrder(order)
			baker.pizzeria.AddTask(preparePizza)
		case preparePizza:
			//fmt.Printf("prepare pizza\n")
			order := baker.pizzeria.GetOrder()
			pizza := baker.Prepare(order)
			baker.pizzeria.AddPizza(pizza, false)
			baker.pizzeria.AddTask(bakePizza)
		case bakePizza:
			//fmt.Printf("bake pizza\n")
			pizza := baker.pizzeria.GetPizza(false)
			p := baker.pizzeria.Bake(pizza)
			baker.pizzeria.AddPizza(p, true)
			baker.pizzeria.AddTask(checkPizza)
		case checkPizza:
			//fmt.Printf("check pizza\n")
			//TODO
			p := baker.pizzeria.GetPizza(true)
			pizza, err := baker.QualityCheck(p)
			if err != nil {
				baker.pizzeria.AddPizza(pizza, false)
				baker.pizzeria.AddTask(bakePizza)
			}
			baker.pizzeria.AddPizza(pizza, true)

		}
	}
}
