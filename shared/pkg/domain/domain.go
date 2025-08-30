package domain

// This package re-exports the internal domain types as public aliases so
// other modules (for example `office`) can depend on the domain without
// importing internal packages directly.

import (
	infra "github.com/brunojet/infra-go/pkg/domain"
	id "github.com/brunojet/store-go/shared/internal/domain"
)

// Exportação das principais entidades para uso externo
type BaseModel = infra.BaseModel
type BaseEntity = infra.BaseEntity
type Anexo = infra.Anexo
type TipoCategoria = id.TipoCategoria
type Categoria = id.Categoria
type TipoIntegracao = id.TipoIntegracao
type ModeloTerminal = id.ModeloTerminal
type Imagem = id.Imagem
type ContatoAplicativo = id.ContatoAplicativo
type Aplicativo = id.Aplicativo
type DetalheAplicativo = id.DetalheAplicativo
type ConfiguracaoAplicativo = id.ConfiguracaoAplicativo
type VersaoAplicativo = id.VersaoAplicativo
type HistoricoPerfilAplicativo = id.HistoricoPerfilAplicativo
type Estagio = id.Estagio
type CatalogoAplicativo = id.CatalogoAplicativo

var EntidadesAutoMigrate = []interface{}{
	&TipoCategoria{},             // base para Categoria
	&Categoria{},                 // depende de TipoCategoria
	&TipoIntegracao{},            // base para Configuracao, CatalogoAplicativo
	&ModeloTerminal{},            // base para Configuracao, CatalogoAplicativo
	&Imagem{},                    // base para VersaoAplicativo
	&ContatoAplicativo{},         // base para Cadastro
	&Aplicativo{},                // base para AppCategoria, Configuracao, Cadastro, CatalogoAplicativo
	&DetalheAplicativo{},         // pode ser compartilhado entre vários cadastros
	&ConfiguracaoAplicativo{},    // depende de TipoIntegracao, ModeloTerminal, Aplicativo
	&VersaoAplicativo{},          // depende de Configuracao, Imagem
	&HistoricoPerfilAplicativo{}, // depende de Aplicativo, Contato, DetalhesAplicativo
	&Estagio{},                   // base para CatalogoAplicativo
	&CatalogoAplicativo{},        // depende de TipoIntegracao, ModeloTerminal, Estagio, Aplicativo, VersaoAplicativo, Cadastro
}
