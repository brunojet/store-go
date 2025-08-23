# 📋 Relatório Comparativo – Jornada de Certificação de Aplicativos

## 🧠 Contexto

O processo de **cadastro do aplicativo e suas configurações** ocorre em uma jornada separada. Nesta etapa de certificação, o parceiro **não deve configurar nada manualmente**, apenas **submeter o APK** para validação contra configurações previamente aprovadas.



---

## 🧩 Cenário 1 – Mais próximo do AsIs

### 🔹 Etapas:
1. Escolhe aplicativo
2. Escolhe modelo
3. Escolhe integração
4. Seleciona APK no PC
5. Valida package name
6. Valida nome e código de versão

### ✅ Vantagens:
- Controle total do parceiro sobre cada etapa.
- Familiaridade com o fluxo atual.
- Permite testar diferentes combinações.

### ❌ Desvantagens:
- Exige mais passos e decisões manuais.
- Maior risco de erro humano.
- Ignora configurações previamente aprovadas.
- Menor automação e consistência.

---

## ⚡ Cenário 2 – Menos passos

### 🔹 Etapas:
1. Seleciona APK no PC
2. Busca configurações pelo package name
3. Valida nome e código de versão
4. Exibe hint do aplicativo
5. Lista configurações vinculadas

### ✅ Vantagens:
- Fluxo rápido e automatizado.
- Reduz drasticamente o risco de erro humano.
- Garante alinhamento com configurações aprovadas.
- Ideal para parceiros que já têm APKs prontos.

### ❌ Desvantagens:
- Depende da existência de configurações vinculadas ao package name.
- Menor controle manual (embora não necessário neste contexto).

---

## 🧠 Cenário 3 – Escolhe configuração antes do APK

### 🔹 Etapas:
1. Escolhe aplicativo
2. Escolhe configuração aprovada
3. Seleciona APK no PC
4. Valida package name, nome e código de versão

### ✅ Vantagens:
- Permite ao parceiro escolher qual configuração deseja certificar.
- Útil em casos com múltiplas configurações por aplicativo.
- Validação robusta contra configuração escolhida.

### ❌ Desvantagens:
- Reintroduz risco de erro humano na escolha da configuração.
- Menos ágil que o cenário 2.
- Exige atenção do parceiro para não selecionar configuração incompatível.

---

## 🆚 Comparativo Geral

| Critério                     | Cenário 1                         | Cenário 2                         | Cenário 3                         |
|-----------------------------|-----------------------------------|-----------------------------------|-----------------------------------|
| Facilidade de uso           | Baixa                             | Muito alta                        | Alta                              |
| Risco de erro do usuário    | Alto                              | Muito baixo                       | Médio                             |
| Agilidade                   | Baixa                             | Máxima                            | Boa                               |
| Controle do parceiro        | Alto                              | Baixo (mas não necessário)        | Médio                             |
| Alinhamento com configurações aprovadas | Baixo              | Alto                              | Alto                              |
| Robustez da validação       | Média                             | Alta                              | Alta                              |

---

## ✅ Veredicto Final

Considerando que o parceiro **não deve configurar manualmente** e que as configurações já foram **aprovadas previamente**, o **Cenário 2 é o mais recomendado** para a jornada de certificação:

- Elimina escolhas manuais.
- Minimiza erros.
- Maximiza agilidade.
- Garante consistência com o processo de aprovação.

🔍 O **Cenário 3 pode ser útil** em casos específicos com múltiplas configurações por app, mas exige atenção extra do parceiro e não oferece os mesmos ganhos de automação.

## Exemplo Portal do Parceiro Dev
<br>

### 📤 Enviar APK para Certificação (botão)
---

### 📱 Aplicativo: App A

| Modelo POS | Package Name               | Configuração | Status         |
|------------|----------------------------|--------------|----------------|
| POS 1      | com.parceiro.appA.pos1     | A1           | 🟡 Em validação |
| POS 2      | com.parceiro.appA.pos2     | A2           | 🟢 Certificado  |

---

### 📱 Aplicativo: App B

| Modelo POS | Package Name               | Configuração | Status         |
|------------|----------------------------|--------------|----------------|
| POS 3      | com.parceiro.appB     | B1           | 🔴 Reprovado    |
| POS 4      | com.parceiro.appB     | B2           | 🟡 Em validação |

---

### 📱 Aplicativo: App C

| Modelo POS | Package Name               | Configuração | Status         |
|------------|----------------------------|--------------|----------------|
| POS 5      | com.parceiro.appC.pos5     | C1           | 🟢 Certificado  |

---

## 🧠 Observações

- O botão de envio é único e centralizado, pois o sistema identifica automaticamente o destino correto com base no APK.
- Cada modelo de POS possui seu próprio `package name`, permitindo validações específicas.
- Os status ajudam o parceiro a acompanhar o progresso de cada certificação.

---

## Pontos de Fricção por Cenário

### Cenário 1
- Repetição de dados a cada nova configuração.
- Retrabalho: erro no cadastro só aparece na certificação.
- Processo longo, mais sujeito a erros humanos.

### Cenário 2
- Travamento se esquecer de cadastrar configuração.
- Rigidez: menos flexível para casos especiais.
- Cadastro inicial errado afeta certificações futuras.

### Cenário 3
- Risco de escolha errada por lista aberta.
- Exige conhecimento técnico para selecionar corretamente.
- Falta de validação automática pode aprovar apps incompatíveis.
