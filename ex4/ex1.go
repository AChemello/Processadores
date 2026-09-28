
// ex1.go — Solução Completa com Timeout
package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	capacidade   = 3
	totalPessoas = 10
)

// Semáforo global limitando a capacidade do banheiro em até 3 pessoas simultâneas
var sem = semaphore.NewWeighted(capacidade)

func usarBanheiro(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("[Pessoa %2d] quer usar o banheiro\n", id)

	// Cada pessoa define aleatoriamente quanto tempo está disposta a esperar (1 a 3 segundos)
	tempoMaxEspera := time.Duration(1+rand.Intn(3)) * time.Second

	// Criação de um contexto com timeout para controlar o limite de tempo da fila
	ctx, cancel := context.WithTimeout(context.Background(), tempoMaxEspera)
	defer cancel() // Boa prática: libera os recursos do contexto ao terminar a função

	// Tenta adquirir a vaga respeitando o limite de tempo estabelecido no contexto
	err := sem.Acquire(ctx, 1)
	if err != nil {
		// Se retornar erro (context.DeadlineExceeded), significa que o tempo acabou e a pessoa desistiu
		fmt.Printf("[Pessoa %2d] XXXXX DESISTIU de esperar após %v e foi embora\n", id, tempoMaxEspera)
		return
	}

	// Se não deu erro, conseguiu adquirir a vaga com sucesso
	fmt.Printf("[Pessoa %2d] >>> ENTROU no banheiro (esperou menos de %v)\n", id, tempoMaxEspera)
	duracao := time.Duration(1+rand.Intn(3)) * time.Second
	time.Sleep(duracao)
	fmt.Printf("[Pessoa %2d] <<< SAIU do banheiro (usou %v)\n", id, duracao)

	// Libera a vaga no semáforo para a próxima pessoa da fila
	sem.Release(1)
}

func main() {
	// Inicializa o gerador de números aleatórios para alternar os testes a cada execução
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup

	for i := 1; i <= totalPessoas; i++ {
		wg.Add(1)
		go usarBanheiro(i, &wg)
	}

	wg.Wait()
	fmt.Println("\nSimulação encerrada.")
}
