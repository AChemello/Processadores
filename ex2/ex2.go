package main

import (
	"fmt"
	"sync"
	"time"
)

func contarCrescente(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		fmt.Printf("Crescente: %d\n", i)
		time.Sleep(1 * time.Second)
	}
}

func contarDecrescente(wg *sync.WaitGroup) {//contador de sincronização
	defer wg.Done()// decrementa ao final
	for i := 10; i >= 1; i-- {
		fmt.Printf("Decrescente: %d\n", i)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	go contarCrescente(&wg)
	go contarDecrescente(&wg)

	wg.Wait()
	fmt.Println("Ambas as goroutines finalizaram!")
}

//go run ex1.go