package main

import (
	"fmt"
	"time"
)

func main() {

	ch := make(chan string, 0) 

	done := make(chan bool) 

	// Goroutine Produtora
	go func() {
		for i := 1; i <= 5; i++ {
			msg := fmt.Sprintf("Mensagem %d", i)
			fmt.Printf("[%s] Produtor: Enviando %s...\n", time.Now().Format("15:04:05"), msg)
			ch <- msg 
			time.Sleep(20 * time.Millisecond) 
		}
		fmt.Printf("[%s] Produtor: Terminou e fechou o canal.\n", time.Now().Format("15:04:05"))
		close(ch)
	}()

	// Goroutine Consumidora
	go func() {
		for msg := range ch {
			time.Sleep(100 * time.Millisecond) 
			fmt.Printf("[%s] Consumidor: Recebeu %s\n", time.Now().Format("15:04:05"), msg)
		}
		done <- true 
	}()

	<-done
	fmt.Println("Fim!")
}


// Parte A — Canal sem buffer (make(chan string))
// a) Observação do comportamento:O produtor avisa que 
// vai enviar a mensagem e o programa sofre uma pausa perceptível de 1 segundo antes de o consumidor 
// exibir que a recebeu. O produtor fica completamente travado a cada iteração, esperando o consumidor 
// terminar sua própria lentidão para poder entregar o dado.

// b) Quem dita o ritmo da comunicação, 
// o produtor ou o consumidor? Por quê?O consumidor. 
// Como o canal não possui buffer, a comunicação é estritamente síncrona (rendez-vous). 
// Isso significa que uma goroutine que envia dados obrigatoriamente bloqueia até que outra 
// goroutine esteja pronta para recebê-los simultaneamente.
//  Como o consumidor possui um atraso de 1 segundo (time.Sleep), 
//  o produtor (mesmo querendo enviar a cada 200ms) é forçado a esperar o ritmo lento do consumidor.
// ---------------------------------------------------------------------------------------------------------------------
//  Parte B — Canal com buffer
//  c) Observação do comportamento (Buffer de tamanho 5):O produtor 
//  dispara todas as suas mensagens de texto na tela quase instantaneamente 
//  (respeitando apenas o seu intervalo curto de 200ms) e fecha o canal muito rápido. O consumidor, 
//  por sua vez, fica processando as mensagens remanescentes que ficaram guardadas na fila, imprimindo 
//  uma por uma a cada 1 segundo, mesmo com o produtor já finalizado.

//  d) O comportamento do produtor 
//  mudou? Por quê?Sim, mudou radicalmente. Com um buffer de tamanho 5, 
//  a comunicação passou a ser assíncrona. Como o número de mensagens enviadas (5) 
//  é exatamente igual à capacidade de armazenamento temporário do canal, 
//  o produtor conseguiu "despejar" todas as strings na fila sem precisar bloquear em nenhuma delas. 
//  Ele se libertou do ritmo do consumidor por haver vagas disponíveis no canal.
//  e) O que acontece quando o buffer enche? (Buffer de tamanho 2):O produtor consegue enviar 
//  as duas primeiras mensagens de forma assíncrona e imediata. No entanto, 
//  ao tentar realizar o envio da terceira mensagem, ele encontra o buffer completamente cheio. 
//  Nesse exato momento, o produtor bloqueia (trava) e o canal volta a agir temporariamente 
//  de forma síncrona. O produtor só é destravado quando o consumidor executa a 
//  leitura da primeira mensagem, esvaziando uma vaga na fila e permitindo que o 
//  fluxo prossiga em blocos.