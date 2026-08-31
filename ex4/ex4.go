package main

import (
	"fmt"
	"sync"
	"time"
)

var mutexA, mutexB sync.Mutex

func tarefa1(wg *sync.WaitGroup) {
	defer wg.Done()
	mutexA.Lock()
	fmt.Println("tarefa1: adquiriu A")
	time.Sleep(100 * time.Millisecond)
	mutexA.Unlock()
	mutexB.Lock() // espera por B, que está com tarefa2
	fmt.Println("tarefa1: adquiriu B")
	mutexB.Unlock()
	
}
func tarefa2(wg *sync.WaitGroup) {
	defer wg.Done()
	mutexB.Lock()
	fmt.Println("tarefa2: adquiriu B")
	time.Sleep(100 * time.Millisecond)
	mutexA.Lock() // espera por A, que está com tarefa1
	fmt.Println("tarefa2: adquiriu A")
	mutexA.Unlock()
	mutexB.Unlock()
}
func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go tarefa1(&wg)
	go tarefa2(&wg)
	wg.Wait()
	fmt.Println("Fim!")
}

// 1: O programa não termina. As linhas impressas antes do travamento são tarefa1: 
// adquiriu A e tarefa2: adquiriu B. A mensagem exibida pelo runtime do Go é o fatal error: 
// all goroutines are asleep - deadlock!.

// 2: Na tarefa1, a goroutine já possui o mutexA e está esperando pelo mutexB. 
// Na tarefa2, a goroutine já possui o mutexB e está esperando pelo mutexA. 
// ]A espera circular acontece porque a tarefa1 segura A e quer B, 
// enquanto a tarefa2 segura B e quer A, criando um ciclo fechado onde nenhuma pode avançar.

// 3: Executando com go run -race ex4.go, o detector não acusa corrida de dados. 
// A diferença é que um data race (corrida de dados) ocorre quando 
// duas goroutines acessam a mesma variável simultaneamente sem sincronização e 
// pelo menos um acesso é escrita. Já um deadlock é um problema de fluxo lógico 
// onde threads/goroutines ficam bloqueadas esperando recursos umas das outras. 
// O -race não ajuda porque o código está perfeitamente sincronizado em termos 
// de acesso à memória (os mutexes evitam o acesso simultâneo), sendo um problema 
// estritamente de ordem de bloqueio.