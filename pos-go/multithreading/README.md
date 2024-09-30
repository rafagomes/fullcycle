# Projeto de Teste de Multithreading em Go

Este projeto tem como objetivo testar o uso de multithreading em Go ao realizar buscas de CEP utilizando duas APIs diferentes simultaneamente: a BrasilAPI e a ViaCep.

## Requisitos

- Go instalado (versão 1.19 ou superior recomendada)

## Passos para rodar

1. Clone o repositório ou copie o código.
2. No terminal, navegue até o diretório do código-fonte.
3. Execute o comando:

   ```bash
   go run main.go
   ```

## Como o projeto demonstra multithreading

O código faz uso de goroutines, que são threads leves no Go, para realizar chamadas simultâneas às duas APIs:

- `go getFromBrasilAPI(cep, channel1)`: Esta linha inicia uma goroutine para buscar o CEP na BrasilAPI de forma assíncrona.
- `go getFromViaCep(cep, channel2)`: Esta linha inicia outra goroutine para buscar o mesmo CEP na ViaCep, também de forma assíncrona.

O uso de canais (`chan Address`) permite que as goroutines enviem os resultados de volta para a função `main`.

O `select` dentro da função `main` escuta ambos os canais ao mesmo tempo e escolhe a primeira resposta que chega ou aguarda um timeout de 1 segundo:

- Se `channel1` responder primeiro, imprime a resposta da BrasilAPI.
- Se `channel2` responder primeiro, imprime a resposta da ViaCep.
- Se nenhuma das APIs responder em 1 segundo, imprime "timeout".

Esta abordagem demonstra o uso eficiente de multithreading, onde múltiplas operações são realizadas simultaneamente, aproveitando o potencial de concorrência do Go.

## Testando casos extremos (Edge Cases)

Para testar cenários onde uma API demora para responder, você pode descomentar uma ou ambas as linhas `time.Sleep(time.Second * 2)` nas funções `getFromBrasilAPI` e `getFromViaCep`. Isso permitirá que você simule situações de atraso e veja como o programa reage nesses casos de multithreading e gerenciamento de tempo de resposta.

## Resultado esperado

- Se uma das APIs responder em menos de 1 segundo, você verá a URL e o JSON da resposta no terminal.
- Se nenhuma das APIs responder em 1 segundo, o programa exibirá `timeout`.
