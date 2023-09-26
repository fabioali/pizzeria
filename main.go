package main

import (
	"fmt"
	"time"
)

// run simple pizzeria for numOrders of orders
func run_interval(numOrders []int, pizzeria Pizzeria) {
	fmt.Printf("===== start run ========\n")
	for _, v := range numOrders {
		tick := time.Now()
		pizzeria.RunBakerySequentially(v)
		tock := time.Now()
		fmt.Printf("%d: %d\n", v, tock.UnixMilli()-tick.UnixMilli())
	}
	fmt.Printf("===== end run ========\n")
}

// run the pizzeria with number of orders an measure the time in milliseconds
func run_interval2(numOrders []int, pizzeria *ConcurrentPizzeria) {
	fmt.Printf("===== start run latency ========\n")
	for _, v := range numOrders {
		pizzeria.Reset()
		tick := time.Now()
		pizzeria.runConcurrentPizzeria(v)
		tock := time.Now()
		fmt.Printf("%d:\t%d\n", v, tock.UnixMilli()-tick.UnixMilli())
		for i := 0; i < v; i++ {
			pizza := pizzeria.donePizzas.Next()
			if pizza == nil {
				fmt.Printf("Not the correct number of pizzas produced\n")
				break
			}
			//do we care about individual latency
			//latency := pizza.(*Pizza).finishTime - pizza.(*Pizza).issueTime
		}
	}
	fmt.Printf("===== end run latency ========\n")
}

// run benchmark of throughput times is how long the pizzeria is open for in milliseconds interval is how often a pizza is ordered(milliseconds)
// measures the number of pizzas produced in the time limit
func run_throughput(times []int, interval int, pizzeria *ConcurrentPizzeria) {
	fmt.Printf("===== start run throughput ========\n")
	for _, v := range times {
		pizzeria.Reset()
		pizzeria.runConcurrentPizzeriaBasedOnTime(int64(v), interval)
		numPizzas := 0
		for pizza := pizzeria.donePizzas.Next(); pizza != nil; pizza = pizzeria.donePizzas.Next() {
			//pizza.(*Pizza)
			numPizzas++
		}
		fmt.Printf("%d:\t%d\n", v, numPizzas)
	}
	fmt.Printf("===== end run throughput ========\n")
}

// Test the correctness of the pizzeria
// Checks if all pizzas are baked, there are no unbaked pizzas or open orders left
// Checks that each of the orders have been fulfilled
func test_correctness(numOrders []int, pizzeria *ConcurrentPizzeria) {
	failed := false
	for _, v := range numOrders {
		pizzeria.Reset()
		test := make([]bool, v)

		pizzeria.runConcurrentPizzeria(v)
		if pizzeria.orders.Next() != nil {
			fmt.Printf("Order not taken\n")
			failed = true
		}
		if pizzeria.preparedPizzas.Next() != nil {
			fmt.Printf("Pizza left in the oven\n")
			failed = true
		}
		for i := 0; i < v; i++ {
			pizza := pizzeria.donePizzas.Next()
			if pizza == nil {
				fmt.Printf("Not the correct number of pizzas produced\n")
				failed = true
				break
			}
			if !pizza.(*Pizza).isBaked {
				fmt.Printf("Pizza not baked\n")
				failed = true
			}
			if test[pizza.(*Pizza).orderId] {
				fmt.Printf("Two pizzas from same order %d\n", pizza.(*Pizza).orderId)
				failed = true
			}
			test[pizza.(*Pizza).orderId] = true
		}
		for i := 0; i < v; i++ {
			if !test[i] {
				fmt.Printf("Missing order %d\n", i)
				failed = true
			}
		}
	}
	if failed {
		fmt.Printf("Tests run failed\n")

	} else {
		fmt.Printf("Tests run succeeded\n")
	}
}
func main() {
	ranges := []int{1, 10, 100, 1000, 10000, 100000}
	fmt.Printf("Begin the baking\n")
	//var pizzeria = GetSimplePizzeria()
	//run_interval(ranges, *pizzeria)
	concurrentPizzeria := getConcurrentPizzeria(300, 500)
	test_correctness(ranges, concurrentPizzeria)
	run_interval2(ranges, concurrentPizzeria)
	run_throughput(ranges[1:4], 100000, concurrentPizzeria)
}
