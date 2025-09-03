# Agregação Avançada em SQL e ORMs

## 1. Agregação de Dados Relacionados
 Consultas SQL podem retornar múltiplas linhas para o mesmo registro quando há relacionamentos N:N ou 1:N (ex: catálogo com várias categorias).
 Para evitar duplicidade, use funções agregadoras:
  - **Postgres:** `array_agg(campo)` para arrays reais.
  - **MySQL:** `GROUP_CONCAT(campo)` para strings separadas por vírgula.

## 1.1. Agregação Condicional por Tipo de Categoria
### Postgres
```sql
SELECT
  nome_app,
  array_agg(ctgr.nome) FILTER (WHERE ctgr.tipo = 'A') AS categorias_tipo_a,
  array_agg(ctgr.nome) FILTER (WHERE ctgr.tipo = 'B') AS categorias_tipo_b
FROM ...
GROUP BY nome_app;
```

### MySQL
```sql
SELECT
  nome_app,
  GROUP_CONCAT(CASE WHEN ctgr.tipo = 'A' THEN ctgr.nome END) AS categorias_tipo_a,
  GROUP_CONCAT(CASE WHEN ctgr.tipo = 'B' THEN ctgr.nome END) AS categorias_tipo_b
FROM ...
GROUP BY nome_app;
```

## 2. Exemplo de Consulta Agregada
```sql
-- Postgres
SELECT nome, array_agg(categoria) AS categorias
FROM produto
JOIN categoria_produto ON ...
GROUP BY nome;

-- MySQL
SELECT nome, GROUP_CONCAT(categoria) AS categorias
FROM produto
JOIN categoria_produto ON ...
GROUP BY nome;
```

## 3. ORMs e Limitações
- ORMs como GORM facilitam CRUD, joins e filtros simples.
- Agregações avançadas (arrays, concatenação) geralmente exigem SQL manual (`db.Raw`).
- Motivo: cada banco tem funções agregadoras diferentes, dificultando abstração universal.

## 4. Modelos Customizados para API
- Para expor dados agregados em JSON, crie structs customizados:
```go
type ProdutoView struct {
    Nome       string   `json:"nome"`
    Categorias []string `json:"categorias"`
}
```

## 5. Recomendações
- Use SQL manual para agregações avançadas.
- Mapeie o resultado para structs customizados.
- Documente consultas e mantenha o código organizado.

---
**Resumo:**
- Agregação avançada é poderosa, mas depende do banco.
- ORMs não cobrem todos os casos; SQL manual é comum e seguro se bem usado.
- Estruture o modelo de resposta para facilitar integração com frontend.
