package main

import (
	"fmt"
	"sync"
)

var (
	contador int
	mu       sync.Mutex
)

func incrementar(n int, wg *sync.WaitGroup) {
	mu.Lock()
	defer wg.Done()
	for i := 0; i < n; i++ {
		contador++
	}
	mu.Unlock()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go incrementar(2000000, &wg)

	go incrementar(2000000, &wg)

	wg.Wait()
	fmt.Printf("Contador: %d (esperado: 4000000)\n", contador)

}

//1 Não
//2 a saida indica a existencia de 3 dataraces
//3 sim
//4time go run exsem/ex3.go

// ----------------------------------------------------------------------

// time go run ex3sem/ex3.go 
// time go run ex3/ex3.go 
// Contador: 3704870 (esperado: 4000000)

// real    0m0.190s
// user    0m0.072s
// sys     0m0.062s
// Contador: 4000000 (esperado: 4000000)

// real    0m0.219s
// user    0m0.153s
// sys     0m0.083s

// -----------------------------------------------------------------

// time go run ex3sem/ex3.go 
// time go run ex3/ex3.go 
// Contador: 2342092 (esperado: 4000000)

// real    0m0.084s
// user    0m0.076s
// sys     0m0.051s
// Contador: 4000000 (esperado: 4000000)

// real    0m0.081s
// user    0m0.059s
// sys     0m0.056s

// ----------------------------------------------------------------------

// time go run ex3sem/ex3.go 
// time go run ex3/ex3.go 

// Contador: 2016910 (esperado: 4000000)

// real    0m0.084s
// user    0m0.088s
// sys     0m0.042s
// Contador: 4000000 (esperado: 4000000)

// real    0m0.078s
// user    0m0.059s
// sys     0m0.055s

//Sim, o mutex adiciona overhead (sobrecarga) em termos de desempenho de CPU e tempo de execução.