### Languages: [Português 🇧🇷](#Rate-Limiter) | [English 🇨🇦](#rate-limiter-api)

---

# Rate Limiter

Este projeto implementa um rate limiter usando Redis e Golang para controlar a quantidade de requisições que podem ser
feitas por IP ou por token.
Ele é configurável e robusto, permitindo que você defina diferentes limites de requisições baseados no IP ou em tokens
autorizados.
Esta API faz parte do desafio técnicos do curso de Pós-Graduação em Engenharia de
Software [GoExpert](https://goexpert.fullcycle.com.br/pos-goexpert/).

## ⚙️ Configuração

Você precisará das seguintes tecnologias abaixo:

- [Docker](https://docs.docker.com/get-docker/) 🐳
- [Docker Compose](https://docs.docker.com/compose/install/) 🐳
- [Postman ☄️](https://www.postman.com/downloads/) ou [VS Code](https://code.visualstudio.com/download) com a
  extensão [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) instalada.
- [GIT](https://git-scm.com/downloads)

## 🚀 Iniciando

1. Clone o repositório e entre no diretório do projeto.
   ```sh
   git clone https://github.com/brunoliveiradev/Rate-Limiter-Challenge-GoExpertPostGrad.git
   cd Rate-Limiter-Challenge-GoExpertPostGrad
   ```

2. Execute o comando abaixo na pasta raiz do projeto para iniciar o redis:
   ```sh
   docker compose up --build -d
   ```

   Para parar os serviços:
   ```sh
   docker compose down
   ```
3. Execute o servidor na porta 8080:

    ```sh
    go run main.go
    ```

4. A **API REST** estará disponível em `http://localhost:8080` 🚀.

## 🧪 Testes

1. Execute o comando abaixo para rodar os testes:
    ```sh
    go test -v ./...
    ```
2. Caso queira, utilize o arquivo `rate_limit.http` para fazer requisições de teste:
    1. Abra o arquivo `rate_limit.http` no seu editor de texto, localizado no caminho `api/rate_limit.http`.
    2. Se preferir utilizar o cURL, você pode copiar o conteúdo do arquivo `rate_limit.http` e colar no terminal ou
       Postman.

3. Na pasta api, há também dois scripts de shell para testar o rate limiter:
    1. `api/test.sh `
    2. `api/test_rate_limiter.sh`

4. Para executar primeiro de permissões de execução para os scripts:
    ```sh
    chmod +x api/test.sh api/test_rate_limiter.sh
    ```
5. Para executar os scripts de teste, use os comandos:
    ```sh
    ./api/test.sh
    ./api/test_rate_limiter.sh
    ```

---

# Rate Limiter API

This project implements a rate limiter using Redis and Golang to control the number of requests that can be made by IP
or by token.
It is configurable and robust, allowing you to define different request limits based on IP or authorized tokens.
This API is part of the technical challenges of the Postgraduate course in Software
Engineering [GoExpert](https://goexpert.fullcycle.com.br/pos-goexpert/).

## ⚙️ Configuration

You will need the following technologies below:

- [Docker](https://docs.docker.com/get-docker/) 🐳
- [Docker Compose](https://docs.docker.com/compose/install/) 🐳
- [Postman ☄️](https://www.postman.com/downloads/) ou [VS Code](https://code.visualstudio.com/download) com a
  extensão [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) instalada.
- [GIT](https://git-scm.com/downloads)

## 🚀 Getting Started

1. Clone the repository and enter the project directory.
   ```sh
   git clone https://github.com/brunoliveiradev/Rate-Limiter-Challenge-GoExpertPostGrad.git
   cd Rate-Limiter-Challenge-GoExpertPostGrad
   ```
2. Run the command below in the project root folder to start the redis:
   ```sh
   docker compose up --build -d
   ```

   To stop the services:
    ```sh
    docker compose down
    ```
3. Run the server on port 8080:
    ```sh
    go run main.go
    ```
4. The **REST API** will be available at `http://localhost:8080` 🚀.

## 🧪 Tests

1. Run the command below to run the tests:
    ```sh
    go test -v ./...
    ```
2. If you want, use the `rate_limit.http` file to make test requests:
    1. Open the `rate_limit.http` file in your text editor, located at `api/rate_limit.http`.
    2. If you prefer to use cURL, you can copy the contents of the `rate_limit.http` file and paste it into the terminal
       or Postman.
3. In the api folder, there are also two shell scripts to test the rate limiter:
    1. `api/test.sh `
    2. `api/test_rate_limiter.sh`
4. To run first give execution permissions to the scripts:
   ```sh
   chmod +x api/test.sh api/test_rate_limiter.sh
   ```
5. To run the test scripts, use the commands:
    ```sh
    ./api/test.sh
    ./api/test_rate_limiter.sh
    ```
