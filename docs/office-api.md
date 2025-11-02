# Office API (esboço expandido)

Este documento descreve um esboço das APIs para Office. Mantive a preferência da companhia por chaves compostas quando faz sentido (por exemplo, configurações identificadas por combinação de aplicativo+integração+modelo), mas também proponho um `idPerfilTerminal` único opcional para facilitar operações CRUD quando desejado. Inclui exemplos, códigos HTTP e um padrão de erros.

Nota: todas as rotas abaixo são sugeridas sob versão de API `/v1`.

## Convenções

- Formato: JSON, camelCase para campos.
- Identificadores: podem ser UUIDs (ex: `id`) ou chaves compostas quando indicado.
- Paginação: `limit` e `offset` nos endpoints de listagem.
- Erros: resposta padrão `{ status, error, message, details? }`.
- Autenticação: presumir JWT ou API Key (não detalhado aqui) — incluir header `Authorization: Bearer <token>`.

Nota sobre atualizações: este projeto padroniza o uso de PATCH para atualizações parciais (em vez de PUT). Recomendamos `Content-Type: application/merge-patch+json` para requests de PATCH e suporte a ETag/If-Match para controle de concorrência (optimistic locking).

## Integracao

Endpoints:

- POST /v1/integracoes -> cria integração
- GET /v1/integracoes/{id} -> consulta integração
- GET /v1/integracoes -> lista integrações (limit/offset)
- PATCH /v1/integracoes/{id} -> atualizações parciais
- DELETE /v1/integracoes/{id} -> remove integração (204 ou 409 se referenciada)

Payload (request/response exemplo):

```json
Request POST /v1/integracoes
{
	"nome": "Nome curto do tipo de integração",
	"descricao": "Descrição da integração que o aplicativo do desenvolvedor utilizará para se comunicar com o aplicativo de pagamento"
}

Response 201 Created
{
	"id": "8a4f1e2a-...",
	"name": "AcquirerX",
	"description": "Integração com AcquirerX",
	"createdAt": "2025-11-01T12:00:00Z"
}
```

Campos:

- id: string (UUID, gerado pelo servidor)
- name: string (required)
- description?: string

Erros:

- 400 Bad Request — validação
- 404 Not Found — buscar id inexistente
- 409 Conflict — nome duplicado

## Modelos de Terminal

Endpoints:

- POST /v1/modelos-terminal
- GET /v1/modelos-terminal/{id}
- GET /v1/modelos-terminal
- PATCH /v1/modelos-terminal/{id}
- DELETE /v1/modelos-terminal/{id}

```json
Exemplo POST /v1/modelos-terminal
{
	"name": "PAX-A920",
	"description": "Terminal PAX A920 versão X"
}

Resposta 201
{
	"id": "b3c3...",
	"name": "PAX-A920",
	"description": "Terminal PAX A920 versão X"
}
```

Campos:

- id: string (UUID)
- name: string
- description?: string

## Perfil terminal

Endpoints:

- POST /v1/perfis-terminal
- GET /v1/perfis-terminal/{id}
- GET /v1/perfis-terminal
- PATCH /v1/perfis-terminal/{id}

```json
Exemplo POST /v1/perfis-terminal
{
    "idIntegracao": "8a4f...",
    "modeloTerminalId": "b3c3...",
    "name": "Perfil Padrão",
    "properties": {
        "timeout": 30,
        "currency": "BRL"
    }
}

Resposta 201
{
    "id": "d2e4...",
    "idIntegracao": "8a4f...",
    "modeloTerminalId": "b3c3...",
    "name": "Perfil Padrão",
    "properties": {
        "timeout": 30,
        "currency": "BRL"
    },
    "createdAt": "2025-11-01T12:10:00Z",
    "updatedAt": "2025-11-01T12:10:00Z"
}
```

Campos:

- id: string (UUID)
- idIntegracao: string (UUID)
- modeloTerminalId: string (UUID)
- name: string
- properties: object (map<string, any>)

