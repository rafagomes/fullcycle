# Rate Limiter em Go

Uma aplicação de rate limiter em Go que controla o número de requisições por segundo
com base em endereço IP ou token de acesso. Este projeto utiliza Redis para persistência e foi desenvolvido seguindo as melhores práticas de desenvolvimento.

## Sumário
- [Visão Geral](#visão-geral)
- [Recursos](#recursos)
- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Instalação e Configuração](#instalação-e-configuração)
- [Uso](#uso)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Testes](#testes)
- [Contribuição](#contribuição)
- [Licença](#licença)
- [Contato](#contato)

## Visão Geral
Este projeto implementa um middleware de rate limiting para aplicações web em Go. A aplicação controla o tráfego de requisições utilizando dois critérios:

- **Endereço IP:** Limita o número de requisições de um mesmo IP dentro de um intervalo de tempo.
- **Token de Acesso:** Limita as requisições baseadas em um token fornecido via header `API_KEY`.

As configurações de limite e bloqueio são definidas através de variáveis de ambiente ou arquivo `.env`. A persistência dos dados (como contadores e bloqueios) é realizada via Redis, porém, a implementação utiliza uma interface que permite trocar o mecanismo de armazenamento facilmente.

## Recursos
- Middleware para rate limiting em aplicações web.
- Suporte a limitação por endereço IP e token de acesso.
- Configuração flexível via variáveis de ambiente ou arquivo `.env`.
- Persistência utilizando Redis, com estratégia de abstração para troca de mecanismo.
- Testes automatizados para garantir a robustez da aplicação.

## Tecnologias Utilizadas
- **Golang:** Linguagem de programação utilizada para o desenvolvimento.
- **Redis:** Banco de dados in-memory para persistência dos dados do rate limiter.
- **Docker & Docker Compose:** Para configuração e execução do ambiente Redis.
- **godotenv:** Gerenciamento de variáveis de ambiente.
- **go-redis:** Cliente Redis para Go.
- **Chat GPT:** Geração parcial de documentaçao.

## Instalação e Configuração

### Pré-requisitos
- [Go](https://golang.org/) instalado (versão 1.18 ou superior).
- [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/) instalados.

### Passos para Instalação
1. **Clone o repositório:**
   ```bash
   git clone https://github.com/seu-usuario/seu-repositorio.git
   cd seu-repositorio
   ```

2. **Configure as variáveis de ambiente:**
   - Crie um arquivo `.env` na raiz do projeto com o seguinte conteúdo:
     ```dotenv
     RATE_LIMIT_IP=5
     RATE_LIMIT_TOKEN=10
     BLOCK_TIME_IP=300
     BLOCK_TIME_TOKEN=300
     REDIS_ADDR=localhost:6379
     REDIS_PASSWORD=
     REDIS_DB=0
     ```

3. **Suba o ambiente Redis utilizando Docker Compose:**
   ```bash
   docker-compose up -d
   ```

4. **Instale as dependências do projeto:**
   ```bash
   go mod tidy
   ```

## Uso

Para iniciar o servidor, execute o seguinte comando:
```bash
go run cmd/main.go
```
O servidor será iniciado na porta **8080**. Você pode testar a aplicação utilizando um navegador, `curl` ou ferramentas como Postman.

### Exemplos de Requisições

- **Requisição via IP (sem token):**
   ```bash
   curl http://localhost:8080/
   ```

- **Requisição com Token de Acesso:**
   ```bash
   curl -H "API_KEY: seu-token" http://localhost:8080/
   ```

Se o limite de requisições for excedido, a resposta retornará o status **429** com a mensagem:
`"you have reached the maximum number of requests or actions allowed within a certain time frame"`.

## Estrutura do Projeto

```plaintext
project-root/
├── .env
├── docker-compose.yml
├── go.mod
├── cmd/
│   └── main.go
└── internal/
    ├── config/
    │   └── config.go
    ├── limiter/
    │   └── limiter.go
    └── store/
        ├── store.go
        └── redis_store.go
```

- **cmd/main.go:** Ponto de entrada da aplicação.
- **internal/config:** Responsável pela leitura e gerenciamento das configurações.
- **internal/limiter:** Contém a lógica de rate limiting, separada do middleware.
- **internal/store:** Define a interface para persistência e sua implementação utilizando Redis.

## Testes

Os testes automatizados garantem que a lógica do rate limiter funcione corretamente em diferentes cenários (IP e token). Para rodar os testes, execute:
```bash
go test ./... -v
```

## Casos de Uso e Exemplos de Chamadas para Validação

A seguir, você encontrará exemplos de chamadas e casos de uso para validar que o rate limiter está funcionando corretamente tanto para requisições baseadas em IP quanto para token de acesso.

### Caso 1: Requisições via IP (sem Token)

Este caso testa o rate limiter para requisições originadas de um mesmo IP. 
As primeiras 5 requisições (dentro de 1 segundo) devem retornar **HTTP 200**, e a 6ª deve retornar **HTTP 429**.

**Como testar:**  
Abra um terminal e execute o seguinte comando para realizar 6 requisições consecutivas:

```bash
for i in {1..12}; do
  curl -s -o /dev/null -w "Requisição $i: %{http_code}\n" http://localhost:8080/
done
```

**Resultado Esperado:**  
- Requisições 1 a 5: Código **200**
- Requisição 6: Código **429**

### Caso 2: Requisições via Token

Este caso valida o rate limiter utilizando um token de acesso fornecido no header `API_KEY`.
Configuramos um limite de 10 requisições por segundo para tokens.
As primeiras 10 requisições (dentro de 1 segundo) devem retornar **HTTP 200**, e a 11ª deve retornar **HTTP 429**.

**Como testar:**  
Execute o seguinte comando para realizar 11 requisições consecutivas com o token `abc123`:

```bash
for i in {1..11}; do
  curl -s -o /dev/null -w "Requisição $i: %{http_code}\n" -H "API_KEY: abc123" http://localhost:8080/
done
```

**Resultado Esperado:**  
- Requisições 1 a 10: Código **200**
- Requisição 11: Código **429**

### Validação do Bloqueio e Tempo de Expiração

Após exceder o limite, o identificador (IP ou token) ficará bloqueado pelo tempo configurado (ex.: 300 segundos).
Durante esse período, todas as requisições retornarão **HTTP 429**.

**Teste de Bloqueio:**  
1. Execute uma sequência que exceda o limite (como mostrado nos casos acima).
2. Tente realizar novas requisições imediatamente após exceder o limite; elas deverão retornar **429**.
3. Aguarde o período de bloqueio (por exemplo, 300 segundos) e realize uma nova requisição; ela deverá retornar **200**.

### Exemplo de Script de Teste Automatizado

Esses testes estão disponíveis em um script de shell que executa os casos de uso descritos acima na pasta `tests/test_rate_limiter.sh.`

**Como Executar o Script:**  
1. Dê permissão de execução:
```bash
chmod +x test_rate_limiter.sh
```

2. Execute o script:
```bash
./test_rate_limiter.sh
```

### Observações Finais
- Utilize ferramentas como `watch` para monitorar as respostas em tempo real durante testes de alta frequência.
- Certifique-se de que o servidor está rodando na porta **8080** e que o Redis está ativo (via Docker Compose ou de outra forma).
- Ajuste os limites e tempos de bloqueio no arquivo `.env` conforme necessário para diferentes cenários de carga.

Estes casos de uso e chamadas ajudarão a validar que o rate limiter está funcionando conforme o esperado.

