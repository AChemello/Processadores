# T1 FPPD: Problemas Clássicos de Concorrência

Integrantes: Arthur Ligabue Chemello (25100290) e Fernando Flores Merten (25100908)
Linguagem: Go (testado com a versão X.Y, saída de `go version`).
Sem dependências externas.

## Estrutura
- `filosofos/`: Problema 1, Jantar dos Filósofos
- `produtor-consumidor/`: Problema 2, Produtor/Consumidor com buffer limitado

## Problema 1: Filósofos
Parâmetros: `-n` (filósofos, padrão 5), `-r` (refeições por filósofo),
`-modo` (`base` | `hierarquia` | `lugares`).

    go run -race ./filosofos -n 5 -r 10 -modo base        # deadlock (esperado)
    go run -race ./filosofos -n 5 -r 20 -modo hierarquia
    go run -race ./filosofos -n 5 -r 20 -modo lugares

## Problema 2: Produtor/Consumidor
Parâmetros: `-p` (produtores), `-c` (consumidores), `-k` (capacidade),
`-n` (itens por produtor), `-modo` (`canal` | `semaforo`),
`-maxprod`/`-maxcons` (tempo máx. simulado por item), `-timeout` (consumidor 0).

    go run -race ./produtor-consumidor -modo canal -k 10
    go run -race ./produtor-consumidor -modo semaforo -k 10
    # teste do consumo com timeout:
    go run -race ./produtor-consumidor -modo canal -k 10 -n 40 -maxprod 100ms -timeout 10ms

## Coleta de dados (sem -race, que deixa a execução mais lenta)
    go build -o pc ./produtor-consumidor
    for m in canal semaforo; do for k in 1 10 100; do ./pc -modo $m -k $k | grep CSV; done; done
    # sem carga simulada:
    ./pc -modo canal -k 100 -n 50000 -maxprod 0 -maxcons 0