// Produtor-consumidor simples com duas goroutines.
//
// Uma goroutine produz itens e envia pelo canal.
// Outra goroutine consome e imprime cada item.
// O produtor fecha o canal ao terminar, sinalizando
// ao consumidor que não haverá mais dados.
package main

import (
	"fmt"
	"time"
)

func produtor(ch chan<- string) {
	for i := 1; i <= 5; i++ {
		item := fmt.Sprintf("item-%d", i)
		fmt.Printf("[Produtor]   Produzindo %s\n", item)
		ch <- item
		time.Sleep(500 * time.Millisecond)
	}
	close(ch)
	fmt.Println("[Produtor] Canal fechado")
}

func consumidor(ch <-chan string, done chan<- bool) {
	for item := range ch {
		fmt.Printf("[Consumidor] Recebeu %s\n", item)
		time.Sleep(1 * time.Second)
	}
	done <- true
}

func main() {
	ch := make(chan string)
	done := make(chan bool)

	go produtor(ch)
	go consumidor(ch, done)

	fmt.Println("[Main] Esperando consumidor...")
	<-done
	fmt.Println("Fim!")
}
