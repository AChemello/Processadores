package main
import (
"fmt"
"sync"
)
var {
 contador int
 mu sync.Mutex
}
func incrementar(n int, wg *sync.WaitGroup) {
defer wg.Done()
for i := 0; i < n; i++ {
 contador++
}
}
func main() {
var wg sync.WaitGroup
 wg.Add(2)
mu.Lock()
go incrementar(2000000, &wg)
mu.Lock()
go incrementar(2000000, &wg)
mu.Unlock()
 wg.Wait()
 fmt.Printf("Contador: %d (esperado: 4000000)\n", contador)
}
