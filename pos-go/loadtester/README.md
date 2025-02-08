# LoadTester CLI

Projeto para execução de testes de carga em um serviço web, desenvolvido em Go.

## Visão Geral

Esta aplicação é uma ferramenta CLI para realizar testes de carga em um serviço web. 
Os parâmetros necessários são:
- `--url`: URL do serviço a ser testado.
- `--requests`: Número total de requisições a serem enviadas.
- `--concurrency`: Quantidade de chamadas simultâneas.

Ao final da execução, o sistema gera um relatório contendo:
- Tempo total de execução.
- Total de requests realizados.
- Quantidade de requests com status HTTP 200.
- Distribuição dos demais códigos de status HTTP (e.g., 404, 500) e erros ocorridos.

## Estrutura do Projeto

O projeto foi organizado seguindo os princípios da Clean Architecture, dividindo as responsabilidades em camadas. 
A estrutura de diretórios é a seguinte:

```
loadtester/
├── cmd/
│   └── loadtester/
│       └── main.go        // Camada de entrada (CLI), responsável por interpretar os argumentos e iniciar a execução.
├── internal/
│   ├── domain/
│   │   └── model.go       // Define as entidades do domínio, como Result e Report.
│   └── usecase/
│       └── loadtester.go  // Implementa a lógica do teste de carga, orquestrando a execução dos requests.
├── go.mod                 // Gerenciamento de dependências.
└── Dockerfile             // Dockerfile multi-stage, otimizado para criar um container leve com a imagem `scratch`.
```

## Detalhes do Codebase

### Entrada de Parâmetros
- Implementada em `cmd/loadtester/main.go`, a aplicação utiliza o pacote `flag` para ler os parâmetros CLI.
- São validados os parâmetros `--url`, `--requests` e `--concurrency` para garantir a execução correta.

### Lógica de Execução (Caso de Uso)
- Em `internal/usecase/loadtester.go`, a lógica de execução dos testes de carga é centralizada no interactor `LoadTester`.
- O sistema utiliza workers concorrentes para enviar os requests HTTP, distribuindo-os conforme o nível de concorrência definido.
- Os resultados dos requests (status e erros) são coletados e consolidados em um relatório.

### Domínio
- As entidades `Result` e `Report` estão definidas em `internal/domain/model.go`.
- `Result` armazena o resultado de cada request, enquanto `Report` agrega os dados finais do teste.

### Docker e Otimizações
- O `Dockerfile` utiliza build multi-stage para compilar o binário em uma imagem base `golang:1.23-alpine`
  e depois copiar o executável para uma imagem `scratch`, garantindo um container com o menor footprint possível.
- O binário é compilado de forma estática (`CGO_ENABLED=0`) e com flags de otimização (`-ldflags="-s -w"`), resultando em um executável leve.

## Como Utilizar

### Executando com Docker
1. **Build da Imagem Docker:**
   ```bash
   docker build -t loadtester .
   ```

2. **Executando o Container:**
   ```bash
   docker run --rm loadtester --url=http://google.com --requests=1000 --concurrency=10
   ```

## Conclusão

Esta aplicação foi desenvolvida para atender à necessidade de realizar testes de carga em serviços web, 
com foco na modularidade, performance e facilidade de manutenção. 
A separação em camadas, o uso de Docker para ambientes isolados e as otimizações realizadas garantem 
que o sistema seja escalável e adequado para testes em ambientes de produção.

