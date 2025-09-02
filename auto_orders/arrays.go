package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var n int
	n = rand.Intn(10)
	a := []int{}
	for i := 0; i < n; i++ {
		b := rand.Intn(15)
		fmt.Printf("Valor: %d - %d |", i, b)
		a = append(a, b)
	}
	fmt.Printf("%d", a)
}
