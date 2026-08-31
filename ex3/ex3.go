package main

import (
	"fmt"
	"sync"
)
<<<<<<< HEAD
var {
 contador int
 mu sync.Mutex
}
=======

var contador int

>>>>>>> dace4d7 (fix: exec3)
func incrementar(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		contador++
	}
}
func main() {
<<<<<<< HEAD
var wg sync.WaitGroup
 wg.Add(2)
mu.Lock()
go incrementar(2000000, &wg)
mu.Lock()
go incrementar(2000000, &wg)
mu.Unlock()
 wg.Wait()
 fmt.Printf("Contador: %d (esperado: 4000000)\n", contador)
=======
	var wg sync.WaitGroup
	wg.Add(2)
	go incrementar(2000000, &wg)
	go incrementar(2000000, &wg)
	wg.Wait()
	fmt.Printf("Contador: %d (esperado: 4000000)\n", contador)
>>>>>>> dace4d7 (fix: exec3)
}
