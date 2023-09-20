package main

/**
* Basic definintions from the task
 */
type Order struct {
}

type Pizza struct {
	isBaked bool
}

type Oven interface {
	Bake(unbakedPizza Pizza) Pizza
}

type PizzaBaker interface {
	ProcessOrder() Order

	Prepare(order Order) *Pizza

	QualityCheck(pizza *Pizza) (*Pizza, error)
}

//- Receive and process order (1ms)
const OrderTime int = 1

//- Prepare pizza (2ms)
const PrepareTime int = 2

//- Bake (5ms)
const BakeTime int = 5

//- Quality check (1ms)
const QualityCheck int = 1
