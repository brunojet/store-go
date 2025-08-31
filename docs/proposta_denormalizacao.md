# Proposta de Denormalização do Modelo de Dados

## Objetivo
Simplificar o modelo de dados, facilitar consultas e reduzir complexidade para o time de dados, mantendo governança e rastreabilidade.

## Pontos de Denormalização Propostos

### 1. Perfil do Aplicativo
- Incluir dados de contato (site, email, telefone, responsável) e detalhes do aplicativo (descrição, observações) diretamente na tabela de perfil (`HistoricoPerfilAplicativo`).
- **Vantagens:** O perfil torna-se autoexplicativo, facilitando consultas, análises e relatórios, além de reduzir a dependência de múltiplas tabelas apenas para chaves estrangeiras.
- **Ônus:** Qualquer alteração nos dados de contato ou detalhes exige a replicação ou atualização da linha inteira do perfil, podendo gerar duplicidade se houver reuso de informações.
- **Frequência de alterações:** Baixa, o que minimiza o impacto da replicação.
- **Governança:** Normalmente, apenas dois registros de perfil estarão ativos simultaneamente (revisão e produção). Perfis anteriores podem ser descartados conforme regras de expurgo, mantendo o modelo enxuto e auditável.

Essa abordagem equilibra simplicidade, clareza e governança, sendo especialmente indicada para cenários onde os dados são pouco mutáveis e o objetivo é facilitar o trabalho do time de dados.

### 2. Estágio como Enum Constraint
- Substituir a tabela de estágios por um campo inteiro com constraint (ex: 10 = desenvolvimento, 20 = certificação, 30 = revisão, 40 = piloto, 50 = produção).
- Vantagem: Reduz número de tabelas e joins, facilita ordenação e consultas, permite evoluir estágios facilmente.
- Ônus: Menos flexível para adicionar metadados extras (descrição, cor, ícone).

### 3. Dados Resumidos no Catálogo
- **Decisão:** Optou-se por manter o catálogo normalizado, sem replicar dados resumidos de aplicativo, versão ou perfil diretamente na tabela de catálogo.
- **Justificativa:** Apesar da denormalização facilitar consultas e relatórios, ela pode gerar duplicidade, inconsistência e dificultar manutenção, especialmente se os dados mudam com frequência ou precisam ser atualizados em vários lugares.
- **Vantagens da normalização:** Garante integridade, evita duplicidade e facilita manutenção, além de ser adequada ao perfil de uso e capacidade do time de dados. O volume de dados atual e a política de expurgo permitem consultas eficientes mesmo com modelo normalizado.
- **Revisão futura:** Caso o perfil de uso mude ou o volume de consultas exija mais performance, a estratégia pode ser reavaliada.

### 4. Histórico com Snapshots
- Armazenar snapshots de dados relevantes (nome, descrição, status) em tabelas de histórico, para facilitar auditoria e consultas temporais.
- Vantagem: Rastreabilidade e auditoria facilitadas.
- Ônus: Crescimento do volume de dados.

## Critérios para Denormalização
- Dados pouco mutáveis ou exclusivos de cada registro.
- Objetivo de simplificar consultas e relatórios.
- Risco de inconsistência baixo ou aceitável.
- Política de retenção e auditoria bem definida.

## Considerações Finais
- Denormalização deve ser documentada e monitorada.
- Avaliar impacto em manutenção, performance e governança.
- Caso o modelo evolua para exigir mais flexibilidade, pode-se migrar para normalização sem grandes impactos.

---
Mantenha este documento como referência para decisões futuras de modelagem e evolução do sistema.
