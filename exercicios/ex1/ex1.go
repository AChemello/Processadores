// Exercício 1 — Solução final utilizando canais
//
// Este programa substitui o uso de sync.WaitGroup por 
// canais (channels) como mecanismo de sincronização.
package main

import (
	"fmt"
	"time"
)

// Requisito b) Cada função agora recebe o canal e envia 'true' ao terminar
func crescente(done chan bool) {
	for i := 1; i <= 10; i++ {
		fmt.Printf("[Crescente] %d\n", i)
		time.Sleep(1 * time.Second)
	}
	done <- true 
}

func decrescente(done chan bool) {
	for i := 10; i >= 1; i-- {
		fmt.Printf("[Decrescente] %d\n", i)
		time.Sleep(1 * time.Second)
	}
	done <- true
}

func main() {
	// Requisito a) Canal adequado criado para receber a sinalização de término
	done := make(chan bool)

	// Inicializa as goroutines passando o canal como argumento
	go crescente(done)
	go decrescente(done)

	// Requisito c) A main() recebe do canal duas vezes antes do print final
	<-done // Bloqueia esperando a primeira goroutine terminar
	<-done // Bloqueia esperando a segunda goroutine terminar

	// Requisito d) O programa não usa sync.WaitGroup
	fmt.Println("Fim!")
}
