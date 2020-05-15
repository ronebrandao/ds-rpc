# gRPC em GO

## Projeto

O projeto é divido em 3 diretórios:
- client
- proto
- server

### client
É responsável por coletar as informações e enviá-las ao servidor

### proto
Define o `service` em `Protocol Buffers` que será utilizado na comunicação entre cliente e servidor.

### server
Obtém as informações enviadas pelo cliente e salva em um banco de dados `postgres`.

A tabela do banco de dados se chama `performance_info` e possui os seguintes campos:

| campo            | tipo   |
|------------------|--------|
| id               | serial |
| cpu              | double |
| memory_used      | double |
| memory_avaliable | double |
| disk_used        | double |
| disk_avaliable   | double |

## Execução
* navegue para a pasta do projeto no `terminal
* execute `cd server` e logo depois `go run main.go`  
* abra um novo terminal, execute `cd client` e logo depois `go run main.go`  

Uma vez em execução, o cliente irá automaticamente enviar informações sobre a máquina a cada dez segundos. 
