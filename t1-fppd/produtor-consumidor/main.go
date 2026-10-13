// T1 (FPPD) — Problema 2: Produtor/Consumidor com buffer limitado
//
// Uso:
//   go run -race ./produtor-consumidor -p 4 -c 4 -k 10 -n 500 -modo canal
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Config struct {
	P, C, K, Itens int           
	MaxProduzir    time.Duration 
	MaxConsumir    time.Duration 
	Timeout        time.Duration 
}

type Resultado int

const (
	RItem Resultado = iota
	RTimeout
	REncerrado
)

type Buffer interface {
	Put(item int)                           
	Take(timeout time.Duration) (int, Resultado) 
	Fechar()                                
	Tamanho() int                           
}

type ResConsumidor struct {
	itens    []int
	timeouts int
}

type BufferCanal struct{ ch chan int }

func NovoBufferCanal(k int) *BufferCanal {
	return &BufferCanal{ch: make(chan int, k)}
}

func (b *BufferCanal) Put(item int) { b.ch <- item }
func (b *BufferCanal) Fechar()      { close(b.ch) }
func (b *BufferCanal) Tamanho() int { return len(b.ch) }

func (b *BufferCanal) Take(timeout time.Duration) (int, Resultado) {
	if timeout <= 0 {
		item, ok := <-b.ch 
		if !ok {
			return 0, REncerrado
		}
		return item, RItem
	}
	select {
	case item, ok := <-b.ch:
		if !ok {
			return 0, REncerrado
		}
		return item, RItem
	case <-time.After(timeout):
		return 0, RTimeout
	}
}

func dorme(limite time.Duration) {
	if limite <= 0 {
		return
	}
	time.Sleep(time.Duration(rand.Int63n(int64(limite))))
}

func produtor(id int, cfg Config, buf Buffer, produzidos *int, wg *sync.WaitGroup) {
	defer wg.Done()
	for s := 0; s < cfg.Itens; s++ {
		dorme(cfg.MaxProduzir)
		buf.Put(id*cfg.Itens + s)
		*produzidos++
	}
}

func consumidor(id int, cfg Config, buf Buffer, res *ResConsumidor, wg *sync.WaitGroup) {
	defer wg.Done()
	var timeout time.Duration
	if id == 0 {
		timeout = cfg.Timeout 
	}
	for {
		item, r := buf.Take(timeout)
		if r == REncerrado {
			return
		}
		if r == RTimeout {
			res.timeouts++ 
			continue
		}
		res.itens = append(res.itens, item)
		dorme(cfg.MaxConsumir)
	}
}

func medirOcupacao(tamanho func() int, parar <-chan struct{}) <-chan float64 {
	res := make(chan float64, 1)
	go func() {
		soma, n := 0, 0
		for {
			select {
			case <-parar:
				if n == 0 {
					res <- 0
				} else {
					res <- float64(soma) / float64(n)
				}
				return
			case <-time.After(time.Millisecond):
				soma += tamanho()
				n++
			}
		}
	}()
	return res
}

func executar(cfg Config, modo string, buf Buffer) {
	total := cfg.P * cfg.Itens
	produzidos := make([]int, cfg.P)
	res := make([]ResConsumidor, cfg.C)

	parar := make(chan struct{})
	ocup := medirOcupacao(buf.Tamanho, parar)

	inicio := time.Now()
	var wgP, wgC sync.WaitGroup
	for i := 0; i < cfg.C; i++ {
		wgC.Add(1)
		go consumidor(i, cfg, buf, &res[i], &wgC)
	}
	for i := 0; i < cfg.P; i++ {
		wgP.Add(1)
		go produtor(i, cfg, buf, &produzidos[i], &wgP)
	}

	wgP.Wait()
	buf.Fechar()
	wgC.Wait()
	dur := time.Since(inicio)
	close(parar)
	ocupMedia := <-ocup

	prod, cons := 0, 0
	for _, p := range produzidos {
		prod += p
	}
	vistos := make([]int, total)
	for i := range res {
		for _, it := range res[i].itens {
			cons++
			if it >= 0 && it < total {
				vistos[it]++
			}
		}
	}
	ok := prod == total && cons == prod
	for it, v := range vistos {
		if v != 1 {
			ok = false
			fmt.Printf("ERRO: item %d consumido %d vez(es)\n", it, v)
		}
	}

	fmt.Println("\nConsumidor | Itens consumidos | Timeouts")
	for i := range res {
		marca := ""
		if i == 0 {
			marca = " (com timeout)"
		}
		fmt.Printf("%10d | %16d | %8d%s\n", i, len(res[i].itens), res[i].timeouts, marca)
	}
	thr := float64(cons) / dur.Seconds()
	fmt.Printf("\nProduzidos=%d Consumidos=%d | tempo=%v | throughput=%.1f itens/s | ocupação média=%.2f/%d\n",
		prod, cons, dur.Round(time.Millisecond), thr, ocupMedia, cfg.K)
	if ok {
		fmt.Println("VERIFICAÇÃO OK: nenhum item perdido ou duplicado; sem goroutines bloqueadas.")
	} else {
		fmt.Println("VERIFICAÇÃO FALHOU.")
		os.Exit(1)
	}
	fmt.Printf("CSV;%s;%d;%d;%d;%d;%.3f;%.1f;%.2f\n", modo, cfg.P, cfg.C, cfg.K, total, dur.Seconds(), thr, ocupMedia)
}

