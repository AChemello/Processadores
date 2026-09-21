// Exercício 3 — Pipeline de processamento
//
// Complete as funções gerador, multiplicador e impressora
// para formar um pipeline de 3 estágios conectados por canais.
package main

import "fmt"

// gerador envia os números de 'valores' no canal 'out' e fecha o canal.
func gerador(valores []int, out chan<- int) {

}

// multiplicador lê do canal 'in', multiplica por 2, e envia no canal 'out'.
// Fecha 'out' ao terminar.
func multiplicador(in <-chan int, out chan<- int) {

}

// impressora lê do canal 'in' e imprime cada valor.
// Envia true no canal 'done' ao terminar.
func impressora(in <-chan int, done chan<- bool) {

}

func main() {
	valores := []int{1, 2, 3, 4, 5}

	c1 := make(chan int)
	c2 := make(chan int)
	done := make(chan bool)

	go gerador(valores, c1)
	go multiplicador(c1, c2)
	go impressora(c2, done)

	<-done
	fmt.Println("Pipeline concluído!")
}
