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
		modelo := domain.ModeloTerminal{BaseEntity: be}
		var existing domain.ModeloTerminal
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
		tipo := domain.TipoIntegracao{BaseEntity: be}
		var existing domain.TipoIntegracao
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

func CadastroCategoriasFromJSONPath(db *gorm.DB, relativePath string) error {
	baseDir := "c:/Projects/store-go/shared/testdata"
	absPath := filepath.Join(baseDir, relativePath)
	file, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var root CategoriasRoot
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&root); err != nil {
		return err
	}
	return cadastroCategoriasFromJSONRec(db, root, nil)
}
func cadastroCategoriasFromJSONRec(db *gorm.DB, root CategoriasRoot, idPai *int64) error {
	// Cadastra o tipo de categoria
	tipoCat, err := ensureTipoCategoria(db, root.TipoCategoria)
	if err != nil {
		return err
	}
	for nome, raw := range root.Categorias {
		// Cadastra a categoria
		categoria := domain.Categoria{
			BaseEntity:      domain.BaseEntity{Nome: nome, Ativo: true},
			IdTipoCategoria: tipoCat.ID,
			IdPai:           idPai,
		}
		var existing domain.Categoria
		err := db.Where("nome = ? AND id_tip_ctgr = ?", nome, tipoCat.ID).First(&existing).Error
		if err != nil {
			if err := db.Create(&categoria).Error; err != nil {
				return err
			}
		}
		// Verifica se há subcategorias
		var subcat CategoriasRoot
		if err := json.Unmarshal(raw, &subcat); err == nil && len(subcat.Categorias) > 0 {
			// Busca o ID da categoria recém-cadastrada
			if err := cadastroCategoriasFromJSONRec(db, subcat, &categoria.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

type CategoriasRoot struct {
	TipoCategoria string                     `json:"tipo_categoria"`
	Categorias    map[string]json.RawMessage `json:"categorias"`
}

func ensureTipoCategoria(db *gorm.DB, nome string) (domain.TipoCategoria, error) {
	var tipo domain.TipoCategoria
	if err := db.Where("nome = ?", nome).First(&tipo).Error; err != nil {
		tipo = domain.TipoCategoria{BaseEntity: domain.BaseEntity{Nome: nome, Ativo: true}}
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

func CadastroAplicativosFromJSON(db *gorm.DB) {
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
		var existing domain.Aplicativo
		if err := db.Where("nome = ?", name).First(&existing).Error; err != nil {
			app := domain.Aplicativo{BaseEntity: domain.BaseEntity{Nome: name, Ativo: true}}
			db.Create(&app)
		}
	}
}

func CadastroConfiguraces(db *gorm.DB) {
	// Buscar todos os modelos de terminal
	var modelos []domain.ModeloTerminal
	db.Find(&modelos)

	// Buscar todos os tipos de integração
	var integracoes []domain.TipoIntegracao
	db.Find(&integracoes)

	// Buscar todos os aplicativos
	var apps []domain.Aplicativo
	db.Find(&apps)

	for _, app := range apps {
		for _, modelo := range modelos {
			for _, integracao := range integracoes {
				// Verifica se já existe configuração
				var existing domain.ConfiguracaoAplicativo
				db.Where("id_aplv = ? AND id_mdl_trml = ? AND id_tip_itgr = ?", app.ID, modelo.ID, integracao.ID).First(&existing)
				if existing.ID == 0 {
					config := domain.ConfiguracaoAplicativo{
						IdAplicativo:     app.ID,
						IdModeloTerminal: modelo.ID,
						IdTipoIntegracao: integracao.ID,
					}
					db.Create(&config)
				}
			}
		}
	}
}

func updateCadastroCategoria(db *gorm.DB, categorias []domain.Categoria, baseModel domain.BaseModel) error {
	// Associa múltiplas categorias ao cadastro via GORM M2M
	cadastro := domain.HistoricoPerfilAplicativo{BaseModel: baseModel}
	if err := db.Model(&cadastro).Association("Categorias").Append(&categorias); err != nil {
		return err
	}
	return nil
}

func associarCategoriasAoCadastro(db *gorm.DB, categorias []domain.Categoria, nomeApp string, baseModel domain.BaseModel) error {
	// Busca nome do app no cadastro
	var cadastro domain.HistoricoPerfilAplicativo
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
	var catsToAssociate []domain.Categoria
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
		if err := updateCadastroCategoria(db, catsToAssociate, baseModel); err != nil {
			return err
		}
	}
	return nil
}

func associarConfiguracaoAoCadastro(db *gorm.DB, configuracoes []domain.ConfiguracaoAplicativo, baseModel domain.BaseModel) error {
	// Associa múltiplas configurações ao cadastro via GORM M2M
	cadastro := domain.HistoricoPerfilAplicativo{BaseModel: baseModel}
	if err := db.Model(&cadastro).Association("ConfiguracoesAplicativo").Append(&configuracoes); err != nil {
		return err
	}
	return nil
}

func SolicitarCadastroAplicativo(db *gorm.DB) error {
	var tipoRegiao, tipoRamos, tipoSubRamos domain.TipoCategoria
	db.Where("nome = ?", "Região").First(&tipoRegiao)
	db.Where("nome = ?", "Ramos").First(&tipoRamos)
	db.Where("nome = ?", "Subramos").First(&tipoSubRamos)
	var todasCategorias []domain.Categoria
	db.Where("id_tip_ctgr IN ?", []int64{tipoRegiao.ID, tipoRamos.ID, tipoSubRamos.ID}).Find(&todasCategorias)

	var apps []domain.Aplicativo
	db.Find(&apps)

	for _, app := range apps {
		// Buscar todas as configurações do app
		var configuracoes []domain.ConfiguracaoAplicativo
		db.Preload("ModeloTerminal", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nome")
		}).Preload("TipoIntegracao", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nome")
		}).Where("id_aplv = ?", app.ID).Find(&configuracoes)

		cadastroConfig := domain.HistoricoPerfilAplicativo{
			IdAplicativo: app.ID,
			ContatoAplicativo: domain.ContatoAplicativo{
				BaseEntity: domain.BaseEntity{
					Nome: app.Nome + " S.A.",
				},
				Email:    "ContatoAplicativo@" + app.Nome + ".com",
				Telefone: "1234-5678",
				Site:     "www." + app.Nome + ".com",
			},
			DetalheAplicativo: domain.DetalheAplicativo{
				Descricao: fmt.Sprintf(
					"Cadastro para %s",
					app.Nome),
			},
		}
		if tx := db.Create(&cadastroConfig); tx.Error != nil {
			return tx.Error
		}
		if err := associarCategoriasAoCadastro(db, todasCategorias, app.Nome, cadastroConfig.BaseModel); err != nil {
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

func SolicitarCadastroVersaoAplicativo(db *gorm.DB) error {
	var apps []domain.Aplicativo
	if err := db.Find(&apps).Error; err != nil {
		return err
	}

	for _, app := range apps {
		// Buscar o cadastro mais recente para o app
		cadastro, err := getCadastroAplicativo(db, app.ID)

		if err != nil || cadastro.ID == 0 {
			continue
		}

		var versoes []domain.VersaoAplicativo
		for _, config := range cadastro.ConfiguracoesAplicativo {
			now := time.Now()
			year := now.Year() % 100
			dayOfYear := now.YearDay()
			seconds := now.Hour()*3600 + now.Minute()*60 + now.Second()
			nomeVersao := fmt.Sprintf("%02d%03d%05d", year, dayOfYear, seconds)
			tamanho := int64(rand.Intn(190)+10) * 1024 * 1024 // 10MB a 200MB

			versao := domain.VersaoAplicativo{
				IdConfiguracaoAplicativo: config.ID,
				BaseEntity:               domain.BaseEntity{Nome: app.Nome},
				Tamanho:                  tamanho,
				NomeVersao:               nomeVersao,
				Imagem: domain.Imagem{
					Anexo: domain.Anexo{
						Nome:          fmt.Sprintf("Imagem da versão %s do aplicativo %s", nomeVersao, app.Nome),
						TipoMime:      "image/png",
						MD5:           "d41d8cd98f00b204e9800998ecf8427e",
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

func getCadastroAplicativo(db *gorm.DB, appID int64) (domain.HistoricoPerfilAplicativo, error) {
	var cadastro domain.HistoricoPerfilAplicativo
	if err := db.Where("id_aplv = ?", appID).
		Order("created_at desc").
		Preload("ConfiguracoesAplicativo").
		First(&cadastro).Error; err != nil {
		return domain.HistoricoPerfilAplicativo{}, err
	}
	return cadastro, nil
}

func SolicitarPublicacaoVersaoAplicativo(db *gorm.DB) error {
	var apps []domain.Aplicativo
	if err := db.Find(&apps).Error; err != nil {
		return err
	}
	var estagios []domain.Estagio
	if err := db.Order("id ASC").Find(&estagios).Error; err != nil {
		return err
	}

	for _, app := range apps {
		cadastro, err := getCadastroAplicativo(db, app.ID)
		if err != nil || cadastro.ID == 0 {
			continue
		}

		// Para cada configuração do cadastro, buscar a versão mais recente e atualizar o catálogo
		for _, config := range cadastro.ConfiguracoesAplicativo {
			var versaoAtual domain.VersaoAplicativo
			if err := db.Where("id_cfg_aplv = ?", config.ID).Order("created_at desc").First(&versaoAtual).Error; err != nil || versaoAtual.ID == 0 {
				continue
			}

			var catalogos []domain.CatalogoAplicativo

			db.Preload("Estagio").Where("id_cfg_aplv = ?", config.ID).Order("id_est DESC").Find(&catalogos)

			lenCatalogos := len(catalogos)
			lenEstagios := len(estagios)

			// Cria novo catálogo se ainda não percorreu todos os estágios
			if lenCatalogos < lenEstagios {
				novoCatalogo := domain.CatalogoAplicativo{
					IdHistoricoPerfilAplicativo:    cadastro.ID,
					IdConfiguracaoPerfilAplicativo: config.ID,
					IdEstagio:                      estagios[lenCatalogos].ID,
					IdVersaoAplicativo: func() int64 {
						if lenCatalogos == 0 {
							return versaoAtual.ID
						}
						return catalogos[0].IdVersaoAplicativo
					}(),
				}
				if err := db.Create(&novoCatalogo).Error; err != nil {
					return err
				}
				catalogos = append([]domain.CatalogoAplicativo{novoCatalogo}, catalogos...)
				lenCatalogos++
			}

			// Avança os estágios das versões
			for i := 0; lenCatalogos > 0 && i < lenCatalogos; i++ {
				catalogoAtual := catalogos[i]
				if i < lenCatalogos-1 {
					catalogoAnterior := catalogos[i+1]
					catalogoAtual.IdVersaoAplicativo = catalogoAnterior.IdVersaoAplicativo
				} else {
					catalogoAtual.IdVersaoAplicativo = versaoAtual.ID
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
		err := db.Debug().Table("vrs_aplv").
			Select(`vrs_aplv.nome as nome_app, vrs_aplv.tamanho, mdl_trml.nome as modelo_terminal_nome, tip_itgr.nome as tipo_integracao_nome, ctgr.nome as categoria_nome, vrs_aplv.id as id_versao`).
			Joins("JOIN ctlg_aplv ON ctlg_aplv.id_vrs_aplv = vrs_aplv.id").                             //Localiza a versao no catalogo
			Joins("JOIN cfg_aplv ON cfg_aplv.id = vrs_aplv.id_cfg_aplv").                               //Localiza a configuracao
			Joins("JOIN hist_pfl_aplv ON hist_pfl_aplv.id = ctlg_aplv.id_hist_pfl_aplv").               //Localiza o cadastro
			Joins("JOIN ctgr_hist_pfl_aplv ON ctgr_hist_pfl_aplv.id_hist_pfl_aplv = hist_pfl_aplv.id"). //M2M cadastro-categoria
			Joins("JOIN ctgr ON ctgr.id = ctgr_hist_pfl_aplv.id_ctgr").                                 //Localiza a categoria
			Joins("JOIN mdl_trml ON mdl_trml.id = cfg_aplv.id_mdl_trml").                               //Localiza o modelo do terminal
			Joins("JOIN tip_itgr ON tip_itgr.id = cfg_aplv.id_tip_itgr").                               //Localiza o tipo de integração
			Where("cfg_aplv.id_mdl_trml = ?", modeloID).
			Where("cfg_aplv.id_tip_itgr = ?", integracaoID).
			Where("ctgr.id = ?", regiaoID).
			Order("vrs_aplv.created_at DESC").
			Scan(&result).Error
		return result, err
	}
	metricas := GetMetricasGlobal()
	return gerarJsonAppsGenerico(db, consultaVersaoApp, "apps_regiao_%s_%s_%s.json", "gerarJsonAppsPorRegiao", metricas)
}

// Exporta apps por região usando CatalogoAplicativo
func gerarJsonAppsPorRegiaoCatalogo(db *gorm.DB) error {
	consultaCatalogoApp := func(db *gorm.DB, regiaoID, modeloID, integracaoID int64) ([]map[string]interface{}, error) {
		var result []map[string]interface{}
		err := db.Table("ctlg_aplv").
			Select(`vrs_aplv.nome as nome_app, vrs_aplv.tamanho, mdl_trml.nome as modelo_terminal_nome, tip_itgr.nome as tipo_integracao_nome, ctgr.nome as categoria_nome, vrs_aplv.id as id_versao`).
			Joins("JOIN vrs_aplv ON vrs_aplv.id = ctlg_aplv.id_vrs_aplv").                              //Localiza a versão no catálogo
			Joins("JOIN cfg_aplv ON cfg_aplv.id = vrs_aplv.id_cfg_aplv").                               //Localiza a configuracao
			Joins("JOIN hist_pfl_aplv ON hist_pfl_aplv.id = ctlg_aplv.id_hist_pfl_aplv").               //Localiza o cadastro
			Joins("JOIN ctgr_hist_pfl_aplv ON ctgr_hist_pfl_aplv.id_hist_pfl_aplv = hist_pfl_aplv.id"). //M2M cadastro-categoria
			Joins("JOIN ctgr ON ctgr.id = ctgr_hist_pfl_aplv.id_ctgr").                                 //Localiza a categoria
			Joins("JOIN mdl_trml ON mdl_trml.id = cfg_aplv.id_mdl_trml").                               //Localiza o modelo do terminal
			Joins("JOIN tip_itgr ON tip_itgr.id = cfg_aplv.id_tip_itgr").                               //Localiza o tipo de integração
			Where("cfg_aplv.id_mdl_trml = ?", modeloID).
			Where("cfg_aplv.id_tip_itgr = ?", integracaoID).
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
	db.Model(&domain.Categoria{}).
		Joins("JOIN tip_ctgr ON tip_ctgr.id = ctgr.id_tip_ctgr").
		Where("tip_ctgr.nome = ?", "Região").
		Pluck("ctgr.id", &idsRegioes)

	var modelos []domain.ModeloTerminal
	db.Find(&modelos)
	var integracoes []domain.TipoIntegracao
	db.Find(&integracoes)

	var totalItems int
	var totalQueries int
	times := make([]int64, 0)
	for _, regiaoID := range idsRegioes {
		var regiao domain.Categoria
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
		CadastroCategoriasFromJSONPath(db, "regioes.json")
		CadastroCategoriasFromJSONPath(db, "ramos_subramos.json")
		CadastroModelosTerminalFromJSON(db)
		CadastroTiposIntegracaoFromJSON(db)
		CadastroAplicativosFromJSON(db)
		CadastroConfiguraces(db)
		CadastroEstagiosFromJSON(db)
	}

	loadVersao := false

	if loadVersao {
		for i := 0; i < 50; i++ {
			SolicitarCadastroAplicativo(db)
			SolicitarCadastroVersaoAplicativo(db)
			SolicitarPublicacaoVersaoAplicativo(db)
		}
	}

	//gerarJsonAppsPorRegiao(db)
	gerarJsonAppsPorRegiaoCatalogo(db)

	// var wg sync.WaitGroup
	// metricasGlobal = &MetricasExecucao{} // resetar métricas
	// for i := 0; i < 10; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		for j := 0; j < 20; j++ {
	// 			gerarJsonAppsPorRegiaoCatalogo(db)
	// 			time.Sleep(100 * time.Millisecond)
	// 		}
	// 	}()
	// }
	// wg.Wait()
	// metricasGlobal.CalcularMetricas("gerarJsonAppsPorRegiao")
}
