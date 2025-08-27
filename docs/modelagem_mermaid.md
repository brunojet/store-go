# Modelagem de Dados (Mermaid)

```mermaid
erDiagram
	APP ||--o{ CADASTRO : "tem cadastro"
	CADASTRO }o--|| CONTATO : "contato principal"
	CADASTRO }o--|| DETALHES_APLICATIVO : "detalhes"
	CADASTRO ||--o{ CONFIGURACAO_CADASTRO : "configurações"
	APP ||--o{ VERSAO_APP : "versões"
	VERSAO_APP ||--o{ CATALOGO_APLICATIVO : "catalogado"
	CATALOGO_APLICATIVO }o--|| ESTAGIO_CATALOGO : "estágio"
	CATALOGO_APLICATIVO }o--|| VERSAO_APP : "versão"
	CATALOGO_APLICATIVO }o--|| DETALHES_APLICATIVO : "detalhes ativos"
	CATALOGO_APLICATIVO }o--|| SCREENSHOT : "screenshots ativos"
	CAT ||--o{ APP_CAT : "classificado"
	APP_CAT }o--|| APP : "referencia app"
	APP_CAT }o--|| CAT : "categoria"
	// ...existing code...
```

## Legenda
- APP: Aplicativo
- CADASTRO: Cadastro do Aplicativo
- CONTATO: Contato principal do cadastro
- DETALHES_APLICATIVO: Detalhes do aplicativo (conteúdos, n screenshots)
- CONFIGURACAO_CADASTRO: Configurações do cadastro
- VERSAO_APP: Versão do Aplicativo
- CATALOGO_APLICATIVO: Controle de exposição/ativação por estágio
- ESTAGIO_CATALOGO: Estágio do catálogo
- SCREENSHOT: Imagens/screenshots do aplicativo
- APP_CAT: Associação App-Categoria
- CAT: Categoria
// ...existing code...

> Diagrama gerado automaticamente por GitHub Copilot.
