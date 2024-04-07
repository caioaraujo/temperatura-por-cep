# Temperatura por CEP
Dado um CEP, verifica a temperatura do local.

## Execução servidor
`go run cmd/server/main.go`

Porta padrão: 8000

## Api
Utilize a rota `GET /temperatura/{cep}` onde o valor do CEP deve ser somente números.

Response ex: `{ "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }`

### Status code
- 200: Success
- 422: Invalid zipcode
- 404: can not find zipcode
