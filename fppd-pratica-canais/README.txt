FPPD — Prática: Channels em Go
===============================

Este pacote contém os códigos-fonte de referência e pontos de partida
para o roteiro de prática sobre canais em Go.

Estrutura:
  exemplos/     Exemplos de referência da aula (canais em Go)
  exercicios/   Códigos-base para os exercícios do roteiro

Exemplos (referência):
  01-canal-sincrono.go        Canal sem buffer — rendez-vous
  02-canal-bufferizado.go     Canal com buffer — envio assíncrono e FIFO
  03-close-range.go           Fechamento de canal e iteração com range
  04-produtor-consumidor.go   Produtor-consumidor com duas goroutines

Exercícios (roteiro):
  ex1-original.go          Exercício 1 — ponto de partida (WaitGroup)
  ex3-base.go              Exercício 3 — estrutura do pipeline

  (ex1.go, ex2.go e ex3.go devem ser criados pelo aluno)

Como executar:
  go run exemplos/01-canal-sincrono.go
  go run exemplos/02-canal-bufferizado.go
  go run exercicios/ex1-original.go

Requisitos:
  Go 1.21 ou superior — https://go.dev/dl
