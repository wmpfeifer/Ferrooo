# Ferrooo

## Linguagem:
- Go (1.24.5 amd64)

## Bibliotecas: 
- SonicJson;
- Fast HTTP.

## Framework:
- Fiber.

## IDEs:
- Jetbrains Goland

## Banco de Dados:
- PostgreSQL (Com índices, versão = 15 Alpine)

## Conteiner:
### Docker.
Para rodar o Docker é preciso primeiro entrar na pasta do projeto pelo terminal e executar o comando:
- ```docker-compose up --build``` esse comnando irá upar os containers do banco de dados e da aplicação.

Caso queira derrubar os containers, execute o seguinte comando:
- ```docker-compose down``` esse comando irá derrubar os containers do banco de dados e da aplicação.

Para testar se o docker está rodando corretamente, rode esse comnando:
- ``Invoke-WebRequest -Uri http://localhost:9999/payments-summary`` e veja se o status é 200 OK.

Pode rodar esse também:
- ``Invoke-WebRequest -Uri http://localhost:9999/health`` e veja se o status é 200 OK.


## Para Versionamento:
- Github. 

Endereço Web do Github:
- ``https://github.com/wmpfeifer/Ferrooo.git``

## Integrantes:
- Willian Menegazzo Pfeifer;
- João Vitor Possenti;