## Aplicativos

Endpoints:

- POST /v1/aplicativos
- GET /v1/aplicativos/{id}
- GET /v1/aplicativos
- PATCH /v1/aplicativos/{id}
- DELETE /v1/aplicativos/{id}

```json
Exemplo POST /v1/aplicativos
{
	"name": "LojaApp",
	"description": "Aplicativo para lojas"
}

Resposta 201
{
	"id": "4d5a...",
	"name": "LojaApp",
	"description": "Aplicativo para lojas"
}
```

Campos:

- id: string (UUID)
- name: string
- description?: string

## Tipo de filtro de aplicativo

Endpoints:

- POST /v1/tipos-filtro-aplicativo
- GET /v1/tipos-filtro-aplicativo/{id}
- GET /v1/tipos-filtro-aplicativo

Payloads e exemplos:

Request POST /v1/tipos-filtro-aplicativo

```json
{
  "nome": "Nome curto do tipo de filtro, que será exibido no frontend do aplicativo",
  "descricao": "Descrição do tipo de filtro dando detalhes sobre seu propósito"
}
```

Response 201 Created

```json
{
  "id": "e6f7...",
  "nome": "Nome curto do tipo de filtro, que será exibido no frontend do aplicativo",
  "descricao": "Descrição do tipo de filtro dando detalhes sobre seu propósito",
  "createdAt": "2025-11-01T12:15:00Z"
}
```

```json
{
  "id": "e6f7...",
  "nome": "Região",
  "descricao": "Filtra aplicativos por região geográfica",
  "createdAt": "2025-11-01T12:15:00Z"
}
```

## Filtros de aplicativo

Endpoints:

- POST /v1/tipos-filtro-aplicativo/{idTipoFiltroAplicativo}/filtros-aplicativo
- GET /v1/tipos-filtro-aplicativo/{idTipoFiltroAplicativo}/filtros-aplicativo
- DELETE /v1/tipos-filtro-aplicativo/{idTipoFiltroAplicativo}/filtros-aplicativo/{filtroId}

Payloads e exemplos:

Request POST /v1/tipos-filtro-aplicativo/e6f7.../filtros-aplicativo

```json
{
  "nome": "Centro-Oeste",
  "descricao": "Filtro para aplicativos disponíveis na região Centro-Oeste do Brasil"
}
```

Response 201 Created

```json
{
  "id": "f8g9...",
  "nome": "Centro-Oeste",
  "descricao": "Filtro para aplicativos disponíveis na região Centro-Oeste do Brasil"
}
```

## Configurações de Aplicativo

Endpoints:

- POST /v1/aplicativos/{aplicativoId}/configuracoes-aplicativo
  - Cria uma nova configuração para o aplicativo.
- GET /v1/aplicativos/{aplicativoId}/configuracoes-aplicativo -> lista (limit, offset, filtros)
- GET /v1/aplicativos/{aplicativoId}/configuracoes-aplicativo/{idPerfilTerminal} -> busca configuração por id
- PATCH /v1/aplicativos/{aplicativoId}/configuracoes-aplicativo/{idPerfilTerminal} -> atualizações parciais
- DELETE /v1/aplicativos/{aplicativoId}/configuracoes-aplicativo/{idPerfilTerminal} -> remove

Payloads e exemplos:

Request POST /v1/aplicativos/4d5a.../configuracoes-aplicativo

```json
{
  "idPerfilTerminal": "d2e4...",
  "packageName": "com.loja.app"
}
```

Response 201 Created

```json
{
  "aplicativoId": "4d5a...",
  "idPerfilTerminal": "c1a2...",
  "packageName": "com.loja.app",
  "createdAt": "2025-11-01T12:05:00Z"
}
```

Listagem (exemplo):

