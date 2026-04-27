# Desafio Prefeitura Rio

Este projeto consiste em uma API REST desenvolvida em Go utilizando o framework Gin, persistência de dados no PostgreSQL e gerenciamento de cache de requisições no Redis.

A aplicação é totalmente containerizada para garantir que o ambiente se mantém em qualquer que seja o local que seja executado.

Na Cidade do Rio de Janeiro, qualquer cidadão consegue criar um chamado por meio de um sistema público da Prefeitura para solicitar reparos de ordem pública como reparos na iluminação pública (postes), buracos nas vias, entre outros.

Esse projeto consiste em um sistema de notificações, responsável por receber atualizações dos chamados via REST API e atualizar os cidadãos em tempo real via WebSocket.

## Requisitos

- [Docker Engine](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Just Task Runner](https://github.com/casey/just#installation) (Opcional)

## Arquitetura em Containers

1. **app_go**: A aplicação Go rodando na porta `8080`.
2. **postgres_db**: Instância do PostgreSQL na porta `5432`.
3. **redis_cache**: Instância do Redis na porta `6379`.

Além dos serviços essenciais, foram acrescentados serviços para facilitar a visualização dos dados no PostgreSQL e Redis por meio de GUIs disponibilizadas pelo **pgadmin** e **redis_insight**.

## Como rodar o projeto

Com o **Docker** e **Docker Compose** instalados, altere as variáveis de ambiente conforme necessário no arquivo `docker-compose.yml` na raiz do projeto e execute o comando: `docker compose up -d`. Tendo o Just instalado, basta executar `just run` para subir todos os contêineres.

## API REST

O projeto conta com as seguintes rotas:

- `GET /health`- Retorna `{ "status": "OK" }` como resultado se a aplicação estiver funcionando.

As requisições abaixo esperam um token JWT como Bearer Token com o CPF do cidadão no campo `preferred_username` para poder identificar quais notificações devem ser lidas/alteradas assinado com o `JWT_SECRET`.
- `GET /notifications` - Retorna todas as notificações criadas:

```
{
	"data": [
		{
			"id": 1,
			"chamado_id": "CH-2024-001235",
			"tipo": "status_change",
			"status_anterior": "finalizada",
			"status_novo": "Removida",
			"titulo": "Buraco na Rua — Atualização",
			"descricao": "Removida",
			"timestamp": "2026-11-15T14:30:00Z",
			"is_read": false
		}
	]
}
```

- `PATCH /notifications/:id/read` - Define a notificação como lida e retorna:
```
{
	"data": ":id",
	"message": "Notification updated successfully"
}
```

- `GET /notifications/unread-count` - Retorna o número de notificações não lidas: `{ "data": 1, "message": "Notification count read successfully" }`

A requisição abaixo exige que o payload venha assinado com o header `X-Signature-256` contendo `sha256=<HMAC-SHA256 do body com o X_SIGNATURE_SECRET>`.

- `POST /` - Rota para criação de novas notificações, exemplo de body esperado:

```
{
  "chamado_id": "CH-2024-001234",
  "tipo": "status_change",
  "cpf": "12345678901",
  "status_anterior": "em_analise",
  "status_novo": "em_execucao",
  "titulo": "Buraco na Rua — Atualização",
  "descricao": "Equipe designada para reparo na Rua das Laranjeiras, 100",
  "timestamp": "2024-11-15T14:30:00Z"
}
```

## WebSocket

Além das rotas REST acima, o projeto conta com uma rota WebSocket que os cidadãos podem conectar para receberem em tempo real as notificações criadas para o seu CPF.

- `GET /ws/:cpf` - Toda vez que uma nova notificação é criada para o CPF passado retorna um body:

```
{ "payload" : 
    { 
    "id": 1,
    "chamado_id" : "CH-2024-001235",
    "tipo" : "status_change",
    "cpf" : "12345678902",
    "status_anterior" : "finalizada",
    "status_novo" : "removida",
    "titulo" : "Buraco na Rua — Atualização",
    "descricao" : "Removida",
    "timestamp" : "2026-11-15T14:30:00Z","is_read":false
    }
}
```

### Redis

O projeto conta com um Redis que armazena temporariamente o header `X-Signature-256` da requisição POST para garantir idempotência, evitando que requisições repetidas insiram uma mesma notificação repetidas vezes no banco de dados.

### GUIs

As interfaces gráficas disponibilizadas pelo PgAdmin e Redis Insight podem ser acessadas em [PgAdmin](http://localhost:8000/login?next=/) e [RedisInsight](http://localhost:5540/), respectivamente, para facilitar a visualização dos dados processados pelo sistema.

### Testes

Para executar os testes basta executar o comando `just test`, ele irá rodar os testes unitários e criará um deploy apartado definido no `docker-compose-test.yml` e `Dockerfile.test` para rodar os testes de integração.

### DLQ

O projeto conta com a existência de uma Dead Letter Queue montada no Postgres ou no Redis, por meio da variável `DLQ_TYPE` com os valores postgres ou redis. Ao tentar salvar uma notificação na base de dados, em caso de falha, o payload da requisição é enviado para a DLQ.

## Considerações

### Branches e Commits

Por estar trabalhando nesse projeto sozinho e de forma temporária, fiz commits direto na main, em contextos profissionais, sempre separo as minhas branches seguindo boas práticas separando o trabalho em branches diferentes e nomeando as branches de acordo com o seu objetivo: `feat/`, `bug/`, `hotfix/`.

### REST API

Em toda a lógica das rotas da API, respeitando a LGPD, em nenhum momento o CPF do usuário é armazenado no PostgreSQL ou no Redis em forma de texto, somente de forma encriptada e usando a estratégia de `Blind Index`.

### Clean Code

Nas rotas API e nas funções usadas pelo banco de dados, o código foi escrito seguindo as regras do clean code, usando interfaces e abstrações para diminuir o acoplamento entre os pacotes.

### WebSocket

O websocket montado armazena os clientes em memória, o que limita muito a quantidade de clientes que podem se conectar, uma forma robusta deve armazenar esses registros dos clientes conectados em um local de rápido acesso (talvez um Redis separado), também para manter a conexão ativa, deve ser implementado uma lógica de heartbeat para que a conexão não seja fechada por inatividade.

### Testes

Foram montados testes unitários e de integração somente a título de ilustração, não estão cobrindo todas as funcionalidades e nem todos os cenários dessas funcionalidades, com mais tempo os testes podem ser evoluídos para além de cobrir tudo, exportar um arquivo para utilizar em ferramentas de BI.

### Observabilidade

No projeto só foi montada uma rota de `/health`, para amadurecer a observabilidade, pode-se adicionar integrações com Datadog, New Relic, Prometheus, entre outros. Além disso, integrar ao SonarQube traria mais qualidade ao código apontando bugs, vulnerabilidades e code smells encontrados.

### CI/CD

Pela natureza do projeto ser pessoal, sem deploy, foi montado um arquivo simples para rodar os testes no Github Actions de forma ilustrativa.

### DLQ

Uma melhor implementação de uma DLQ pode ser feita utilizando serviços separados voltados somente para esse propósito, como Kafka ou RabbitMQ, dependendo do cenário de requisitos que essa DLQ deve corresponder. Outro ponto de atenção é que a DLQ armazena o CPF do usuário em forma de texto, isso só foi mantido por conta do tempo e natureza de teste do projeto, em contexto de produção o CPF seria criptografado assim como é feita na rota de `POST`.

### Worker

Para que a funcionalidade da DLQ fique completa dever ser construído um worker para executar os registros da DLQ e em caso de sucesso, removê-los da fila. Para isso, seria interessante montar uma rota exclusiva para esse worker que não tenha as barreiras da rota de `POST`, já que a requisição só cai na DLQ após ser validada.

### Observações

Por ser meu primeiro contato com a linguagem Go e suas ferramentas, bibliotecas e frameworks posso não ter utilizado o que há de melhor, mais moderno e seguro que o mercado utiliza, busquei utilizar as versões mais recentes e mais leves para ser um projeto rápido e com as tecnologias atuais.