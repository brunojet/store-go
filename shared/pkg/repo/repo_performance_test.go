package repo

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/brunojet/store-go/shared/pkg/domain"
	"github.com/google/uuid"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Carrega um arquivo JSON de shared/testdata dado apenas o nome do arquivo
func LoadJSONFromTestdata(filename string) ([]byte, error) {
	baseDir := "c:/Projects/store-go/shared/testdata"
	absPath := filepath.Join(baseDir, filename)
	return os.ReadFile(absPath)
}

// Struct para acumular métricas globais
type MetricasExecucao struct {
	Tempos       []int64
	TotalItems   int
	TotalQueries int
	mu           sync.Mutex
}

func (m *MetricasExecucao) Registrar(tempos []int64, totalItems, totalQueries int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Tempos = append(m.Tempos, tempos...)
	m.TotalItems += totalItems
	m.TotalQueries += totalQueries
}

func (m *MetricasExecucao) CalcularMetricas(logPrefix string) {
	if m.TotalQueries == 0 || len(m.Tempos) == 0 {
		fmt.Printf("[MÉTRICA GLOBAL] %s: nenhuma consulta registrada\n", logPrefix)
		return
	}
	best, worst := m.Tempos[0], m.Tempos[0]
	for _, t := range m.Tempos {
		if t < best {
			best = t
		}
		if t > worst {
			worst = t
		}
	}
	sorted := make([]int64, len(m.Tempos))
	copy(sorted, m.Tempos)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	p := func(percent float64) int64 {
		idx := int(float64(len(sorted)) * percent)
		if idx >= len(sorted) {
			idx = len(sorted) - 1
		}
		return sorted[idx]
	}
	fmt.Printf("[MÉTRICA GLOBAL] %s: %d consultas, %d itens exportados, tempo total %d ms, média %.2f ms, melhor %d ms, pior %d ms, p90 %d ms, p95 %d ms, p99 %d ms\n",
		logPrefix, m.TotalQueries, m.TotalItems, soma(m.Tempos), float64(soma(m.Tempos))/float64(m.TotalQueries), best, worst, p(0.90), p(0.95), p(0.99))
}

func soma(arr []int64) int64 {
	var s int64
	for _, v := range arr {
		s += v
	}
	return s
}

var metricasGlobal = &MetricasExecucao{}

func GetMetricasGlobal() *MetricasExecucao {
	return metricasGlobal
}

// Função auxiliar para cadastrar tipo de categoria se não existir
// Estrutura para ler o novo JSON de regiões
// Estrutura genérica para ler JSON de categories
// Função genérica recursiva para cadastrar tipos e categories
// Função genérica para carregar e cadastrar categories a partir de um JSON
// Cadastra modelos de terminal a partir do JSON
func CadastroModelosTerminalFromJSON(db *gorm.DB) error {
	data, err := LoadJSONFromTestdata("modelos_terminal.json")
	if err != nil {
		return err
	}
	var entidades []domain.BaseEntity
	if err := json.Unmarshal(data, &entidades); err != nil {
		return err
	}
	for _, be := range entidades {
		modelo := domain.TerminalModel{BaseEntity: be}
		var existing domain.TerminalModel
		if err := db.Where("name = ?", modelo.Name).First(&existing).Error; err != nil {
			db.Create(&modelo)
		}
	}
	return nil
}

// Cadastra tipos de integração a partir do JSON
func CadastroTiposIntegracaoFromJSON(db *gorm.DB) error {
	data, err := LoadJSONFromTestdata("tipos_integracao.json")
	if err != nil {
		return err
	}
	var entidades []domain.BaseEntity
	if err := json.Unmarshal(data, &entidades); err != nil {
		return err
	}
	for _, be := range entidades {
		tipo := domain.IntegrationType{BaseEntity: be}
		var existing domain.IntegrationType
		if err := db.Where("name = ?", tipo.Name).First(&existing).Error; err != nil {
			db.Create(&tipo)
		}
	}
	return nil
}

