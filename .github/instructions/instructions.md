---
---
description: Carregado quando o agente estiver trabalhando em alterações de código no repositório "Golang API Ecommerce"
applyTo: "internal/**, cmd/**, *.go"
---

Resumo
-------
Este arquivo fornece contexto do projeto e regras práticas para agentes LLM que geram ou revisam código neste repositório. Carregue essas instruções quando o agente for solicitado a modificar, refatorar, testar ou documentar código em `internal/`, `cmd/` ou em arquivos Go.

Contexto do projeto
---------------------
- Módulo: `github.com/edupooter/golang-api-ecommerce`
- Go: 1.25
- Arquitetura escolhida: Hexagonal (ports/adapters/services/handlers)
- Banco: opção In-Memory ou SQLite (driver `modernc.org/sqlite`). Use `SQLITE_PATH` para apontar o arquivo DB.
- Swagger: docs geradas com `swag` (github.com/swaggo/swag) e servidas por `github.com/swaggo/http-swagger` em `/swagger/`.

Comandos úteis
--------------
- Rodar servidor (memória):

```bash
go run ./cmd/server
```

- Rodar servidor com SQLite (PowerShell):

```powershell
$env:SQLITE_PATH = 'products.db'
go run ./cmd/server
```

- Gerar/atualizar Swagger:

```bash
# com o binário swag instalado
swag init -g cmd/server/main.go -o docs

# ou sem instalar
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o docs
```

- Testes e build:

```bash
go test ./...
go build ./...
```

Decisões e princípios (cumprir sempre)
-------------------------------------
- DDD / SOLID / DRY / KISS / YAGNI: mantenha handlers leves, coloque regras de negócio em `internal/service` e modelos em `internal/model`.
- Dependências e injeção: use interfaces (ports) e injete adaptadores via construtores para facilitar testes.
- Concorrência/estoque: use `UPDATE ... WHERE stock >= ?` (condicional) para o DB; proteja repositórios em memória com mutex.
- Transações: inicie com transações simples ao persistir Order + débito de estoque; evolua para padrões mais complexos apenas quando necessário.
- Erros: padronize erros exportados (ex.: `repo.ErrNotFound`, `ports.ErrInsufficientStock`) e centralize mapeamento HTTP em `internal/server/errors.go`.

Boas práticas para o agente LLM
-------------------------------
- Antes de alterar código, verifique `go.mod`, `cmd/server/main.go`, `internal/server/router.go` e os handlers afetados.
- Ao adicionar endpoints HTTP: crie handlers concisos em `internal/handler`, delegue lógica a services e escreva testes unitários.
- Regenerar Swagger sempre que anotações `// @...` forem modificadas.
- Ao adicionar dependências: justifique o motivo, prefira bibliotecas pequenas e mantidas.
- Cobertura de testes: unidades para domain/serviço com mocks; testes de integração opcionais com SQLite em arquivo temporário.

Checklist de PR
---------------
- Rodar `go test ./...` — todos os testes devem passar.
- Rodar `go vet` e `golangci-lint` (se disponível).
- Executar `go fmt` e `go mod tidy`.
- Regenerar `docs/` com `swag init` após mudanças em anotações.
- Incluir/atualizar testes que cobrem o comportamento novo/refatorado.

Arquivos de referência
----------------------
- `cmd/server/main.go` — ponto de entrada, variáveis de ambiente.
- `internal/server/router.go` — registro de rotas (Swagger, handlers).
- `internal/handler/` — handlers HTTP.
- `internal/service/` — regras de negócio e orquestrações.
- `internal/repo/` — adaptadores (memory, sqlite).
- `internal/model/` — entidades e value objects.
- `docs/` — arquivos gerados pelo `swag`.

Notas finais
-----------
Estas instruções são um guia prático; quando houver conflito entre uma decisão técnica e uma necessidade urgente do usuário, priorize segurança, testes e mínima invasão do código existente. Atualize este arquivo sempre que decisões de arquitetura importantes forem tomadas.

