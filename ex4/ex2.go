// ex2.go — Exercício 2: Largada Sincronizada com Variáveis de Condição
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const numCorredores = 6

// Declaração das três variáveis globais exigidas pelo roteiro
var (
	mu     sync.Mutex
	cond   = sync.NewCond(&mu)
	pronto = false // Inicialmente false
)

func corredor(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("[Corredor %d] na linha de largada\n", id)

	// Adquire o lock e espera em loop 'for' enquanto 'pronto' for false
	mu.Lock()
	for !pronto {
		cond.Wait() // Libera o lock, dorme, e readquire o lock ao acordar
	}
	mu.Unlock() // Libera o lock ao sair do loop

	// Início da corrida após a liberação do juiz
	duracao := time.Duration(500+rand.Intn(2000)) * time.Millisecond
	fmt.Printf("[Corredor %d] LARGOU!\n", id)
	time.Sleep(duracao)
	fmt.Printf("[Corredor %d] chegou (tempo: %v)\n", id, duracao)
}

func main() {
	// Inicializa o gerador de números aleatórios
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup

	for i := 1; i <= numCorredores; i++ {
		wg.Add(1)
		go corredor(i, &wg)
	}

	// Juiz prepara a pista
	fmt.Println("\n[Juiz] Preparando a pista...")
	time.Sleep(2 * time.Second)

	// Adquire o lock e altera pronto para true
	mu.Lock()
	pronto = true

	// CORREÇÃO DE LOG: Printamos o sinal do Juiz ANTES de acordar os corredores.
	// Como ainda estamos segurando o Lock, nenhuma goroutine consegue avançar
	// além do cond.Wait() até liberarmos o mu.Unlock() abaixo.
	fmt.Println("[Juiz] VAI!")
	fmt.Println() // Linha em branco correta sem quebrar o linter do Go vet

	cond.Broadcast() // Acorda TODOS os corredores ao mesmo tempo
	mu.Unlock()

	wg.Wait()
	fmt.Println("\nCorrida encerrada.")
}
