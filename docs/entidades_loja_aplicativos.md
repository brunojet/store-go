# Entidades do Modelo de Loja e Aplicativos

## TipoCategoria
Representa o tipo de uma categoria de aplicativo ou produto. Utilizada como base para a entidade Categoria.

---

## Categoria
Categoria de aplicativo ou produto, vinculada a um TipoCategoria.

---

## TipoIntegracao
Define o tipo de integração disponível para aplicativos ou sistemas. Base para ConfiguracaoAplicativo e CatalogoAplicativo.

---

## ModeloTerminal
Modelo de terminal (hardware/dispositivo) suportado pelo sistema. Base para ConfiguracaoAplicativo e CatalogoAplicativo.

---

## Imagem
Armazena imagens relacionadas a aplicativos, versões ou outros elementos. Base para VersaoAplicativo.

---

## ContatoAplicativo
Informações de contato associadas a um aplicativo. Base para HistoricoPerfilAplicativo.

---

## Aplicativo
Entidade principal que representa um aplicativo cadastrado na loja. Possui relacionamentos diretos com HistoricoPerfilAplicativo e ConfiguracaoAplicativo, refletindo o histórico de perfis e as configurações vinculadas ao aplicativo. Base para CatalogoAplicativo.

---

## DetalheAplicativo
Detalhes adicionais sobre o aplicativo, podendo ser compartilhados entre vários históricos de perfil.

---

## ConfiguracaoAplicativo
Configuração específica de um aplicativo, vinculada a TipoIntegracao, ModeloTerminal e Aplicativo.

---

## VersaoAplicativo
Versão específica de um aplicativo, vinculada a ConfiguracaoAplicativo e Imagem.

---

## HistoricoPerfilAplicativo
Histórico de perfis e alterações de um aplicativo, vinculado a Aplicativo, ContatoAplicativo e DetalheAplicativo.

---

## Estagio
Representa o estágio (status/fase) de um aplicativo ou catálogo. Base para CatalogoAplicativo.

---

## CatalogoAplicativo
Catálogo de aplicativos disponíveis, vinculado a TipoIntegracao, ModeloTerminal, Estagio, Aplicativo, VersaoAplicativo e HistoricoPerfilAplicativo.

---
