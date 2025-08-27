# Modelagem de Dados (Mermaid)

```mermaid
---
config:
  layout: elk
---
erDiagram
	direction TB
	MODELO_TERMINAL {
	}
	TIPO_INTEGRACAO {
	}
	TIP_CAT {
	}
	CAT {
	}
	APP {
	}
	CONFIGURACAO {
	}
	CADASTRO {
	}
	CONTATO {
	}
	DETALHES_APLICATIVO {
	}
	IMAGEM_DETALHE {
	}
	IMAGEM {
	}
	ANEXO {
	}
	CONFIGURACAO_CADASTRO {
	}
	APP_CAT {
	}
	VERSAO_APP {
	}
	CATALOGO_APLICATIVO {
	}
	ESTAGIO {
	}

	CADASTRO}o--o{APP_CAT:"classificado"
	CAT}o--||TIP_CAT:"tipo categoria"
	CAT||--o{APP_CAT:"classificado"
	CONFIGURACAO}o--||TIPO_INTEGRACAO:"tipo integração"
	CONFIGURACAO}o--||MODELO_TERMINAL:"modelo terminal"
	CONFIGURACAO}o--||APP:"aplicativo"
	APP||--o{CADASTRO:"tem cadastro"
	CADASTRO}o--||CONTATO:"contato principal"
	CADASTRO}o--||DETALHES_APLICATIVO:"detalhes"
	DETALHES_APLICATIVO||--o{IMAGEM_DETALHE:"imagens detalhes"
	IMAGEM_DETALHE}o--||IMAGEM:"reuso de imagem"
	IMAGEM||--||ANEXO:"imagem física"
	CADASTRO||--o{CONFIGURACAO_CADASTRO:"configurações"
	CONFIGURACAO_CADASTRO}o--||CONFIGURACAO:"referencia configuracao"
	VERSAO_APP}o--||CONFIGURACAO:"configuração da versão"
	VERSAO_APP}o--||IMAGEM:"imagem da versão"
	CATALOGO_APLICATIVO}o--||ESTAGIO:"estágio"
	CATALOGO_APLICATIVO}o--||VERSAO_APP:"versão"
	CATALOGO_APLICATIVO}o--||CADASTRO:"cadastro"
	CATALOGO_APLICATIVO}o--||TIPO_INTEGRACAO:"tipo integração"
	CATALOGO_APLICATIVO}o--||MODELO_TERMINAL:"modelo terminal"
	CATALOGO_APLICATIVO}o--||APP:"aplicativo"
```

## Legenda
## Legenda
- APP: Aplicativo
- CADASTRO: Registro histórico do aplicativo, vinculado a contato e detalhes
- CONTATO: Dados de contato principal do cadastro
- DETALHES_APLICATIVO: Informações detalhadas do aplicativo (conteúdo, screenshots)
- IMAGEM_DETALHE: Associação entre detalhes do aplicativo e imagens reutilizáveis
- IMAGEM: Representação lógica da imagem
- ANEXO: Arquivo físico vinculado à imagem (1:1)
- CONFIGURACAO_CADASTRO: Associação entre cadastro e configuração técnica
- CONFIGURACAO: Configuração técnica do aplicativo (modelo, integração, app)
- MODELO_TERMINAL: Tipo/modelo de terminal suportado
- TIPO_INTEGRACAO: Tipo de integração suportada
- VERSAO_APP: Versão do aplicativo, vinculada à configuração e imagem
- CATALOGO_APLICATIVO: Exposição/ativação do aplicativo por estágio, versão, modelo, integração e cadastro
- ESTAGIO: Estágio do catálogo (ex: homologação, produção)
- APP_CAT: Associação muitos-para-muitos entre aplicativo e categoria
- CAT: Categoria do aplicativo
- TIP_CAT: Tipo de categoria

> Diagrama gerado automaticamente por GitHub Copilot.

## Requisitos
1. Catalogo de aplicativos deve ter combinação (modelo terminal, tipo integração, estágio e aplicativo) únicos.
2. Configuração de aplicativos deve ter combinação (modelo terminal, tipo integração e aplicativo) únicos.
3. O cadastro de aplicativos de ter reúso dos contatos, detalhes dos aplicativos, caso não sejam alterados.
4. Os anexos e imagens devem ter reúso entre todas as entidades.
