<!-- Documento: portais-integration.md
	Gerado/atualizado: primeira versão com visões e jornadas
	Objetivo: registrar como os portais (Relacionamento, Dev, Certificação, MDM, Loja) se integram
-->

# Portais e integrações

Este documento descreve a visão dos portais que compõem o ecossistema de distribuição e certificação de aplicativos, seus papéis e jornadas principais de integração entre sistemas: Portal de Relacionamento (novo), Portal Dev (IAAS), Portal Certificação (antigo portal dev ServiceNow), MDM e Loja de Aplicativos.

## Sumário

- Portal de Relacionamento (Novo)
- Portal Dev (IAAS)
- Portal Certificação
- MDM (Mobile Device Management)
- Loja de Aplicativos
- Jornadas de integração (onboarding, certificação, distribuição piloto/prod)
- Interações técnicas e contratos mínimos de API

## Portal de Relacionamento (Novo)

Objetivo: centralizar o relacionamento com parceiros — onboarding, autorização KYP (Know Your Partner), provisionamento de contas e disponibilização de ambientes/maquinas para desenvolvimento e testes.

Funcionalidades principais:
- Onboarding do parceiro: coleta de dados legais, contato, obrigações contratuais e papéis (ex.: desenvolvedor, integrador).
- Integração KYP: verificação de identidade e validação de empresa (documentos, razão social, acordos).
- Integração de conta/credenciais: criação e revogação de contas de acesso às plataformas (Portal Dev, MDM, Loja, APIs internas).
- Solicitação e disponibilização de maquinhas POS para desenvolvimento e homologação: disponibilização controlada de terminais POS (físicos ou simuladores aprovados), gerenciamento de lotes/emprestimos, firmware e imagens de testes, e logística (envio/retorno) para que o parceiro possa desenvolver e testar integrações com o SDK.
- Cadastro de aplicativos (dono da informação): parceiro cadastra metadados do app (nome, descrição, versão, bundle/id, arquivos binários ou link para repositório, equipe de contato).
- Fluxos para distribuir aplicativos para ambientes de piloto e produção (coordenação com MDM e Loja).

Público-alvo: equipes de negócios e operações que gerenciam a relação com parceiros; times de suporte técnico.

APIs mínimas (exemplos):
- POST /partners -> cria parceiro
- GET /partners/{id} -> consulta parceiro
- POST /partners/{id}/kyc -> submete resultado KYP
- POST /partners/{id}/accounts -> cria credenciais em plataformas
- POST /applications -> registra aplicativo (referência ao arquivo/binário ou artefato)
- POST /applications/{id}/release -> solicita liberação piloto/produção

## Portal Dev (Documentação, SDKs e Suporte)

Objetivo: fornecer artefatos técnicos, documentação e suporte para parceiros integrarem com o ecossistema.

Funcionalidades principais:
- Repositório de SDKs, bibliotecas, exemplos e sample apps com guias de integração passo-a-passo.
- Documentação técnica (API reference, tutoriais, melhores práticas, changelogs e notas de versão).
- Suporte técnico: sistema de tickets, FAQ, canais de comunicação (chat/email) e guias de troubleshooting.

APIs mínimas (exemplos):
- GET /sdk/{platform}/{version} -> obtém artefato SDK
- GET /docs/{topic} -> obtém documentação
- GET /examples/{name} -> baixa sample app/exemplo
- POST /support/tickets -> cria chamado de suporte

## Portal Certificação (antigo portal dev ServiceNow)

Objetivo: gerenciar o processo de certificação de aplicativos do parceiro — testes, relatórios e emissão de selos/aptidão para distribuição.

Funcionalidades principais:
- Submissão de pacote para certificação (automática ou manual).
- Execução de baterias de testes (segurança, performance, compatibilidade).
- Registro de evidências e comunicação com o parceiro sobre pendências.
- Emissão de certificado/ selo de conformidade e publicação do status para a Loja e Portal de Relacionamento.

APIs mínimas (exemplos):
- POST /certifications/submit -> submete build para análise
- GET /certifications/{id}/status -> status e relatórios
- POST /certifications/{id}/evidence -> upload de evidências

## MDM (Mobile Device Management)

Objetivo: distribuição e gerenciamento de aplicações (privadas e públicas) em dispositivos gerenciados pelo comerciante/cliente.

