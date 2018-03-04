package main

import (
	"fmt"
	"sync"
)

func doSomething(i int, wg *sync.WaitGroup) {

	fmt.Println("Function in background", i)
	wg.Done()
}

func main() {
	nums := []int{2, 3, 4}
	var wg sync.WaitGroup
	for _, x := range nums {
		wg.Add(1)
		go doSomething(x, &wg)
	}
	wg.Wait()
}
