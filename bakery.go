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
	Bake(unbakedPizza *Pizza) *Pizza
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

type Pizzeria struct {
	Bakers         []PizzaBaker
	Ovens          []Oven
	orders         []Order
	preparedPizzas []*Pizza
	donePizzas     []*Pizza
}

//requires pizzeria has at least one baker and one oven
func (pizzeria Pizzeria) RunBakerySequentially(numOrders int) {
	for i := 0; i < numOrders; i++ {
		pizzeria.orders = append(pizzeria.orders, pizzeria.Bakers[0].ProcessOrder())
	}
	for _, o := range pizzeria.orders {
		pizzeria.preparedPizzas = append(pizzeria.preparedPizzas, pizzeria.Bakers[0].Prepare(o))
	}
	for _, p := range pizzeria.preparedPizzas {
		pizzeria.donePizzas = append(pizzeria.donePizzas, pizzeria.Ovens[0].Bake(p))
	}
}

func GetSimplePizzeria() *Pizzeria {
	var ret = new(Pizzeria)
	ret.Bakers = []PizzaBaker{BasicBaker{orderTime: OrderTime, prepareTime: PrepareTime, checkTime: QualityCheck}}
	ret.Ovens = []Oven{timedOven(BakeTime)}
	return ret
}
