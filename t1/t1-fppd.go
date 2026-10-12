// T1 (FPPD) — Problema 1: Jantar dos Filósofos
//
// Representação dos garfos (como em aula10/03-canais/filosofos):
// cada garfo é um canal bufferizado de capacidade 1 que guarda um token.
//
//	pegar o garfo    = receber do canal  (<-garfos[i])    -> wait / P
//	devolver o garfo = enviar ao canal   (garfos[i] <- {}) -> signal / V
//
// Se o garfo está com outro filósofo, o canal está vazio e quem tenta
// receber bloqueia: exclusão mútua + espera, tudo no próprio canal.
//
// Uso:
//
//	go run -race ./filosofos -n 5 -r 10 -modo base
//
// Modos (um por passo do trabalho):
//
//	base  -> todos pegam primeiro o garfo da esquerda (DEADLOCK)
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Tempos simulados. São constantes (iguais em todos os modos) para que a
// comparação de fairness entre as estratégias seja justa.
const (
	maxPensar = 20 * time.Millisecond
	maxComer  = 20 * time.Millisecond

	// Pausa entre pegar o 1º e o 2º garfo. Como é maior que maxPensar,
	// todos os filósofos conseguem pegar o 1º garfo antes de qualquer um
	// tentar o 2º: isso torna o deadlock da versão base determinístico
	// na prática (o mesmo papel do time.Sleep nos slides de deadlock).
	pausaEntreGarfos = 50 * time.Millisecond
)

// Config reúne os parâmetros recebidos pela linha de comando.
type Config struct {
	N       int           // número de filósofos
	R       int           // refeições (iterações) por filósofo
	Timeout time.Duration // tempo sem terminar => declara deadlock (modo base)
}

// Stats guarda os dados de execução. Tudo é protegido por um Mutex
// (aula 05/09), então pode ser lido pela main mesmo com filósofos ativos
// ou bloqueados, sem data race.
type Stats struct {
	mu        sync.Mutex
	refeicoes []int
	espera    []time.Duration
	esperaMax []time.Duration
	estado    []string
}

func NovoStats(n int) *Stats {
	s := &Stats{
		refeicoes: make([]int, n),
		espera:    make([]time.Duration, n),
		esperaMax: make([]time.Duration, n),
		estado:    make([]string, n),
	}
	for i := range s.estado {
		s.estado[i] = "pensando"
	}
	return s
}

func (s *Stats) SetEstado(id int, e string) {
	s.mu.Lock()
	s.estado[id] = e
	s.mu.Unlock()
}

func (s *Stats) RegistrarRefeicao(id int, espera time.Duration) {
	s.mu.Lock()
	s.refeicoes[id]++
	s.espera[id] += espera
	if espera > s.esperaMax[id] {
		s.esperaMax[id] = espera
	}
	s.mu.Unlock()
}

func (s *Stats) ImprimirResumo() {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println("\nFilósofo | Refeições | Espera média | Espera máxima")
	for i := range s.refeicoes {
		var media time.Duration
		if s.refeicoes[i] > 0 {
			media = s.espera[i] / time.Duration(s.refeicoes[i])
		}
		fmt.Printf("%8d | %9d | %12v | %v\n", i, s.refeicoes[i], media, s.esperaMax[i])
	}
}

// dorme simula pensar/comer por um tempo aleatório em [0, limite).
func dorme(limite time.Duration) {
	time.Sleep(time.Duration(rand.Int63n(int64(limite))))
}

// novosGarfos cria N canais bufferizados (cap. 1), cada um já com seu
// token: todos os garfos começam livres.
func novosGarfos(n int) []chan struct{} {
	garfos := make([]chan struct{}, n)
	for i := range garfos {
		garfos[i] = make(chan struct{}, 1)
		garfos[i] <- struct{}{}
	}
	return garfos
}

// ---------------------------------------------------------------------
// VERSÃO BASE: todos pegam o garfo da ESQUERDA primeiro  => deadlock.
//
// Filósofo i usa o garfo i (esquerdo) e o garfo (i+1)%N (direito).
// Todos seguem a mesma ordem: esquerdo, depois direito. É exatamente a
// simetria do Algoritmo 6.10 / 8.5 do Ben-Ari sem correção.
// ---------------------------------------------------------------------
func filosofoBase(id int, cfg Config, garfos []chan struct{}, st *Stats, wg *sync.WaitGroup) {
	defer wg.Done()
	esq := id
	dir := (id + 1) % cfg.N

	for r := 0; r < cfg.R; r++ {
		dorme(maxPensar) // pensar (seção não-crítica)

		inicio := time.Now()
		st.SetEstado(id, fmt.Sprintf("esperando garfo esquerdo (%d)", esq))
		<-garfos[esq] // pega o esquerdo
		st.SetEstado(id, fmt.Sprintf("segura garfo %d, esperando garfo direito (%d)", esq, dir))

		time.Sleep(pausaEntreGarfos) // abre a janela para a espera circular

		<-garfos[dir] // pega o direito  <- aqui o sistema trava
		espera := time.Since(inicio)

		st.SetEstado(id, "comendo")
		dorme(maxComer) // comer (seção crítica)
		st.RegistrarRefeicao(id, espera)

		garfos[esq] <- struct{}{} // devolve os garfos
		garfos[dir] <- struct{}{}
		st.SetEstado(id, "pensando")
	}
}