type Sem chan struct{}

func NovoSem(valorInicial, capacidade int) Sem {
	s := make(Sem, capacidade)
	for i := 0; i < valorInicial; i++ {
		s <- struct{}{} 
	}
	return s
}

func (s Sem) Wait()   { <-s }
func (s Sem) Signal() { s <- struct{}{} }

const sentinela = -1

type BufferSem struct {
	buf      []int
	in, out  int
	notFull  Sem
	notEmpty Sem
	mutexIn  sync.Mutex
	mutexOut sync.Mutex
	c        int 
}

func NovoBufferSem(k, c int) *BufferSem {
	return &BufferSem{
		buf:      make([]int, k),
		notFull:  NovoSem(k, k),
		notEmpty: NovoSem(0, k),
		c:        c,
	}
}

func (b *BufferSem) Put(item int) {
	b.notFull.Wait() 
	b.mutexIn.Lock()
	b.buf[b.in] = item
	b.in = (b.in + 1) % len(b.buf)
	b.mutexIn.Unlock()
	b.notEmpty.Signal()
}

func (b *BufferSem) Take(timeout time.Duration) (int, Resultado) {
	if timeout <= 0 {
		b.notEmpty.Wait() 
	} else {
		select { 
		case <-b.notEmpty:
		case <-time.After(timeout):
			return 0, RTimeout
		}
	}
	b.mutexOut.Lock()
	item := b.buf[b.out]
	b.out = (b.out + 1) % len(b.buf)
	b.mutexOut.Unlock()
	b.notFull.Signal() 

	if item == sentinela {
		return 0, REncerrado
	}
	return item, RItem
}

func (b *BufferSem) Fechar() {
	for i := 0; i < b.c; i++ {
		b.Put(sentinela)
	}
}

func (b *BufferSem) Tamanho() int { return len(b.notEmpty) } 

func main() {
	p := flag.Int("p", 4, "número de produtores")
	c := flag.Int("c", 4, "número de consumidores")
	k := flag.Int("k", 10, "capacidade do buffer")
	n := flag.Int("n", 500, "itens produzidos por produtor")
	maxProd := flag.Duration("maxprod", 2*time.Millisecond, "tempo máx. para produzir um item")
	maxCons := flag.Duration("maxcons", 2*time.Millisecond, "tempo máx. para consumir um item")
	timeout := flag.Duration("timeout", 20*time.Millisecond, "timeout do consumidor 0")
	modo := flag.String("modo", "canal", "versão do buffer: canal | semaforo")
	flag.Parse()

	if *p < 1 || *c < 1 || *k < 1 || *n < 1 {
		fmt.Fprintln(os.Stderr, "p, c, k e n devem ser >= 1")
		os.Exit(2)
	}
	cfg := Config{P: *p, C: *c, K: *k, Itens: *n, MaxProduzir: *maxProd, MaxConsumir: *maxCons, Timeout: *timeout}

	var buf Buffer
	switch *modo {
	case "canal":
		buf = NovoBufferCanal(cfg.K)
	case "semaforo":
		buf = NovoBufferSem(cfg.K, cfg.C)
	default:
		fmt.Fprintf(os.Stderr, "modo %q ainda não implementado\n", *modo)
		os.Exit(2)
	}

	fmt.Printf("=== Produtor/Consumidor | modo=%s P=%d C=%d K=%d itens/produtor=%d ===\n",
		*modo, cfg.P, cfg.C, cfg.K, cfg.Itens)
	executar(cfg, *modo, buf)
}