````json
{
	"items": [
		{
			"aplicativoId": "4d5a...",
			"idPerfilTerminal": "c1a2...",
			"packageName": "com.loja.app",
			"createdAt": "2025-11-01T12:05:00Z"
		}
	],
	"total": 1,
	"limit": 20,
	"offset": 0
}
GET /v1/aplicativos/4d5a.../configuracoes-aplicativo?limit=20&offset=0

Response 200
```json
{
	"items": [ { /* config objects como acima */ } ],
	"total": 3,
	"limit": 20,
	"offset": 0
}
````

Campos (recomendados):

- idAplicativo: string (UUID) PK parte da chave composta, FK para `aplicativos`
- idPerfil: string (UUID) — PK parte da chave composta, FK para `perfis-terminal`
- packageName: string
- createdAt, updatedAt

## Perfis de aplicativo

Endpoints:

- POST /v1/aplicativos/{aplicativoId}/perfis-aplicativo
- GET /v1/aplicativos/{aplicativoId}/perfis-aplicativo/{idPerfil}
- GET /v1/aplicativos/{aplicativoId}/perfis-aplicativo
- GET /v1/aplicativos/{aplicativoId}/perfis-aplicativo/{idPerfil}/status
- PATCH /v1/aplicativos/{aplicativoId}/perfis-aplicativo/{idPerfil}/status/revisao
  - Somente no status 'pronto'
- PATCH /v1/aplicativos/{aplicativoId}/perfis-aplicativo/{idPerfil}/status/producao
  - Somente no status 'emRevisao'
- PATCH /v1/aplicativos/{aplicativoId}/perfis-aplicativo/{idPerfil}/status/arquivar
  - Pode ser arquivado em qualquer status
- PATCH /v1/aplicativos/{aplicativoId}/perfis-aplicativo/{idPerfil}
  - Somente no status 'pendente' ou 'pronto'

Payloads e exemplos:
Request POST /v1/aplicativos/4d5a.../perfis-aplicativo

```json
{
  "nome": "Nome do aplicativo que deverá ser exibido na loja",
  "descricao": "Deve descrever o que é aplicativo e suas características",
  "filtros": [
    "f8g9...", // IDs dos filtros aplicáveis
    "f8g9..."
  ],
  "icone": {
    "nomeArquivo": "icone.png",
    "tipoMime": "image/png",
    "tamanho": 204800,
    "largura": 512,
    "altura": 512,
    "hash": "sha256-abc123..."
  },
  "screenshots": [
    {
      "nomeArquivo": "screenshot1.png",
      "tipoMime": "image/png",
      "tamanho": 102400,
      "largura": 1024,
      "altura": 768,
      "hash": "sha256-def456..."
    }
  ],
  "contato": {
    "nome": "Suporte LojaApp",
    "email": "contato@lojaapp.com",
    "telefone": "+55 11 91234-5678",
    "site": "https://www.lojaapp.com/suporte",
    "horariosAtendimento": {
      "semana": {
        "horaInicio": 9,
        "horaFim": 18
      },
      "sabado": {
        "horaInicio": 8, // Preencher com nulo se não atender
        "horaFim": 12 // Preencher com nulo se não atender
      },
      "domingoFeriados": {
        "horaInicio": 8, // Preencher com nulo se não atender
        "horaFim": 12 // Preencher com nulo se não atender
      }
    }
  }
}
```

Response 201 Created

```json
{
  "idPerfil": "c1a2...",
  "status": "pendente | pronto | emRevisao | emProducao | arquivado",
  "arquivos": [
    {
      "nomeArquivo": "icone.png",
      "tipoMime": "image/png",
      "tamanho": 204800,
      "largura": 512,
      "altura": 512,
      "hash": "sha256-abc123...",
      "status": "uploadPendente",
      "upload": {
        "method": "PUT",
        "presignedUrl": "https://storage.example.com/...?X-Amz-...",
        "expiresIn": 300,
        "requiredHeaders": { "Content-Type": "image/png" }
      }
    },
    {
      "nomeArquivo": "screenshot1.png",
      "tipoMime": "image/png",
      "tamanho": 102400,
      "largura": 1024,
      "altura": 768,
      "hash": "sha256-def456...",
      "status": "uploadPendente",
      "upload": {
        "method": "PUT",
        "presignedUrl": "https://storage.example.com/...?X-Amz-...",
        "expiresIn": 300,
        "requiredHeaders": { "Content-Type": "image/png" }
      }
    }
  ]
}
```

Request GET /v1/aplicativos/4d5a.../perfis-aplicativo/c1a2.../status

```json
{
  "idPerfil": "c1a2...",
  "status": "pendente | pronto | emRevisao | emProducao | arquivado",
  "ultimaAtualizacao": "2025-11-01T12:20:00Z",
  "pendencias": [
    {
      "nomeArquivo": "icone.png",
      "tipoMime": "image/png",
      "tamanho": 204800,
      "largura": 512,
      "altura": 512,
      "hash": "sha256-abc123...",
      "status": "uploadPendente | falha",
      "upload": {
        "method": "PUT",
        "presignedUrl": "https://storage.example.com/...?X-Amz-...",
        "expiresIn": 300,
        "requiredHeaders": { "Content-Type": "image/png" }
      }
    },
    {
      "nomeArquivo": "screenshot1.png",
      "tipoMime": "image/png",
      "tamanho": 102400,
      "largura": 1024,
      "altura": 768,
      "hash": "sha256-def456...",
      "uploadUrl": "https://storage.example.com/upload/screenshot1.png",
      "status": "emProcessamento"
    }
  ],
  "contadorPendencias": 2
}
```

Request GET /v1/aplicativos/4d5a.../perfis-aplicativo/c1a2...

Response 200 OK
```json
{
  "idAplicativo": "4d5a...",
  "idPerfil": "c1a2...",
  "nome": "Nome do aplicativo que deverá ser exibido na loja",
  "descricao": "Deve descrever o que é aplicativo e suas características",
  "filtros": [
    "f8g9...", // IDs dos filtros aplicáveis
    "f8g9..."
  ],
  "icone": {
    "nomeArquivo": "icone.png",
    "tipoMime": "image/png",
    "tamanho": 204800,
    "largura": 512,
    "altura": 512,
    "hash": "sha256-abc123...",
    "status": "uploadPendente",
    "upload": {
      "method": "PUT",
      "presignedUrl": "https://storage.example.com/...?X-Amz-...",
      "expiresIn": 300,
      "requiredHeaders": { "Content-Type": "image/png" }
    }
  },
  "screenshots": [
    {
      "nomeArquivo": "screenshot1.png",
      "tipoMime": "image/png",
      "tamanho": 102400,
      "largura": 1024,
      "altura": 768,
      "hash": "sha256-def456...",
      "status": "pronto"
    }
  ],
  "contato": {
    "nome": "Suporte LojaApp",
    "email": "contato@lojaapp.com",
    "telefone": "+55 11 91234-5678",
    "site": "https://www.lojaapp.com/suporte",
    "horariosAtendimento": {
      "semana": {
        "horaInicio": 9,
        "horaFim": 18
      },
      "sabado": {
        "horaInicio": 8, // Preencher com nulo se não atender
        "horaFim": 12 // Preencher com nulo se não atender
      },
      "domingoFeriados": {
        "horaInicio": 8, // Preencher com nulo se não atender
        "horaFim": 12 // Preencher com nulo se não atender
      }
    }
  }
}
```

Request GET /v1/aplicativos/4d5a.../perfis-aplicativo

Response 200 OK
```json
{
  "idAplicativo": "4d5a...",
  "idPerfil": "c1a2...",
  "nome": "Nome do aplicativo que deverá ser exibido na loja",
  "descricao": "Deve descrever o que é aplicativo e suas características",
  "status": "pendente | pronto | emRevisao | emProducao | arquivado",
  "ultimaAtualizacao": "2025-11-01T12:20:00Z",
  "contadorPendencias": 2
}
```

