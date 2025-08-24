package domain

// This package re-exports the internal domain types as public aliases so
// other modules (for example `office`) can depend on the domain without
// importing internal packages directly.

import id "github.com/brunojet/store-go/shared/internal/domain"

type Categoria = id.Categoria
type TipoCategoria = id.TipoCategoria
type ModeloTerminal = id.ModeloTerminal
type TipoIntegracao = id.TipoIntegracao
type Imagem = id.Imagem
