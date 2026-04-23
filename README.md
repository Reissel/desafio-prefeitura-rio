# Desafio Prefeitura Rio

Este projeto consiste em uma API REST desenvolvida em Go utilizando o framework Gin, persistência de dados no PostgreSQL e gerenciamento de cache de requisições no Redis.

A aplicação é totalmente containerizada para garantir que o ambiente se mantém em qualquer que seja o local que seja executado.

Na Cidade do Rio de Janeiro, qualquer cidadão consegue criar um chamado por meio de um sistema público da Prefeitura para solicitar reparos de ordem pública como reparos na iluminação pública (postes), buracos nas vias, entre outros.

Esse projeto consiste em um sistema de notificações, responsável por receber atualizações dos chamados via REST API e atualizar os cidadãos em tempo real via WebSocket.

## Requisitos

- Docker Engine.
- Docker Compose.

## Arquitetura em Containers

1. **app_go**: A aplicação Go rodando na porta `8080`.
2. **postgres_db**: Instância do PostgreSQL na porta `5432`.
3. **redis_cache**: Instância do Redis na porta `6379`.

Além dos serviços essenciais, foram acrescentados serviços para facilitar a visualização dos dados no PostgreSQL e Redis por meio de GUIs disponibilizadas pelo **pgadmin** e **redis_insight**.

## Como rodar o projeto

Como o único pré-requisito é possuir o **Docker** e **Docker Compose** instalados, altere as variáveis de ambiente conforme necessário no arquivo `docker-compose.yml` na raiz do projeto e execute o comando: `docker compose up -d`.

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