// relatarDeadlock imprime o estado de cada filósofo no momento em que o
// watchdog dispara, e a ligação com as condições de Coffman.
func relatarDeadlock(cfg Config, st *Stats) {
	st.mu.Lock()
	defer st.mu.Unlock()
	fmt.Printf("\n*** DEADLOCK: nenhum filósofo concluiu em %v ***\n", cfg.Timeout)
	for i, e := range st.estado {
		fmt.Printf("  Filósofo %d: %s\n", i, e)
	}
	fmt.Println("\nCondições de Coffman presentes:")
	fmt.Println("  1. Exclusão mútua  : cada garfo só pode estar com um filósofo")
	fmt.Println("  2. Hold and wait   : cada um segura o esquerdo e pede o direito")
	fmt.Println("  3. Não preempção   : ninguém tira o garfo de outro à força")
	fmt.Println("  4. Espera circular : 0 espera 1, 1 espera 2, ..., N-1 espera 0")
}

// executarBase roda a versão base. Como ela DEVE travar, usamos um
// watchdog com select + time.After (aula 09): se a main não receber o
// sinal de término a tempo, declaramos deadlock.
//
// Observação: o runtime de Go só aborta com "all goroutines are asleep"
// quando NÃO há timer pendente (aula 04). O time.After do watchdog é um
// timer pendente, então quem detecta o deadlock aqui somos nós.
func executarBase(cfg Config) {
	garfos := novosGarfos(cfg.N)
	st := NovoStats(cfg.N)

	var wg sync.WaitGroup
	for i := 0; i < cfg.N; i++ {
		wg.Add(1)
		go filosofoBase(i, cfg, garfos, st, &wg)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Todos terminaram (o deadlock NÃO ocorreu nesta execução).")
		st.ImprimirResumo()
	case <-time.After(cfg.Timeout):
		relatarDeadlock(cfg, st)
		st.ImprimirResumo()
		os.Exit(1) // as goroutines continuam bloqueadas: é o deadlock
	}
}

func filosofoHierarquia(id int, cfg Config, garfos []chan struct{}, st *Stats, wg *sync.WaitGroup) {
	defer wg.Done()
	primeiro, segundo := id, (id+1)%cfg.N
	if primeiro > segundo {
		primeiro, segundo = segundo, primeiro
	}

	for r := 0; r < cfg.R; r++ {
		dorme(maxPensar)

		inicio := time.Now()
		st.SetEstado(id, fmt.Sprintf("esperando garfo %d", primeiro))
		<-garfos[primeiro]
		time.Sleep(pausaEntreGarfos)
		st.SetEstado(id, fmt.Sprintf("segura %d, esperando %d", primeiro, segundo))
		<-garfos[segundo]
		espera := time.Since(inicio)

		st.SetEstado(id, "comendo")
		dorme(maxComer)
		st.RegistrarRefeicao(id, espera)

		garfos[primeiro] <- struct{}{}
		garfos[segundo] <- struct{}{}
		st.SetEstado(id, "pensando")
	}
}

func filosofoLugares(id int, cfg Config, garfos []chan struct{}, seats chan struct{}, st *Stats, wg *sync.WaitGroup) {
	defer wg.Done()
	esq := id
	dir := (id + 1) % cfg.N

	for r := 0; r < cfg.R; r++ {
		dorme(maxPensar)

		inicio := time.Now()
		st.SetEstado(id, "esperando lugar na mesa")
		<-seats // wait(seats)
		st.SetEstado(id, fmt.Sprintf("esperando garfo esquerdo (%d)", esq))
		<-garfos[esq]
		time.Sleep(pausaEntreGarfos)
		st.SetEstado(id, fmt.Sprintf("segura %d, esperando %d", esq, dir))
		<-garfos[dir]
		espera := time.Since(inicio)

		st.SetEstado(id, "comendo")
		dorme(maxComer)
		st.RegistrarRefeicao(id, espera)

		garfos[esq] <- struct{}{}
		garfos[dir] <- struct{}{}
		seats <- struct{}{} // signal(seats)
		st.SetEstado(id, "pensando")
	}
}

func executar(cfg Config, modo string) {
	garfos := novosGarfos(cfg.N)
	st := NovoStats(cfg.N)

	var seats chan struct{}
	if modo == "lugares" {
		seats = make(chan struct{}, cfg.N-1)
		for i := 0; i < cfg.N-1; i++ {
			seats <- struct{}{} // valor inicial do semáforo = N-1
		}
	}

	inicio := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < cfg.N; i++ {
		wg.Add(1)
		switch modo {
		case "hierarquia":
			go filosofoHierarquia(i, cfg, garfos, st, &wg)
		case "lugares":
			go filosofoLugares(i, cfg, garfos, seats, st, &wg)
		}
	}
	wg.Wait()

	fmt.Printf("Todos terminaram em %v, sem goroutines bloqueadas.\n", time.Since(inicio))
	st.ImprimirResumo()
}

func main() {
	n := flag.Int("n", 5, "número de filósofos (>= 2)")
	r := flag.Int("r", 10, "refeições (iterações) por filósofo")
	modo := flag.String("modo", "base", "estratégia: base")
	timeout := flag.Duration("timeout", 2*time.Second, "tempo para declarar deadlock (modo base)")
	flag.Parse()

	if *n < 2 || *r < 1 {
		fmt.Fprintln(os.Stderr, "uso: go run ./filosofos -n N (N>=2) -r R (R>=1) -modo base")
		os.Exit(2)
	}
	cfg := Config{N: *n, R: *r, Timeout: *timeout}

	fmt.Printf("=== Jantar dos Filósofos | modo=%s N=%d R=%d ===\n", *modo, cfg.N, cfg.R)

	switch *modo {
	case "base":
		executarBase(cfg)
	case "hierarquia", "lugares":
		executar(cfg, *modo)
	default:
		fmt.Fprintf(os.Stderr, "modo %q ainda não implementado\n", *modo)
		os.Exit(2)
	}
}
// falta fazer o relatorio e estudar o codigo para apresentar (fazer um relatorio - explicando linha por linha)