package main

import (
	"fmt"
	"time"
)

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
func run_interval2(numOrders []int, pizzeria *ConcurrentPizzeria) {
	fmt.Printf("===== start run ========\n")
	for _, v := range numOrders {
		pizzeria.Reset()
		tick := time.Now()
		pizzeria.runConcurrentPizzeria(v)
		tock := time.Now()
		fmt.Printf("%d:\t%d\n", v, tock.UnixMilli()-tick.UnixMilli())
	}
	fmt.Printf("===== end run ========\n")
}
func test_correctness(numOrders []int, pizzeria *ConcurrentPizzeria) {
	failed := false
	for _, v := range numOrders {
		pizzeria.Reset()
		test := make([]bool, v)

		pizzeria.runConcurrentPizzeria(v)
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
	concurrentPizzeria := getConcurrentPizzeria(1000, 800)
	run_interval2(ranges, concurrentPizzeria)
	//test_correctness(ranges, concurrentPizzeria)
}

//@ensures ret == a + b
func add(a int, b int) (ret int) {
	return a + b
}
