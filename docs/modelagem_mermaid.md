# Modelagem de Dados (Mermaid)

```mermaid
erDiagram
    APP ||--o{ CFG : "possui"
    CFG ||--o{ VERSAO_APP : "possui"
    APP ||--o{ APP_CAT : "classificado"
    APP_CAT }o--|| CAT : "categoria"
    CFG }o--|| MODELO_TRML : "usa modelo"
    CFG }o--|| TIPO_INT : "usa integração"
    VERSAO_APP ||--o{ CAT_APP : "catalogado"
    CAT_APP }o--|| APP : "referencia app"
    CAT_APP }o--|| MODELO_TRML : "modelo catálogo"
    CAT_APP }o--|| TIPO_INT : "integração catálogo"
    CAT_APP }o--|| EST_CAT : "estágio catálogo"
    CAT_APP }o--|| VERSAO_APP : "versão catálogo"
    CAT ||--o{ TIPO_CAT : "tipo"
    EST_CAT ||--o{ CAT_APP : "usado em"
```

## Legenda
- APP: Aplicativo
- CFG: Configuração
- VERSAO_APP: Versão do Aplicativo
- APP_CAT: Associação App-Categoria
- CAT: Categoria
- MODELO_TRML: Modelo de Terminal
- TIPO_INT: Tipo de Integração
- CAT_APP: Catálogo de Aplicativo
- EST_CAT: Estágio do Catálogo
- TIPO_CAT: Tipo de Categoria

> Diagrama gerado automaticamente por GitHub Copilot.