Funcionalidades principais:
- Distribuição de apps privados (solicitados pelo parceiro) e apps públicos (iniciados pelo comércio via loja de aplicativos).
- Gerenciamento de grupos de dispositivos, perfis e políticas de instalação/atualização.
- Integração com o processo de certificação para aprovar builds antes da distribuição.

APIs mínimas (exemplos):
- POST /mdm/apps -> disponibiliza app no MDM
- POST /mdm/groups/{id}/deploy -> deploy para grupo
- GET /mdm/apps/{id}/status -> status de distribuição

## Loja de Aplicativos

Objetivo: vitrine pública de aplicativos disponíveis para os comerciantes; interface para pesquisa, avaliação, instalação (quando aplicável) e solicitação para pilotos.

Funcionalidades principais:
- Vitrine pública com filtros, buscas e categorias.
- Páginas de detalhe do aplicativo com informações, screenshots, changelog e histórico de versões.
- Fluxos de solicitação para piloto (solicitar acesso) e botão para iniciar instalação (dependendo do canal: MDM, link de download, deep link).
- Integração com Portal de Relacionamento para exibir dados do parceiro e com Portal Certificação para mostrar selo de conformidade.

APIs mínimas (exemplos):
- GET /store/apps -> listagem pública
- GET /store/apps/{id} -> detalhe do app
- POST /store/apps/{id}/request-pilot -> solicita participação em piloto

## Jornadas de integração

As jornadas descrevem como os sistemas interagem nas operações mais comuns.

### 1) Onboarding do Parceiro

1. Parceiro acessa o Portal de Relacionamento e submete formulário de cadastro.
2. Portal de Relacionamento cria registro (POST /partners) e inicia processo KYP.
3. KYP é realizado (manual/automático) e resultado é persistido (POST /partners/{id}/kyc).
4. Se aprovado, o Portal de Relacionamento solicita criação de contas no Portal Dev e fornece credenciais iniciais (POST /partners/{id}/accounts).
5. Notificação por e-mail/portal e dashboard com próximos passos (upload de app, documentação requerida).

### 2) Cadastro e publicação inicial de um aplicativo (parceiro)

1. Parceiro cadastra o aplicativo no Portal de Relacionamento (POST /applications) informando metadados e upload do artefato.
2. Portal de Relacionamento cria o registro e notifica o Portal Certificação (POST /certifications/submit) — pode encaminhar o pacote ou apenas o link/artefato.
3. Portal Certificação executa testes e retorna relatórios; em caso de reprovação, gera checklist de correções.
4. Após aprovação, o Portal Certificação atualiza o status e emite selo; a Loja de Aplicativos passa a mostrar o app com selo de conformidade.
5. Se for app privado, o Portal de Relacionamento solicita ao MDM disponibilização para grupos/pilotos (POST /mdm/apps e POST /mdm/groups/{id}/deploy).

### 3) Distribuição para Piloto e Produção

1. Parceiro ou time de produto solicita liberação de piloto via Portal de Relacionamento (POST /applications/{id}/release?env=pilot).
2. Portal de Relacionamento verifica certificação e autorizações e cria uma release para o MDM, com target groups.
3. MDM realiza deploy para o grupo de pilotos e reporta status de instalação.
4. Após resultado do piloto, equipe solicita promoção para produção — processo similar, com escopos e aprovação manual se necessário.

### 4) Atualização contínua (versões)

1. Parceiro submete nova versão no Portal de Relacionamento.
2. Processo de certificação pode ser disparado automaticamente com filtro de risco (por exemplo, breaking API changes disparam testes completos).
3. Se aprovado, novas versões são publicadas na Loja e enviadas ao MDM para atualização controlada.

## Interações técnicas e contratos mínimos de API

- API de eventos/assincronismo: os portais devem oferecer webhooks ou integração via mensageria (por exemplo, eventos: partner.created, app.submitted, certification.finished, release.approved, mdm.deploy.status).
- Segurança: OAuth2 + JWT para APIs; scopes por recurso (partners:read, apps:write, certifications:read).
- Idempotência e versionamento: endpoints de submissão devem ser idempotentes (retry-safe). APIs públicas versionadas via URL (/v1/...).

## Next steps

- Detalhar contratos de API (spec OpenAPI) para cada serviço.
- Modelar data flow (diagrama) entre portais e MDM/Loja.
- Criar jornadas técnicas com exemplos de payloads e erros comuns.

---

Versão: rascunho inicial — sinta-se livre para pedir alterações, adicionar requisitos ou jornadas adicionais.
