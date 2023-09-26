package main

import (
	"sync"

	fifo "github.com/foize/go.fifo"
	pq "github.com/jupp0r/go-priority-queue"
)

/**
* Basic definintions from the task
 */
type Order struct {
	id int
}

type Pizza struct {
	isBaked bool
	orderId int
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
	for _, p := range pizzeria.preparedPizzas {
		pizzeria.Bakers[0].QualityCheck(p)
	}
}

func GetSimplePizzeria() *Pizzeria {
	var ret = new(Pizzeria)
	ret.Bakers = []PizzaBaker{BasicBaker{orderTime: OrderTime, prepareTime: PrepareTime, checkTime: QualityCheck}}
	ret.Ovens = []Oven{timedOven(BakeTime)}
	return ret
}

type task int
type taskWrapper struct {
	Task task
	id   int
}

const (
	takeOrder    = 0
	preparePizza = 1
	bakePizza    = 2
	checkPizza   = 3
)

// make Concurrent pizzeria
type ConcurrentPizzeria struct {
	Bakers         []ConcurrentBaker
	Ovens          []ConcurrentOven
	orders         fifo.Queue
	preparedPizzas fifo.Queue
	donePizzas     fifo.Queue
	taskList       pq.PriorityQueue
	taskIds        int // also secured via taskLok
	numberOfOrders int
	taskLock       sync.Mutex
}

//TODO verify that all pizzas are baked and all orders are completed
func (pizzeria *ConcurrentPizzeria) runConcurrentPizzeria(numOrders int) {
	pizzeria.taskLock.Lock()
	pizzeria.numberOfOrders = numOrders
	pizzeria.taskLock.Unlock()
	for i := 0; i < numOrders; i++ {
		pizzeria.AddTask(takeOrder)
	}
	var wg = sync.WaitGroup{}
	for _, v := range pizzeria.Bakers {
		wg.Add(1)
		go func(baker ConcurrentBaker) {
			defer wg.Done()
			baker.run()
		}(v)
	}
	//TODO: synchronize the stop
	wg.Wait()
	if pizzeria.orders.Len() > 0 {
		//this fails because the fifo library is faulty (the size keeps increasing even if it is empty)
		//fmt.Printf("Not all orders are taken: %d\n", pizzeria.orders.Len())
	}
}
func (pizzeria *ConcurrentPizzeria) GetTask() (taskWrapper, error) {

	pizzeria.taskLock.Lock()
	ret, err := pizzeria.taskList.Pop()
	pizzeria.taskLock.Unlock()
	if err != nil {
		return taskWrapper{}, err
	}
	return taskWrapper(ret.(taskWrapper)), nil
}
func (pizzeria *ConcurrentPizzeria) AddTask(t task) {

	pizzeria.taskLock.Lock()
	id := pizzeria.taskIds
	pizzeria.taskIds += 1
	pizzeria.taskList.Insert(taskWrapper{t, id}, 10-float64(t))
	pizzeria.taskLock.Unlock()
}

func (pizzeria *ConcurrentPizzeria) Bake(pizza *Pizza) *Pizza {
	for {
		for _, v := range pizzeria.Ovens {
			if v.TryLock() {
				pizza = v.Bake(pizza)
				v.Unlock()
				return pizza
			}
		}
	}

}
func (pizzeria *ConcurrentPizzeria) AddPizza(pizza *Pizza, baked bool) {
	if baked {
		pizzeria.donePizzas.Add(pizza)
	} else {
		pizzeria.preparedPizzas.Add(pizza)
	}
}
func (pizzeria *ConcurrentPizzeria) GetPizza(baked bool) *Pizza {
	var p interface{}
	if baked {
		p = pizzeria.donePizzas.Next()
	} else {
		p = pizzeria.preparedPizzas.Next()
	}
	if p == nil {
		return nil
	}
	return p.(*Pizza)
}

func (pizzeria *ConcurrentPizzeria) AddOrder(order Order) {
	pizzeria.orders.Add(order)
}
func (pizzeria *ConcurrentPizzeria) GetOrder() Order {
	return pizzeria.orders.Next().(Order)
}
func (pizzeria *ConcurrentPizzeria) isDone() bool {
	ret := false
	pizzeria.taskLock.Lock()
	if pizzeria.donePizzas.Len() >= pizzeria.numberOfOrders {
		ret = true
	}
	pizzeria.taskLock.Unlock()
	return ret
}
func (pizzeria *ConcurrentPizzeria) Reset() {
	pizzeria.taskLock.Lock()
	pizzeria.numberOfOrders = 0
	pizzeria.taskIds = 0
	pizzeria.taskLock.Unlock()
}
func getConcurrentPizzeria(numBakers int, numOven int) *ConcurrentPizzeria {
	ret := new(ConcurrentPizzeria)
	for i := 0; i < numBakers; i++ {
		ret.Bakers = append(ret.Bakers, newConcurrentBaker(ret))
	}
	for i := 0; i < numOven; i++ {
		ret.Ovens = append(ret.Ovens, newOven(BakeTime))
	}
	ret.taskList = pq.New()
	ret.donePizzas = *fifo.NewQueue()
	ret.preparedPizzas = *fifo.NewQueue()
	ret.orders = *fifo.NewQueue()
	return ret
}
