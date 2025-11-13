# URL Shortener

Microsserviço de encurtamento de URLs utilizando Docker, Redis, PostgreSQL e FastAPI.
Focando em criar uma solução de encurtamento de URLs de alto desempenho e fácil uso (apenas criar os arquivos .env e subir os Docker Compose).

## Como iniciar

Clone o repo:
```bash
git clone https://github.com/vassao1/urlshortener.git
```

Crie um arquivo `.env` seguindo o exemplo em ```.envexample```, bote-o tanto na pasta app quanto na pasta principal (sei lá, só funcionou assim kkkkkkk).

Suba o Compose do encurtador.
```bash
docker compose up -d --build
```

Suba o compose do Caddy.
```bash
docker compose -f compose-caddy.yaml up -d
```

A API estará disponível em: http://localhost:80 (ou configura o caddy pra servir no seu dominio, sei la kkkkkkk)

## Endpoints

### POST /shorten
Cria uma URL encurtada

```bash
curl -X POST http://localhost:8000/short \
  -H "Content-Type: application/json" \
  -d '{"url": "github.com"}'
```

Resposta:
```json
{
  "short_url": "http://localhost:8000/short/aBcD"
}
```

Status: 201 Created

### GET /short/{hash}
Redireciona para a URL original

```bash
curl -L http://localhost:80/short/aBcD
```

Redireciona para: https://github.com
Status: 301 Moved Permanently

# Notas:

## Escalabilidade:
Dá pra aumentar a pool de conexões com a db (não sei qual é o limite disso) pelo main.py, dá também pra aumentar a quantidade de replicas do container de encurtamento (literalmente só mudar o argumento replicas do compose) e automaticamente o Caddy vai saber que tem mais um serviço de encurtamento pronto para uso. <br>
Dá pra mexer também na política de load balancing do Caddy, utilizando tanto round robin quanto least connected, nos testes eles tiveram ambos resultados diferentes conforme a quantidade de requisições por segundo.

## "E a persistência dos dados?"
Os dados são guardados nos volumes do Docker. Só fazer backup dessa pasta.

## Caddy/reverse proxy:
Pode utilizar o compose atual do Caddy para fazer reverse proxy de outras aplicações também. <br>
Não vou ensinar como utiliza o Caddy ou redes do Docker, mas é possível tanto criar novos microsserviços e conectá-los na rede `microservices` para utilização deles por esse mesmo reverse proxy quanto simplesmente utilizar o Caddy para reverse proxying de outras partes da sua aplicação fora do Docker.

# Perspectivas futuras:
Implementar um cache, só isso. <br>
Mas, no momento, por conta do uso do status 301 do HTTP, cada requisição só é feita uma vez por navegador, o resto o navegador automaticamente puxa do seu cache e redireciona pro link correto, o que pra um projeto de **no máximo** médio porte, tá bom até. <br>
Adicionar também a opção de criar links personalizados no futuro.