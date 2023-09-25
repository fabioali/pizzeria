package main

import (
	"fmt"
	"time"
)

func run_intervall(numOrders []int, pizzeria Pizzeria) {
	fmt.Printf("===== start run ========\n")
	for _, v := range numOrders {
		tick := time.Now()
		pizzeria.RunBakerySequentially(v)
		tock := time.Now()
		fmt.Printf("%d: %d\n", v, tock.UnixMilli()-tick.UnixMilli())
	}
	fmt.Printf("===== end run ========\n")
}
func run_intervall2(numOrders []int, pizzeria *ConcurrentPizzeria) {
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
func main() {
	ranges := []int{1, 10, 100, 1000, 10000, 100000}
	fmt.Printf("Begin the baking\n")
	//var pizzeria = GetSimplePizzeria()
	//run_intervall(ranges, *pizzeria)
	concurrentPizzeria := getConcurrentPizzeria(300, 400)
	run_intervall2(ranges, concurrentPizzeria)

}

//@ensures ret == a + b
func add(a int, b int) (ret int) {
	return a + b
}
