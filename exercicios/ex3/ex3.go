// Exercício 3 — Pipeline de processamento
//
// Complete as funções gerador, multiplicador e impressora
// para formar um pipeline de 3 estágios conectados por canais.
package main

import "fmt"

// gerador envia os números de 'valores' no canal 'out' e fecha o canal.
func gerador(valores []int, out chan<- int) {
	// Requisito a) Iterando sobre os valores fornecidos
	for _, num := range valores {
		out <- num
	}
	// Requisito b) Fecha o canal ao terminar
	close(out) 
}

// multiplicador lê do canal 'in', multiplica por 2, e envia no canal 'out'.
// Fecha 'out' ao terminar.
func multiplicador(in <-chan int, out chan<- int) {
	// Requisito a) Use range para iterar sobre o canal de entrada
	for num := range in {
		out <- num * 2
	}
	// Requisito b) Fecha o canal ao terminar
	close(out) 
}

// impressora lê do canal 'in' e imprime cada valor.
// Envia true no canal 'done' ao terminar.
func impressora(in <-chan int, done chan<- bool) {
	// Requisito a) Use range para iterar sobre o canal de entrada
	for num := range in {
		fmt.Println(num)
	}
	// Requisito c) Envia true no canal done ao terminar
	done <- true 
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

// Resposta: Os canais deste pipeline podem ser sem buffer. 
// Como cada estágio do pipeline roda em sua própria goroutine de 
// forma concorrente e independente, o mecanismo de sincronismo estrito (rendez-vous) 
// funciona perfeitamente. Assim que o gerador coloca um número no canal, 
// ele bloqueia até que o multiplicador o receba; este, por sua vez, processa 
// o dado e o repassa imediatamente para a impressora. O fluxo de dados 
// ocorre de forma contínua e segura sem a necessidade de armazenamento temporário.
// Nota de Otimização (Opcional): Embora os canais sem buffer atendam 
// perfeitamente à lógica do problema, o uso de canais com buffer (ex: tamanho 2 ou 3) 
// poderia ser adotado como uma otimização de performance (desacoplamento temporal). 
// O buffer permitiria que o gerador continuasse lendo os próximos itens da lista 
// sem ter que esperar o processamento lento e a impressão dos estágios seguintes, 
// melhorando a vazão total se os estágios tivessem tempos de execução muito discrepantes.