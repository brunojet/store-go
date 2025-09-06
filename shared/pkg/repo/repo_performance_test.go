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

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
// Estrutura genérica para ler JSON de categorias
// Função genérica recursiva para cadastrar tipos e categorias
// Função genérica para carregar e cadastrar categorias a partir de um JSON
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
		if err := db.Where("nome = ?", modelo.Nome).First(&existing).Error; err != nil {
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
		if err := db.Where("nome = ?", tipo.Nome).First(&existing).Error; err != nil {
			db.Create(&tipo)
		}
	}
	return nil
}

func CadastroEstagiosFromJSON(db *gorm.DB) error {
	data, err := LoadJSONFromTestdata("estagios_catalogo.json")
	if err != nil {
		return err
	}
	var estagios []domain.Estagio
	if err := json.Unmarshal(data, &estagios); err != nil {
		return err
	}
	for _, estagio := range estagios {
		var existing domain.Estagio
		if err := db.Where("nome = ?", estagio.Nome).First(&existing).Error; err != nil {
			db.Create(&estagio)
		}
	}
	return nil
}

func CadastroCategorysFromJSONPath(db *gorm.DB, relativePath string) error {
	baseDir := "c:/Projects/store-go/shared/testdata"
	absPath := filepath.Join(baseDir, relativePath)
	file, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var root CategorysRoot
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&root); err != nil {
		return err
	}
	return cadastroCategorysFromJSONRec(db, root, nil)
}
func cadastroCategorysFromJSONRec(db *gorm.DB, root CategorysRoot, idPai *int64) error {
	// Cadastra o tipo de categoria
	tipoCat, err := ensureCategoryType(db, root.CategoryType)
	if err != nil {
		return err
	}
	for nome, raw := range root.Categorys {
		// Cadastra a categoria
		categoria := domain.Category{
			BaseEntity:     domain.BaseEntity{Nome: nome, Ativo: true},
			CategoryTypeId: tipoCat.ID,
			ParentId:       idPai,
		}
		var existing domain.Category
		err := db.Where("nome = ? AND category_type_id = ?", nome, tipoCat.ID).First(&existing).Error
		if err != nil {
			if err := db.Create(&categoria).Error; err != nil {
				return err
			}
		}
		// Verifica se há subcategorias
		var subcat CategorysRoot
		if err := json.Unmarshal(raw, &subcat); err == nil && len(subcat.Categorys) > 0 {
			// Busca o ID da categoria recém-cadastrada
			if err := cadastroCategorysFromJSONRec(db, subcat, &categoria.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

type CategorysRoot struct {
	CategoryType string                     `json:"tipo_categoria"`
	Categorys    map[string]json.RawMessage `json:"categorias"`
}

func ensureCategoryType(db *gorm.DB, nome string) (domain.CategoryType, error) {
	var tipo domain.CategoryType
	if err := db.Where("nome = ?", nome).First(&tipo).Error; err != nil {
		tipo = domain.CategoryType{BaseEntity: domain.BaseEntity{Nome: nome, Ativo: true}}
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
		if err := db.Where("nome = ?", name).First(&existing).Error; err != nil {
			app := domain.Application{BaseEntity: domain.BaseEntity{Nome: name, Ativo: true}}
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
				if existing.ID == 0 {
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

func updateCadastroCategory(db *gorm.DB, categorias []domain.Category, baseModel domain.BaseModel) error {
	// Associa múltiplas categorias ao cadastro via GORM M2M
	cadastro := domain.HistoricoPerfilApplication{BaseModel: baseModel}
	if err := db.Model(&cadastro).Association("Categorys").Append(&categorias); err != nil {
		return err
	}
	return nil
}

func associarCategorysAoCadastro(db *gorm.DB, categorias []domain.Category, nomeApp string, baseModel domain.BaseModel) error {
	// Busca nome do app no cadastro
	var cadastro domain.HistoricoPerfilApplication
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
	// Acumula categorias a associar
	var catsToAssociate []domain.Category
	// Associa regiões
	if regioes, ok := appInfo["regioes"].([]interface{}); ok {
		for _, reg := range regioes {
			regStr, ok := reg.(string)
			if !ok {
				continue
			}
			for _, cat := range categorias {
				if cat.Nome == regStr {
					catsToAssociate = append(catsToAssociate, cat)
				}
			}
		}
	}
	// Associa ramo
	if ramo, ok := appInfo["ramo"].(string); ok {
		for _, cat := range categorias {
			if cat.Nome == ramo {
				catsToAssociate = append(catsToAssociate, cat)
				break
			}
		}
	}
	// Associa subramo
	if subramo, ok := appInfo["subramo"].(string); ok {
		for _, cat := range categorias {
			if cat.Nome == subramo {
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
	cadastro := domain.HistoricoPerfilApplication{BaseModel: baseModel}
	if err := db.Model(&cadastro).Association("ConfiguracoesApplication").Append(&configuracoes); err != nil {
		return err
	}
	return nil
}

func SolicitarCadastroApplication(db *gorm.DB) error {
	var tipoRegiao, tipoRamos, tipoSubRamos domain.CategoryType
	db.Where("nome = ?", "Região").First(&tipoRegiao)
	db.Where("nome = ?", "Ramos").First(&tipoRamos)
	db.Where("nome = ?", "Subramos").First(&tipoSubRamos)
	var todasCategorys []domain.Category
	db.Where("category_type_id IN ?", []int64{tipoRegiao.ID, tipoRamos.ID, tipoSubRamos.ID}).Find(&todasCategorys)

	var apps []domain.Application
	db.Find(&apps)

	for _, app := range apps {
		// Buscar todas as configurações do app
		var configuracoes []domain.ApplicationConfiguration
		db.Preload("TerminalModel", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nome")
		}).Preload("IntegrationType", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nome")
		}).Where("application_id = ?", app.ID).Find(&configuracoes)

		cadastroConfig := domain.HistoricoPerfilApplication{
			ApplicationId: app.ID,
			ApplicationContact: domain.ApplicationContact{
				BaseEntity: domain.BaseEntity{
					Nome: app.Nome + " S.A.",
				},
				Email:    "ApplicationContact@" + app.Nome + ".com",
				Telefone: "1234-5678",
				Site:     "www." + app.Nome + ".com",
			},
			ApplicationDetail: domain.ApplicationDetail{
				Descricao: fmt.Sprintf(
					"Cadastro para %s",
					app.Nome),
			},
		}
		if tx := db.Create(&cadastroConfig); tx.Error != nil {
			return tx.Error
		}
		if err := associarCategorysAoCadastro(db, todasCategorys, app.Nome, cadastroConfig.BaseModel); err != nil {
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

func SolicitarCadastroVersaoApplication(db *gorm.DB) error {
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

		var versoes []domain.VersaoApplication
		for _, config := range cadastro.ConfiguracoesApplication {
			now := time.Now()
			year := now.Year() % 100
			dayOfYear := now.YearDay()
			seconds := now.Hour()*3600 + now.Minute()*60 + now.Second()
			nomeVersao := fmt.Sprintf("%02d%03d%05d", year, dayOfYear, seconds)
			tamanho := int64(rand.Intn(190)+10) * 1024 * 1024 // 10MB a 200MB

			versao := domain.VersaoApplication{
				ApplicationConfigurationId: config.ID,
				BaseEntity:                 domain.BaseEntity{Nome: app.Nome},
				Tamanho:                    tamanho,
				NomeVersao:                 nomeVersao,
				Imagem: domain.Imagem{
					Anexo: domain.Anexo{
						Nome:          fmt.Sprintf("Imagem da versão %s do aplicativo %s", nomeVersao, app.Nome),
						TipoMime:      "image/png",
						SHA256:        "d41d8cd98f00b204e9800998ecf8427e",
						Tamanho:       2048,
						Armazenamento: "S3",
						Caminho:       fmt.Sprintf("apps/%d/%s/icon.png", app.ID, nomeVersao),
						Presente:      true,
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

func getCadastroApplication(db *gorm.DB, appID int64) (domain.HistoricoPerfilApplication, error) {
	var cadastro domain.HistoricoPerfilApplication
	if err := db.Where("application_id = ?", appID).
		Order("created_at desc").
		Preload("ConfiguracoesApplication").
		First(&cadastro).Error; err != nil {
		return domain.HistoricoPerfilApplication{}, err
	}
	return cadastro, nil
}

func SolicitarPublicacaoVersaoApplication(db *gorm.DB) error {
	var apps []domain.Application
	if err := db.Find(&apps).Error; err != nil {
		return err
	}
	var estagios []domain.Estagio
	if err := db.Order("id ASC").Find(&estagios).Error; err != nil {
		return err
	}

	for _, app := range apps {
		cadastro, err := getCadastroApplication(db, app.ID)
		if err != nil || cadastro.ID == 0 {
			continue
		}

		// Para cada configuração do cadastro, buscar a versão mais recente e atualizar o catálogo
		for _, config := range cadastro.ConfiguracoesApplication {
			var versaoAtual domain.VersaoApplication
			if err := db.Where("id_cfg_aplv = ?", config.ID).Order("created_at desc").First(&versaoAtual).Error; err != nil || versaoAtual.ID == 0 {
				continue
			}

			var catalogos []domain.CatalogoApplication

			db.Preload("Estagio").Where("id_cfg_aplv = ?", config.ID).Order("id_est DESC").Find(&catalogos)

			lenCatalogos := len(catalogos)
			lenEstagios := len(estagios)

			// Cria novo catálogo se ainda não percorreu todos os estágios
			if lenCatalogos < lenEstagios {
				novoCatalogo := domain.CatalogoApplication{
					IdHistoricoPerfilApplication:    cadastro.ID,
					IdConfiguracaoPerfilApplication: config.ID,
					IdEstagio:                       estagios[lenCatalogos].ID,
					IdVersaoApplication: func() int64 {
						if lenCatalogos == 0 {
							return versaoAtual.ID
						}
						return catalogos[0].IdVersaoApplication
					}(),
				}
				if err := db.Create(&novoCatalogo).Error; err != nil {
					return err
				}
				catalogos = append([]domain.CatalogoApplication{novoCatalogo}, catalogos...)
				lenCatalogos++
			}

			// Avança os estágios das versões
			for i := 0; lenCatalogos > 0 && i < lenCatalogos; i++ {
				catalogoAtual := catalogos[i]
				if i < lenCatalogos-1 {
					catalogoAnterior := catalogos[i+1]
					catalogoAtual.IdVersaoApplication = catalogoAnterior.IdVersaoApplication
					if catalogoAtual.Estagio.Nome == "Produção" {
						catalogoAtual.IdHistoricoPerfilApplication = cadastro.ID
					}
				} else {
					catalogoAtual.IdVersaoApplication = versaoAtual.ID
				}
				if err := db.Save(&catalogoAtual).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func gerarJsonAppsPorRegiao(db *gorm.DB) error {
	consultaVersaoApp := func(db *gorm.DB, regiaoID, modeloID, integracaoID int64) ([]map[string]interface{}, error) {
		var result []map[string]interface{}
		err := db.Table("vrs_aplv").
			Select(`vrs_aplv.nome as nome_app, vrs_aplv.tamanho, vrs_aplv.id as id_versao, mdl_trml.nome as modelo_terminal_nome, integration_type.nome as tipo_integracao_nome, ctgr.nome as categoria_nome`).
			Joins("JOIN ctlg_aplv ON ctlg_aplv.id_vrs_aplv = vrs_aplv.id").                                                                                     //Localiza a versao no catalogo
			Joins("JOIN cfg_aplv ON cfg_aplv.id = vrs_aplv.id_cfg_aplv").                                                                                       //Localiza a configuracao
			Joins("JOIN application_profile_history ON application_profile_history.id = ctlg_aplv.id_application_profile_history").                             //Localiza o cadastro
			Joins("JOIN ctgr_application_profile_history ON ctgr_application_profile_history.id_application_profile_history = application_profile_history.id"). //M2M cadastro-categoria
			Joins("JOIN ctgr ON ctgr.id = ctgr_application_profile_history.id_ctgr").                                                                           //Localiza a categoria
			Joins("JOIN mdl_trml ON mdl_trml.id = cfg_aplv.terminal_model_id").                                                                                 //Localiza o modelo do terminal
			Joins("JOIN integration_type ON integration_type.id = cfg_aplv.integration_type_id").                                                               //Localiza o tipo de integração
			Where("cfg_aplv.terminal_model_id = ?", modeloID).
			Where("cfg_aplv.integration_type_id = ?", integracaoID).
			Where("ctgr.id = ?", regiaoID).
			Order("vrs_aplv.created_at DESC").
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
		err := db.Table("ctlg_aplv").
			Select(`vrs_aplv.nome as nome_app, vrs_aplv.tamanho, vrs_aplv.id as id_versao, mdl_trml.nome as modelo_terminal_nome, integration_type.nome as tipo_integracao_nome, ctgr.nome as categoria_nome`).
			Joins("JOIN vrs_aplv ON vrs_aplv.id = ctlg_aplv.id_vrs_aplv").                                                                                      //Localiza a versão no catálogo
			Joins("JOIN cfg_aplv ON cfg_aplv.id = ctlg_aplv.id_cfg_aplv").                                                                                      //Localiza a configuracao
			Joins("JOIN application_profile_history ON application_profile_history.id = ctlg_aplv.id_application_profile_history").                             //Localiza o cadastro
			Joins("JOIN ctgr_application_profile_history ON ctgr_application_profile_history.id_application_profile_history = application_profile_history.id"). //M2M cadastro-categoria
			Joins("JOIN ctgr ON ctgr.id = ctgr_application_profile_history.id_ctgr").                                                                           //Localiza a categoria
			Joins("JOIN mdl_trml ON mdl_trml.id = cfg_aplv.terminal_model_id").                                                                                 //Localiza o modelo do terminal
			Joins("JOIN integration_type ON integration_type.id = cfg_aplv.integration_type_id").                                                               //Localiza o tipo de integração
			Where("cfg_aplv.terminal_model_id = ?", modeloID).
			Where("cfg_aplv.integration_type_id = ?", integracaoID).
			Where("ctgr.id = ?", regiaoID).
			Order("vrs_aplv.created_at DESC").
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
		Joins("JOIN category_type ON category_type.id = ctgr.category_type_id").
		Where("category_type.nome = ?", "Região").
		Pluck("ctgr.id", &idsRegioes)

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
				fileName := fmt.Sprintf(fileNameFmt, regiao.Nome, modelo.Nome, integracao.Nome)
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
				//fmt.Printf("%s %s, Modelo %s, Integracao %s: %d itens, consulta e exportacao demorou %d ms\n", logPrefix, regiao.Nome, modelo.Nome, integracao.Nome, len(apps), duracao)
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
	setup := false
	if setup {
		CadastroCategorysFromJSONPath(db, "regioes.json")
		CadastroCategorysFromJSONPath(db, "ramos_subramos.json")
		CadastroModelosTerminalFromJSON(db)
		CadastroTiposIntegracaoFromJSON(db)
		CadastroApplicationsFromJSON(db)
		CadastroConfiguraces(db)
		CadastroEstagiosFromJSON(db)
	}

	loadVersao := false

	if loadVersao {
		for i := 0; i < 50; i++ {
			SolicitarCadastroApplication(db)
			SolicitarCadastroVersaoApplication(db)
			SolicitarPublicacaoVersaoApplication(db)
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
