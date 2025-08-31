# Guia de Desenvolvimento do Chassis da Loja

Este documento tem como objetivo orientar os desenvolvedores para garantir a integridade, flexibilidade e governança do núcleo (chassis) do sistema de loja.

## Princípios do Chassis
- Não acople regras de negócio diretamente nas entidades básicas.
- Mantenha nomes, relacionamentos e padrões conforme a documentação oficial.
- Toda alteração estrutural deve ser revisada e documentada.

## Boas Práticas
- Utilize migrations e testes automatizados para validar mudanças no núcleo.
- Documente novos campos, entidades e relacionamentos imediatamente.
- Prefira composições e extensões ao invés de modificações diretas nas entidades básicas.

## Processos de Revisão
- Toda Pull Request que altera o núcleo deve passar por revisão técnica e de arquitetura.
- Sincronize alterações entre código e documentação para evitar divergências.

## Fail-Proof
- Implemente testes de integridade para entidades básicas (ex: constraints, validações mínimas).
- Evite dependências externas no núcleo; mantenha o core isolado.

## Flexibilidade
- Estruture o núcleo para permitir plug-ins e módulos de negócio sem retrabalho.
- Use interfaces e abstrações para facilitar adaptações futuras.

## Checklist para Alterações
- [ ] Nome e relacionamento seguem o padrão?
- [ ] Documentação foi atualizada?
- [ ] Testes automatizados cobrem a alteração?
- [ ] Revisão técnica foi realizada?

---

## Avaliação do Projeto Atual

O que foi feito até agora está totalmente alinhado com estas recomendações:
- Nomenclatura e relacionamentos padronizados.
- Separação clara entre núcleo e regras de negócio.
- Documentação e código sincronizados.
- Estrutura flexível e preparada para crescimento.

Nenhuma regra foi ferida até o momento. O projeto está em conformidade com as melhores práticas para um chassis fail-proof e flexível.

---

## Exemplo de JWT para Identificação de Terminal

O JWT pode conter os seguintes claims (informações):

- `sub`: Identificador único do terminal (UUID ou ID no banco)
- `iat`: Data/hora de emissão do token
- `exp`: Data/hora de expiração do token
- `tipo_integracao_id`: ID do tipo de integração do terminal
- `modelo_terminal_id`: ID do modelo do terminal
- `estagio_id`: ID do estágio/contexto do terminal
- `perfil`: Perfil ou classificação do terminal (ex: piloto, produção)
- `app_version`: Versão do aplicativo/firmware do terminal
- Outros campos relevantes para filtros, contexto ou auditoria

### Exemplo de payload JWT (JSON)
```json
{
  "sub": "b1a2c3d4-e5f6-7890-abcd-1234567890ef",
  "iat": 1693382400,
  "exp": 1693386000,
  "tipo_integracao_id": "1",
  "modelo_terminal_id": "5",
  "estagio_id": "2",
  "perfil": "piloto",
  "app_version": "2.1.0"
}
```

Esses dados são assinados e enviados no header (ex: `Authorization: Bearer <jwt>` ou `X-Terminal-Token: <jwt>`). O backend valida a assinatura, expiração e utiliza os claims para aplicar filtros e regras nas requisições do terminal.

## Sobre o Identificador `sub` no JWT

O claim `sub` (subject) pode ser um UUIDv5 gerado de forma determinística usando os dados do terminal cadastrado como base (ex: número de série, MAC address, modelo, etc.) e um namespace fixo.

**Vantagens do UUIDv5:**
- Identificador único, padronizado e seguro para cada terminal.
- Não expõe diretamente dados sensíveis, pois é um hash dos dados-base.
- Permite reconstruir o mesmo UUID para o terminal, se necessário, usando os mesmos dados e namespace.
- Evita duplicidade de identificadores para o mesmo terminal.

**Exemplo de geração de UUIDv5:**
- Namespace: UUID fixo definido para o sistema (ex: "loja-terminals-namespace").
- Dados-base: concatenação de atributos únicos do terminal (ex: "serial:1234;mac:AA:BB:CC:DD:EE:FF;modelo:POS-X").
- O resultado é um UUIDv5 único e determinístico para cada terminal.

**Uso no JWT:**
- O campo `sub` do JWT recebe o UUIDv5 gerado.
- O backend pode validar e rastrear o terminal sem expor dados sensíveis.

## Extensão do JWT com Dados de Contexto do Terminal

O JWT pode ser facilmente estendido para incluir novos claims conforme o sistema evolui. Por exemplo, no momento da ativação do terminal, é possível consultar outros sistemas (ex: cadastro de pontos de venda) e obter dados adicionais, como:

- `ramo_atividade`: ramo de atuação do ponto de venda
- `atividades`: lista de atividades ou serviços oferecidos
- `localizacao`: localização geográfica ou endereço
- Outros dados relevantes para personalização

Esses dados podem ser incluídos como claims no JWT, permitindo que o backend utilize essas informações para:
- Personalizar e filtrar listas de aplicativos, ofertas ou funcionalidades exibidas para cada terminal
- Adaptar regras de negócio conforme o perfil do ponto de venda
- Facilitar auditoria, rastreabilidade e análise de uso

### Exemplo de payload JWT estendido
```json
{
  "sub": "b1a2c3d4-e5f6-7890-abcd-1234567890ef",
  "iat": 1693382400,
  "exp": 1693386000,
  "tipo_integracao_id": "1",
  "modelo_terminal_id": "5",
  "estagio_id": "2",
  "perfil": "piloto",
  "app_version": "2.1.0",
  "ramo_atividade": "restaurante",
  "atividades": ["delivery", "balcao"],
  "localizacao": "Rua Exemplo, 123, Cidade, UF"
}
```

O JWT, assim, carrega todo o contexto necessário para personalizar e otimizar a experiência do terminal, sem necessidade de múltiplas consultas externas a cada requisição.

---

Mantenha este guia como referência para futuras evoluções!
