Office module — estrutura proposta (detalhada)

Objetivo
--------
Criar um scaffold mínimo para o módulo `office` que implemente CRUD para as entidades:

- `TipoIntegracao`
- `ModeloTerminal`
- `TipoCategoria`
- `Categoria`

O módulo `office` deve consumir os modelos de domínio já existentes em `shared` e as implementações de repositório em `infra`.

Princípios principais
---------------------
- Domínio puro (shared) não depende de infra.
- Interfaces de repositório (ports) ficam em `shared` para reuso por múltiplos módulos.
- Implementações concretas (GORM) ficam em `infra` e implementam as interfaces de `shared`.
- Migrações (AutoMigrate) pertencem a `infra`.

Checklist curto
---------------
- [ ] Domain models em `shared/internal/domain` (já existem).
- [ ] Interfaces de repositório em `shared/pkg/repo` (criar).
- [ ] Implementação GORM em `infra/pkg/repo` (usar generic repo onde aplicável).
- [ ] HTTP handlers e serviços em `office`.
- [ ] Migração em `infra/migrate.go` importando `shared`.

Estrutura de diretórios proposta
--------------------------------

```
store-go/
├─ shared/
│  ├─ internal/domain/        # domain models (ModeloTerminal, Categoria, ...)
│  └─ pkg/repo/               # repository interfaces (ports)
├─ infra/
│  ├─ pkg/repo/               # GORM adapters implementing shared ports
│  ├─ models/                 # optional: persistence-specific structs (gorm tags)
│  └─ migrate.go              # AutoMigrate ordenado; importa shared domain
├─ office/                    # novo módulo a criar
│  ├─ application/
│  │  ├─ port/                # (opcional) portas locais ou extensões
│  │  └─ service/             # services / use-cases
│  ├─ http/                   # handlers e roteador
│  ├─ repo/                   # wiring/adapters locais (pequeno)
│  └─ tests/                  # integração + serviços
└─ docs/
	 └─ office_structure.md     # este arquivo
```

Arquivos-chave e exemplos rápidos
---------------------------------

- `shared/pkg/repo/modelo_terminal_repository.go` (interface)

	Exemplo de contrato (Go):

	```go
	package repo

	import (
			"context"
			"github.com/brunojet/store-go/shared/internal/domain"
	)

	type ModeloTerminalRepository interface {
			Create(ctx context.Context, m *domain.ModeloTerminal) error
			Update(ctx context.Context, m *domain.ModeloTerminal) error
			GetByID(ctx context.Context, id int64) (*domain.ModeloTerminal, error)
			List(ctx context.Context, p ListParams) ([]domain.ModeloTerminal, error)
			Delete(ctx context.Context, id int64) error
			WithTx(tx interface{}) ModeloTerminalRepository // optional, infra-specific
	}
	```

- `infra/pkg/repo/modelo_terminal_gorm.go`

	- Implementa `ModeloTerminalRepository` usando GORM. Mapeia domain ↔ persistence.
	- Usa o `generic repository` onde aplicável para evitar duplicação.

- `office/application/service/modelo_terminal_service.go`

	- Validações, regras de negócio, DTOs↔domain mapping.
	- Recebe o `shared` repository por injeção de dependência.

- `office/http/handler_modelo_terminal.go`

	- Rotas HTTP (ex.: POST /modelo-terminals, GET /modelo-terminals, GET /modelo-terminals/:id, PUT, DELETE).
	- Usa o service e converte domain → JSON para resposta.

- `infra/migrate.go` (na raiz de infra)

	- AutoMigrate em ordem segura: `TipoCategoria`, `TipoIntegracao`, `Categoria`, `ModeloTerminal`.
	- IMPORTANTE: remover `migrate.go` de `shared/internal/repo` e colocá-lo em `infra` para evitar dependências invertidas.

Mapping domain ↔ persistence
----------------------------

- Prefer fazer mapas explícitos em `infra/pkg/repo/mappers.go`:

	- func toPersistence(d *domain.ModeloTerminal) *models.ModeloTerminal
	- func toDomain(m *models.ModeloTerminal) *domain.ModeloTerminal

	Isso evita poluir o domínio com tags GORM.

Migrações e ordem
-----------------

- `infra/migrate.go` deve chamar AutoMigrate com os tipos do domínio ou com structs de `infra/models` mapeadas a partir do domínio.
- Exemplo de ordem:
	1. TipoCategoria
	2. TipoIntegracao
	3. Categoria (FK para TipoCategoria)
	4. ModeloTerminal (FK para Categoria e TipoIntegracao)

Testes recomendados
-------------------

- `shared`: unit tests de validação/BeforeCreate para cada entidade.
- `office` services: unit tests com repositório mockado (table-driven para cenários de erro e sucesso).
- `office` integration_test: roda `infra` adapters reais com sqlite in-memory; usa `infra/migrate.go`.

Exemplo de fluxo para criar um ModeloTerminal via API
----------------------------------------------------

1. HTTP handler recebe POST /modelo-terminals com JSON.
2. Handler converte para DTO e chama `TerminalService.Create(ctx, dto)`.
3. Service valida DTO, cria `domain.ModeloTerminal` e chama `repo.Create(ctx, domainModel)`.
4. Repo GORM persiste (mapeia domain → persistence) e retorna ID.
5. Service retorna 201 Created com localização e payload.

Próximo passos sugeridos (escolha)
---------------------------------

- [A] Gerar scaffold do módulo `office` (arquivos: interface, service, handler, wiring, testes básicos).
- [B] Mover `shared/internal/repo/migrate.go` → `infra/migrate.go` e ajustar imports.
- [C] Implementar um `infra` adapter mínimo para `ModeloTerminal` usando o generic repo.

Indique a letra da opção que prefere (A, B ou C) e eu aplico o patch correspondente.