func CadastroCategoriesFromJSONPath(db *gorm.DB, relativePath string) error {
	baseDir := "c:/Projects/store-go/shared/testdata"
	absPath := filepath.Join(baseDir, relativePath)
	file, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var root CategoriesRoot
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&root); err != nil {
		return err
	}
	return cadastroCategoriesFromJSONRec(db, root, nil)
}
func cadastroCategoriesFromJSONRec(db *gorm.DB, root CategoriesRoot, idPai *int64) error {
	// Cadastra o tipo de categoria
	tipoCat, err := ensureCategoryType(db, root.CategoryType)
	if err != nil {
		return err
	}
	for nome, raw := range root.Categories {
		// Cadastra a categoria
		categoria := domain.Category{
			BaseEntity:     domain.BaseEntity{Name: nome, Active: true},
			CategoryTypeId: tipoCat.ID,
			ParentId:       idPai,
		}
		var existing domain.Category
		err := db.Where("name = ? AND category_type_id = ?", nome, tipoCat.ID).First(&existing).Error
		if err != nil {
			if err := db.Create(&categoria).Error; err != nil {
				return err
			}
		}
		// Verifica se há subcategories
		var subcat CategoriesRoot
		if err := json.Unmarshal(raw, &subcat); err == nil && len(subcat.Categories) > 0 {
			// Busca o ID da categoria recém-cadastrada
			if err := cadastroCategoriesFromJSONRec(db, subcat, &categoria.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

type CategoriesRoot struct {
	CategoryType string                     `json:"category_type"`
	Categories   map[string]json.RawMessage `json:"categories"`
}

func ensureCategoryType(db *gorm.DB, nome string) (domain.CategoryType, error) {
	var tipo domain.CategoryType
	if err := db.Where("name = ?", nome).First(&tipo).Error; err != nil {
		tipo = domain.CategoryType{BaseEntity: domain.BaseEntity{Name: nome, Active: true}}
		if err := db.Create(&tipo).Error; err != nil {
			return tipo, err
		}
	}
	return tipo, nil
}

// Estrutura para ler o JSON de regiões
type Regiao struct {
	Regiao string `json:"regiao"`
}

// Estrutura para ler o JSON de ramos e subramos
type RamoSubramo struct {
	Tipo     string   `json:"tipo"`
	Ramo     string   `json:"ramo"`
	Subramos []string `json:"subramos"`
}

func CadastroApplicationsFromJSON(db *gorm.DB) {
	data, err := LoadJSONFromTestdata("aplicativos.json")
	if err != nil {
		fmt.Printf("Erro ao ler aplicativos.json: %v\n", err)
		return
	}
	var apps map[string]map[string]interface{}
	if err := json.Unmarshal(data, &apps); err != nil {
		fmt.Printf("Erro ao decodificar aplicativos.json: %v\n", err)
		return
	}
	for name := range apps {
		var existing domain.Application
		if err := db.Where("name = ?", name).First(&existing).Error; err != nil {
			app := domain.Application{BaseEntity: domain.BaseEntity{Name: name, Active: true}}
			db.Create(&app)
		}
	}
}

func CadastroConfiguraces(db *gorm.DB) {
	// Buscar todos os modelos de terminal
	var modelos []domain.TerminalModel
	db.Find(&modelos)

	// Buscar todos os tipos de integração
	var integracoes []domain.IntegrationType
	db.Find(&integracoes)

	// Buscar todos os aplicativos
	var apps []domain.Application
	db.Find(&apps)

	for _, app := range apps {
		for _, modelo := range modelos {
			for _, integracao := range integracoes {
				// Verifica se já existe configuração
				var existing domain.ApplicationConfiguration
				db.Where("application_id = ? AND terminal_model_id = ? AND integration_type_id = ?", app.ID, modelo.ID, integracao.ID).First(&existing)
				if existing.ApplicationId == 0 {
					config := domain.ApplicationConfiguration{
						ApplicationId:     app.ID,
						TerminalModelId:   modelo.ID,
						IntegrationTypeId: integracao.ID,
					}
					db.Create(&config)
				}
			}
		}
	}
}

func updateCadastroCategory(db *gorm.DB, categories []domain.Category, baseModel domain.BaseModel) error {
	// Associa múltiplas categories ao cadastro via GORM M2M
	cadastro := domain.ApplicationProfileHistory{BaseModel: baseModel}
	if err := db.Model(&cadastro).Association("Categories").Append(&categories); err != nil {
		return err
	}
	return nil
}

func associarCategoriesAoCadastro(db *gorm.DB, categories []domain.Category, nomeApp string, baseModel domain.BaseModel) error {
	// Busca nome do app no cadastro
	var cadastro domain.ApplicationProfileHistory
	if err := db.Where("id = ?", baseModel.ID).First(&cadastro).Error; err != nil {
		return err
	}

	data, err := LoadJSONFromTestdata("aplicativos.json")
	if err != nil {
		return err
	}
	var apps map[string]map[string]interface{}
	if err := json.Unmarshal(data, &apps); err != nil {
		return err
	}
	appInfo, ok := apps[nomeApp]
	if !ok {
		return nil
	}
	// Acumula categories a associar
	var catsToAssociate []domain.Category
	// Associa regiões
	if regioes, ok := appInfo["regioes"].([]interface{}); ok {
		for _, reg := range regioes {
			regStr, ok := reg.(string)
			if !ok {
				continue
			}
			for _, cat := range categories {
				if cat.Name == regStr {
					catsToAssociate = append(catsToAssociate, cat)
				}
			}
		}
	}
	// Associa ramo
	if ramo, ok := appInfo["ramo"].(string); ok {
		for _, cat := range categories {
			if cat.Name == ramo {
				catsToAssociate = append(catsToAssociate, cat)
				break
			}
		}
	}
	// Associa subramo
	if subramo, ok := appInfo["subramo"].(string); ok {
		for _, cat := range categories {
			if cat.Name == subramo {
				catsToAssociate = append(catsToAssociate, cat)
				break
			}
		}
	}
	if len(catsToAssociate) > 0 {
		if err := updateCadastroCategory(db, catsToAssociate, baseModel); err != nil {
			return err
		}
	}
	return nil
}

func associarConfiguracaoAoCadastro(db *gorm.DB, configuracoes []domain.ApplicationConfiguration, baseModel domain.BaseModel) error {
	// Associa múltiplas configurações ao cadastro via GORM M2M
	cadastro := domain.ApplicationProfileHistory{BaseModel: baseModel}
	if err := db.Model(&cadastro).Association("ApplicationConfigurations").Append(&configuracoes); err != nil {
		return err
	}
	return nil
}

func SolicitarCadastroApplication(db *gorm.DB) error {
	var tipoRegiao, tipoRamos, tipoSubRamos domain.CategoryType
	db.Where("name = ?", "Região").First(&tipoRegiao)
	db.Where("name = ?", "Ramos").First(&tipoRamos)
	db.Where("name = ?", "Subramos").First(&tipoSubRamos)
	var todasCategories []domain.Category
	db.Where("category_type_id IN ?", []int64{tipoRegiao.ID, tipoRamos.ID, tipoSubRamos.ID}).Find(&todasCategories)

	var apps []domain.Application
	db.Find(&apps)

	for _, app := range apps {
		// Buscar todas as configurações do app
		var configuracoes []domain.ApplicationConfiguration
		db.Preload("TerminalModel", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).Preload("IntegrationType", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).Where("application_id = ?", app.ID).Find(&configuracoes)

		cadastroConfig := domain.ApplicationProfileHistory{
			ApplicationContact: &domain.ApplicationContact{
				BaseEntity: domain.BaseEntity{
					Name: app.Name + " S.A.",
				},
				Email: "ApplicationContact@" + app.Name + ".com",
				Phone: "1234-5678",
				Site:  "www." + app.Name + ".com",
			},
			ApplicationDetail: &domain.ApplicationDetail{
				Description: fmt.Sprintf(
					"Cadastro para %s",
					app.Name),
			},
		}
		if tx := db.Create(&cadastroConfig); tx.Error != nil {
			return tx.Error
		}
		if err := associarCategoriesAoCadastro(db, todasCategories, app.Name, cadastroConfig.BaseModel); err != nil {
			return err
		}

		if len(configuracoes) > 0 {
			if err := associarConfiguracaoAoCadastro(db, configuracoes, cadastroConfig.BaseModel); err != nil {
				return err
			}
		}
	}

	return nil
}

func SolicitarCadastroApplicationVersionHistory(db *gorm.DB) error {
	var apps []domain.Application
	if err := db.Find(&apps).Error; err != nil {
		return err
	}

	for _, app := range apps {
		// Buscar o cadastro mais recente para o app
		cadastro, err := getCadastroApplication(db, app.ID)

		if err != nil || cadastro.ID == 0 {
			continue
		}

		var versoes []domain.ApplicationVersionHistory
		for _, config := range cadastro.ApplicationConfigurations {
			now := time.Now()
			year := now.Year() % 100
			dayOfYear := now.YearDay()
			seconds := now.Hour()*3600 + now.Minute()*60 + now.Second()
			nomeVersao := fmt.Sprintf("%02d%03d%05d", year, dayOfYear, seconds)
			Size := int64(rand.Intn(190)+10) * 1024 * 1024 // 10MB a 200MB

			versao := domain.ApplicationVersionHistory{
				ApplicationId:     app.ID,
				IntegrationTypeId: config.IntegrationTypeId,
				TerminalModelId:   config.TerminalModelId,
				BaseEntity: domain.BaseEntity{
					Name: app.Name,
				},
				Size:        Size,
				VersionName: nomeVersao,
				Image: &domain.Image{
					StorageObject: domain.StorageObject{
						Path:     uuid.NewString(),
						Name:     fmt.Sprintf("Imagem %s de %s", nomeVersao, app.Name),
						MimeType: "image/png",
					},
				},
			}
			versoes = append(versoes, versao)
		}
		if len(versoes) > 0 {
			if err := db.Create(&versoes).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// getCadastroApplication localiza o cadastro (ApplicationProfileHistory) mais recente
// para um determinado application_id. Estratégia:
//   - Faz JOIN direto na tabela M2M (application_profile_history_configuration) para
//     filtrar por application_id.
//   - Ordena por application_profile_history.id DESC (pega o profile mais novo por ID).
//   - Preload("ApplicationConfigurations") para já retornar o cadastro com as
//     configurações necessárias para os passos seguintes.
func getCadastroApplication(db *gorm.DB, appID int64) (domain.ApplicationProfileHistory, error) {
	var profile domain.ApplicationProfileHistory
	if err := db.Preload("ApplicationConfigurations").
		Model(&domain.ApplicationProfileHistory{}).
		Joins("JOIN application_profile_history_configuration apc ON apc.profile_id = application_profile_history.id").
		Where("apc.application_id = ?", appID).
		Order("application_profile_history.id DESC").
		First(&profile).Error; err != nil {
		return domain.ApplicationProfileHistory{}, err
	}
	return profile, nil
}

func SolicitarPublicacaoApplicationVersionHistory(db *gorm.DB) error {
	var apps []domain.Application
	if err := db.Find(&apps).Error; err != nil {
		return err
	}
	stages := []domain.Stage{
		domain.StageDevelopment,
		domain.StageTesting,
		domain.StagePilot,
		domain.StageProduction,
	}

	for _, app := range apps {
		cadastro, err := getCadastroApplication(db, app.ID)
		if err != nil || cadastro.ID == 0 {
			continue
		}

		// Para cada configuração do cadastro, buscar a versão mais recente e atualizar o catálogo
		for _, config := range cadastro.ApplicationConfigurations {
			var versaoAtual domain.ApplicationVersionHistory
			if err := db.Where("application_id = ? AND terminal_model_id = ? AND integration_type_id = ?",
				config.ApplicationId, config.TerminalModelId, config.IntegrationTypeId).
				Order("id desc").First(&versaoAtual).Error; err != nil || versaoAtual.ID == 0 {
				continue
			}

			var catalogs []domain.ApplicationCatalog
			db.Where("application_id = ? AND integration_type_id = ? AND terminal_model_id = ?", config.ApplicationId, config.IntegrationTypeId, config.TerminalModelId).Order("stage DESC").Find(&catalogs)

			lenCatalogs := len(catalogs)
			lenStages := len(stages)

			// Cria novo catálogo se ainda não percorreu todos os estágios
			if lenCatalogs < lenStages {
				newCatalog := domain.ApplicationCatalog{
					ApplicationId:        config.ApplicationId,
					IntegrationTypeId:    config.IntegrationTypeId,
					TerminalModelId:      config.TerminalModelId,
					Stage:                stages[lenCatalogs],
					ApplicationProfileId: &cadastro.ID,
					ApplicationVersionId: &versaoAtual.ID,
				}
				// Idempotente: evita erro de PK duplicada caso já exista o registro para este estágio
				if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&newCatalog).Error; err != nil {
					return err
				}
				// Garante que teremos o registro em memória (se já existia, buscamos)
				var created domain.ApplicationCatalog
				if err := db.Where("application_id = ? AND integration_type_id = ? AND terminal_model_id = ? AND stage = ?",
					newCatalog.ApplicationId, newCatalog.IntegrationTypeId, newCatalog.TerminalModelId, newCatalog.Stage,
				).First(&created).Error; err != nil {
					return err
				}
				catalogs = append([]domain.ApplicationCatalog{created}, catalogs...)
				lenCatalogs++
			}

			// Avança os estágios das versões
			for i := 0; lenCatalogs > 0 && i < lenCatalogs; i++ {
				catalogAtual := catalogs[i]
				if i < lenCatalogs-1 {
					catalogAnterior := catalogs[i+1]
					catalogAtual.ApplicationVersionId = catalogAnterior.ApplicationVersionId
					// Se for produção, atualiza o profile
					if catalogAtual.Stage == domain.StageProduction {
						catalogAtual.ApplicationProfileId = &cadastro.ID
					}
				} else {
					catalogAtual.ApplicationVersionId = &versaoAtual.ID
				}
				// Nota: anteriormente evitamos Save devido ao Stage=0 (valor zero em PK composta)
				// que levava o GORM a inferir INSERT. Como o enum Stage não inicia mais em 0,
				// Save volta a ser seguro aqui e executará UPDATE corretamente.
				if err := db.Save(&catalogAtual).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Cria índices recomendados para melhorar a performance das consultas utilizadas nos testes
func EnsurePerfIndexes(db *gorm.DB) error {
	stmts := []string{
		// M2M cadastro-categoria: acelerar filtro por categoria -> perfil
		"CREATE INDEX IF NOT EXISTS idx_aphc_category_profile ON application_profile_history_category (category_id, profile_id)",
		// Catálogo: acelerar junção pelo profile
		"CREATE INDEX IF NOT EXISTS idx_ac_application_profile ON application_catalog (application_profile_id)",
		// Opcional/composto: acelerar filtros por integração + modelo + profile
		"CREATE INDEX IF NOT EXISTS idx_ac_it_tm_profile ON application_catalog (integration_type_id, terminal_model_id, application_profile_id)",
	}
	for _, sql := range stmts {
		if err := db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

func gerarJsonAppsPorRegiao(db *gorm.DB) error {
	consultaVersaoApp := func(db *gorm.DB, regiaoID, modeloID, integracaoID int64) ([]map[string]interface{}, error) {
		var result []map[string]interface{}
		err := db.Table("application_version").
			Select(`application_version.name as nome_app, application_version.Size, application_version.id as id_versao, mdl_trml.name as modelo_terminal_nome, integration_type.name as tipo_integracao_nome, category.name as categoria_nome`).
			Joins("JOIN application_catalog ON application_catalog.application_version_id = application_version.id").                                                   //Localiza a versao no catalogo
			Joins("JOIN cfg_aplv ON id = application_version.id_cfg_aplv").                                                                                             //Localiza a configuracao
			Joins("JOIN application_profile_history ON application_profile_history.id = application_catalog.id_application_profile_history").                           //Localiza o cadastro
			Joins("JOIN category_application_profile_history ON category_application_profile_history.id_application_profile_history = application_profile_history.id"). //M2M cadastro-categoria
			Joins("JOIN category ON category.id = category_application_profile_history.id_category").                                                                   //Localiza a categoria
			Joins("JOIN mdl_trml ON mdl_trml.id = terminal_model_id").                                                                                                  //Localiza o modelo do terminal
			Joins("JOIN integration_type ON integration_type.id = integration_type_id").                                                                                //Localiza o tipo de integração
			Where("terminal_model_id = ?", modeloID).
			Where("integration_type_id = ?", integracaoID).
			Where("category.id = ?", regiaoID).
			Order("application_version.created_at DESC").
			Scan(&result).Error
		return result, err
	}
	metricas := GetMetricasGlobal()
	return gerarJsonAppsGenerico(db, consultaVersaoApp, "apps_regiao_%s_%s_%s.json", "gerarJsonAppsPorRegiao", metricas)
}

// Exporta apps por região usando CatalogoApplication
func gerarJsonAppsPorRegiaoCatalogo(db *gorm.DB) error {
	consultaCatalogoApp := func(db *gorm.DB, regiaoID, modeloID, integracaoID int64) ([]map[string]interface{}, error) {
		var result []map[string]interface{}
		err := db.Table("application_catalog ac").
			Select(`
				av.name  as nome_app,
				av.size  as size,
				av.id    as id_versao,
				tm.name  as modelo_terminal_nome,
				it.name  as tipo_integracao_nome,
				cat.name as categoria_nome
			`).
			Joins("JOIN application_version av ON av.id = ac.application_version_id").
			Joins("JOIN application_profile_history aph ON aph.id = ac.application_profile_id").
			Joins("JOIN application_profile_history_category aphc ON aphc.profile_id = aph.id").
			Joins("JOIN category cat ON cat.id = aphc.category_id").
			Joins("JOIN terminal_model tm ON tm.id = ac.terminal_model_id").
			Joins("JOIN integration_type it ON it.id = ac.integration_type_id").
			Where("ac.terminal_model_id = ?", modeloID).
			Where("ac.integration_type_id = ?", integracaoID).
			Where("cat.id = ?", regiaoID).
			Order("av.created_at DESC").
			Scan(&result).Error
		return result, err
	}
	metricas := GetMetricasGlobal()
	return gerarJsonAppsGenerico(db, consultaCatalogoApp, "apps_regiao_catalogo_%s_%s_%s.json", "gerarJsonAppsPorRegiaoCatalogo", metricas)
}

// Função genérica para exportar apps por região, modelo e integração
func gerarJsonAppsGenerico(
	db *gorm.DB,
	consulta func(db *gorm.DB, regiaoID, modeloID, integracaoID int64) ([]map[string]interface{}, error),
	fileNameFmt string,
	logPrefix string,
	metricas *MetricasExecucao,
) error {
	var idsRegioes []int64
	db.Model(&domain.Category{}).
		Joins("JOIN category_type ON category_type.id = category.category_type_id").
		Where("category_type.name = ?", "Região").
		Pluck("category.id", &idsRegioes)

	var modelos []domain.TerminalModel
	db.Find(&modelos)
	var integracoes []domain.IntegrationType
	db.Find(&integracoes)

	var totalItems int
	var totalQueries int
	times := make([]int64, 0)
	for _, regiaoID := range idsRegioes {
		var regiao domain.Category
		db.First(&regiao, regiaoID)
		for _, modelo := range modelos {
			for _, integracao := range integracoes {
				start := time.Now()
				apps, err := consulta(db, regiaoID, modelo.ID, integracao.ID)
				if err != nil {
					return err
				}
				fileName := fmt.Sprintf(fileNameFmt, regiao.Name, modelo.Name, integracao.Name)
				file, err := os.Create(fileName)
				if err != nil {
					return err
				}
				defer file.Close()
				encoder := json.NewEncoder(file)
				encoder.SetIndent("", "  ")
				if err := encoder.Encode(apps); err != nil {
					return err
				}
				duracao := time.Since(start).Milliseconds()
				totalItems += len(apps)
				totalQueries++
				times = append(times, duracao)
				//fmt.Printf("%s %s, Modelo %s, Integracao %s: %d itens, consulta e exportacao demorou %d ms\n", logPrefix, regiao.Name, modelo.Name, integracao.Name, len(apps), duracao)
			}
		}
	}
	if metricas != nil {
		metricas.Registrar(times, totalItems, totalQueries)
	}
	return nil
}

func TestEntrypoint(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable search_path=store_go"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Erro ao inicializar o banco: %v", err)
	}
	// Garante migrações e índices via PostMigrate das entidades de domínio
	if err := domain.AutoMigrate(db); err != nil {
		t.Fatalf("Erro ao executar AutoMigrate/PostMigrate: %v", err)
	}
	setup := false
	if setup {
		CadastroCategoriesFromJSONPath(db, "regioes.json")
		CadastroCategoriesFromJSONPath(db, "ramos_subramos.json")
		CadastroModelosTerminalFromJSON(db)
		CadastroTiposIntegracaoFromJSON(db)
		CadastroApplicationsFromJSON(db)
		CadastroConfiguraces(db)
	}

	loadVersao := false

	if loadVersao {
		for i := 0; i < 50; i++ {
			SolicitarCadastroApplication(db)
			SolicitarCadastroApplicationVersionHistory(db)
			SolicitarPublicacaoApplicationVersionHistory(db)
		}
	}

	performanceTest := true
	if performanceTest {
		var wg sync.WaitGroup
		metricasGlobal = &MetricasExecucao{} // resetar métricas
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 20; j++ {
					gerarJsonAppsPorRegiaoCatalogo(db)
					//gerarJsonAppsPorRegiao(db)
					time.Sleep(100 * time.Millisecond)
				}
			}()
		}
		wg.Wait()
		metricasGlobal.CalcularMetricas("gerarJsonAppsPorRegiao")
	}

}
