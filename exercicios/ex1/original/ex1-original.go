// Exercício 1 — Ponto de partida
//
// Este programa usa sync.WaitGroup para a main() esperar
// duas goroutines. Sua tarefa: substituir o WaitGroup por
// canais como mecanismo de sinalização.
package main

import (
	"fmt"
	"sync"
	"time"
)

func crescente(wg *sync.WaitGroup) {
	for i := 1; i <= 10; i++ {
		fmt.Printf("[Crescente] %d\n", i)
		time.Sleep(1 * time.Second)
	}
	wg.Done()
}

func decrescente(wg *sync.WaitGroup) {
	for i := 10; i >= 1; i-- {
		fmt.Printf("[Decrescente] %d\n", i)
		time.Sleep(1 * time.Second)
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	go crescente(&wg)
	go decrescente(&wg)

	wg.Wait()
	fmt.Println("Fim!")
}
