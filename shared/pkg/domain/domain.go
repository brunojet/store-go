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
type CategoryType = id.CategoryType
type Category = id.Category
type TipoIntegracao = id.TipoIntegracao
type TerminalModel = id.TerminalModel
type Imagem = id.Imagem
type ContatoAplicativo = id.ContatoAplicativo
type Aplicativo = id.Aplicativo
type DetalheAplicativo = id.DetalheAplicativo
type ConfiguracaoAplicativo = id.ConfiguracaoAplicativo
type VersaoAplicativo = id.VersaoAplicativo
type HistoricoPerfilAplicativo = id.HistoricoPerfilAplicativo
type Estagio = id.Estagio
type CatalogoAplicativo = id.CatalogoAplicativo
type StorageObject = id.StorageObject
type Image = id.Image
type ImageType = id.ImageType
type Video = id.Video
type VideoType = id.VideoType

var EntidadesAutoMigrate = []interface{}{
	&StorageObject{},             // base para Image, Video
	&Image{},                     // depende de StorageObject
	&Video{},                     // depende de StorageObject
	&CategoryType{},              // base para Category
	&Category{},                  // depende de CategoryType
	&TipoIntegracao{},            // base para Configuracao, CatalogoAplicativo
	&TerminalModel{},             // base para Configuracao, CatalogoAplicativo
	&Imagem{},                    // base para VersaoAplicativo
	&ContatoAplicativo{},         // base para Cadastro
	&Aplicativo{},                // base para AppCategory, Configuracao, Cadastro, CatalogoAplicativo
	&DetalheAplicativo{},         // pode ser compartilhado entre vários cadastros
	&ConfiguracaoAplicativo{},    // depende de TipoIntegracao, TerminalModel, Aplicativo
	&VersaoAplicativo{},          // depende de Configuracao, Imagem
	&HistoricoPerfilAplicativo{}, // depende de Aplicativo, Contato, DetalhesAplicativo
	&Estagio{},                   // base para CatalogoAplicativo
	&CatalogoAplicativo{},        // depende de TipoIntegracao, TerminalModel, Estagio, Aplicativo, VersaoAplicativo, Cadastro
}
