# Contexto

Gerar um projeto em go no modelo monolito modular com as seguintes partes:

**Backoffice**: Responsável pela gestão interna da aplicação, incluindo funcionalidades administrativas e de suporte.
**Store**: Módulo responsável pela disponibilização dos aplicativos da loja.
**Shared**: Módulo responsável por conter funcionalidades e componentes compartilhados entre os outros módulos, como por exemplo a persistencia dos dados em banco de dados e sistenma de arquivos.

Tenho em mente que shared é um modulo comum aos dois e o projeto deve suportar o build combinado de backoffice + shared, store + shared e todos os módulos juntos.

---

## Proposta de Estrutura e Técnicas de Isolamento/Modularização

**Estrutura Recomendada:**
- `cmd/backoffice/` — Entry point do backoffice (`main.go`)
- `cmd/store/` — Entry point da store (`main.go`)
- `internal/backoffice/` — Lógica e módulos do backoffice
- `internal/store/` — Lógica e módulos da store
- `internal/shared/` — Funcionalidades compartilhadas (persistência, arquivos, utilitários)
- `pkg/` — Pacotes reutilizáveis (se necessário expor para outros projetos)
- `configs/` — Arquivos de configuração
- `docs/` — Documentação
- `test/` — Testes integrados e mocks

**Técnicas para Isolamento e Modularização:**
- Utilizar o padrão Go de pacotes internos (`internal/`) para garantir encapsulamento e evitar dependências indesejadas entre módulos.
- Interfaces para abstrair dependências entre módulos (ex: persistência, serviços externos).
- Injeção de dependências via construtores, facilitando testes e substituição de implementações.
- Builds combinados: cada entry point em `cmd/` pode importar apenas os módulos necessários (`backoffice + shared`, `store + shared`, ou todos juntos), permitindo builds independentes ou integrados.
- Testes unitários e de integração em cada módulo, com uso de mocks para dependências compartilhadas.
- Documentação clara das dependências e contratos entre módulos.

**Vantagens:**
- Modularidade, isolamento e fácil manutenção.
- Flexibilidade para builds combinados conforme necessidade.
- Reuso de código e componentes compartilhados sem acoplamento excessivo.

---

## Estrutura mínima para validação do modelo modular Go

**Estrutura criada:**
- `cmd/backoffice/main.go`: Entry point do backoffice
- `cmd/store/main.go`: Entry point da store
- `cmd/all/main.go`: Entry point para build combinado
- `internal/shared/hello.go`: Módulo compartilhado

**Recomendação de ajuste:**
- O import dos módulos internos deve usar o caminho relativo ao módulo Go. Certifique-se de inicializar o projeto com `go mod init store-go` na raiz para que os imports funcionem corretamente.
- Para builds e testes, utilize `go run ./cmd/backoffice`, `go run ./cmd/store` e `go run ./cmd/all`.

**Observação:**
- O projeto está pronto para validar o modelo modular e builds combinados. Se precisar adicionar mais módulos ou ajustar a estrutura, basta seguir o padrão acima.
