package domain

// This package re-exports the internal domain types as public aliases so
// other modules (for example `office`) can depend on the domain without
// importing internal packages directly.

import id "github.com/brunojet/store-go/shared/internal/domain"

type Categoria = id.Categoria
type TipoCategoria = id.TipoCategoria
type TipoIntegracao = id.TipoIntegracao
type ModeloTerminal = id.ModeloTerminal
type Aplicativo = id.Aplicativo
type AppCategoria = id.AppCategoria
type Configuracao = id.Configuracao
type Contato = id.Contato
type Imagem = id.Imagem
type Estagio = id.Estagio
type VersaoAplicativo = id.VersaoAplicativo
type Cadastro = id.Cadastro
type ConfiguracaoCadastro = id.ConfiguracaoCadastro
type CatalogoAplicativo = id.CatalogoAplicativo

// Lista de entidades para AutoMigrate, na ordem do arquivo
// Ordem correta para AutoMigrate, respeitando dependências entre entidades
// 1. Bases sem dependências
// 2. Entidades dependentes
// 3. Entidades com múltiplas FKs
var EntidadesAutoMigrate = []interface{}{
	&TipoCategoria{},        // base para Categoria
	&Categoria{},            // depende de TipoCategoria
	&TipoIntegracao{},       // base para Configuracao, CatalogoAplicativo
	&ModeloTerminal{},       // base para Configuracao, CatalogoAplicativo
	&Imagem{},               // base para VersaoAplicativo
	&Contato{},              // base para Cadastro
	&Aplicativo{},           // base para AppCategoria, Configuracao, Cadastro, CatalogoAplicativo
	&Configuracao{},         // depende de TipoIntegracao, ModeloTerminal, Aplicativo
	&VersaoAplicativo{},     // depende de Configuracao, Imagem
	&Cadastro{},             // depende de Aplicativo, Contato, DetalhesAplicativo
	&ConfiguracaoCadastro{}, // depende de Cadastro, Configuracao
	&AppCategoria{},         // depende de Aplicativo, Categoria
	&Estagio{},              // base para CatalogoAplicativo
	&CatalogoAplicativo{},   // depende de TipoIntegracao, ModeloTerminal, Estagio, Aplicativo, VersaoAplicativo, Cadastro
}
