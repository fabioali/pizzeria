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

func main() {
	ranges := []int{1, 10, 100, 1000, 10000}
	fmt.Printf("Begin the baking\n")
	var pizzeria = GetSimplePizzeria()
	run_intervall(ranges, *pizzeria)
}

//@ensures ret == a + b
func add(a int, b int) (ret int) {
	return a + b
}
