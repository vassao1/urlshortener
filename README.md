# URL Shortener

Microsserviço de encurtamento de URLs utilizando Docker, Redis, PostgreSQL e FastAPI.

## Como Iniciar

Clone o repo:
```bash
git clone https://github.com/vassao1/urlshortener.git
```

Crie um arquivo `.env` seguindo o exemplo em ```.envexample```, bote-o tanto na pasta app quanto na pasta principal (sei lá, só funcionou assim kkkkkkk).

Suba o Compose do Docker.
```bash
docker compose up -d --build
```

A API estará disponível em: http://localhost:8000 (tenho que setar o caddy ainda)

## Endpoints

### POST /shorten
Cria uma URL encurtada

```bash
curl -X POST http://localhost:8000/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "github.com"}'
```

Resposta:
```json
{
  "short_url": "http://localhost:8000/aBcD"
}
```

Status: 201 Created

### GET /{hash}
Redireciona para a URL original

```bash
curl -L http://localhost:8000/aBcD
```

Redireciona para: https://github.com
Status: 301 Moved Permanently