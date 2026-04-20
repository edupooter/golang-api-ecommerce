# Golang API Ecommerce

Projeto mínimo de API de e-commerce em Go com repositório em memória, opção SQLite e testes unitários.

Comandos básicos

- Executar servidor (em memória):

```bash
go run ./cmd/server
```

- Executar testes:

```bash
go test ./...
```

Suporte a SQLite

- Para usar SQLite defina a variável de ambiente `SQLITE_PATH` apontando para o arquivo do banco (será criado se não existir). Exemplo PowerShell:

```powershell
$env:SQLITE_PATH = 'products.db'
go run ./cmd/server
```

- CMD:

```cmd
set SQLITE_PATH=products.db
go run ./cmd/server
```

- Bash:

```bash
export SQLITE_PATH=products.db
go run ./cmd/server
```

Ao usar `SQLITE_PATH` o servidor usará o repositório SQLite (`internal/repo/sqlite.go`). Sem `SQLITE_PATH` usa o repositório em memória.

REST Client (VS Code)

Incluí um conjunto de requisições para a extensão REST Client e um arquivo de ambiente.

- Arquivos adicionados na raiz do projeto:
	- `requests.http` — exemplos de `GET`, `POST`, `PUT`, `DELETE` para o recurso `/products`.
	- `rest-client.env.json` — define o ambiente `Local` com `baseUrl` = `http://localhost:8080`.

- Como usar:
	1. Inicie o servidor (`go run ./cmd/server` ou via `server.exe`).
	2. Abra `requests.http` no VS Code.
	3. No canto superior direito do editor selecione o ambiente `Local` (ou use a variável embutida `@baseUrl`).
	4. Clique em "Send Request" acima de qualquer requisição.

Notas

- O driver SQLite usado é `modernc.org/sqlite` (implementação pura em Go), a dependência já foi adicionada ao módulo e o projeto buildou com sucesso.
- O arquivo do banco (`products.db`) será criado no diretório onde o servidor for executado.

