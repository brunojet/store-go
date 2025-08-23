// Para cada configuração de cada app, resgata as 4 versões mais recentes e insere no catálogo com os estágios
package pkg

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"store-go/shared/internal"
	"store-go/shared/internal/domain"
	"sync"
	"testing"
	"time"

	"gorm.io/gorm"
)

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

func tipoCategoria(db *gorm.DB) {

	//Preencher TipoCategoria
	tipos := []domain.TipoCategoria{
		{BaseEntity: domain.BaseEntity{Nome: "Região", Ativo: true}},
		{BaseEntity: domain.BaseEntity{Nome: "Ramos", Ativo: true}},
		{BaseEntity: domain.BaseEntity{Nome: "SubRamos", Ativo: true}},
	}
	for _, tipo := range tipos {
		var existing domain.TipoCategoria
		if err := db.Where("nome = ?", tipo.Nome).First(&existing).Error; err != nil {
			db.Create(&tipo)
		}
	}
}

func categoriasRegiao(db *gorm.DB) {
	var tipoRegiao domain.TipoCategoria

	db.Where("nome = ?", "Região").First(&tipoRegiao)

	categorias := []domain.Categoria{
		{BaseEntity: domain.BaseEntity{Nome: "Norte", Ativo: true}, IdTipoCategoria: tipoRegiao.ID},
		{BaseEntity: domain.BaseEntity{Nome: "Nordeste", Ativo: true}, IdTipoCategoria: tipoRegiao.ID},
		{BaseEntity: domain.BaseEntity{Nome: "Centro-Oeste", Ativo: true}, IdTipoCategoria: tipoRegiao.ID},
		{BaseEntity: domain.BaseEntity{Nome: "Sudeste", Ativo: true}, IdTipoCategoria: tipoRegiao.ID},
		{BaseEntity: domain.BaseEntity{Nome: "Sul", Ativo: true}, IdTipoCategoria: tipoRegiao.ID},
	}
	for _, categoria := range categorias {
		var existing domain.Categoria
		if err := db.Where("nome = ? AND id_tipo_categoria = ?", categoria.Nome, categoria.IdTipoCategoria).First(&existing).Error; err != nil {
			db.Create(&categoria)
		}
	}
}

func categoriasRamo(db *gorm.DB) {
	var tipoRamos domain.TipoCategoria
	// Busca o tipo de categoria "Ramos"
	db.Where("nome = ?", "Ramos").First(&tipoRamos)

	ramosMacro := []string{
		"Varejo", "Atacado", "Veículos", "Alimentação", "Saúde",
		"Tecnologia", "Entretenimento", "Educação", "Finanças", "Imóveis",
		"Transporte", "Comunicação", "Utilidades", "Compras", "Serviços",
	}
	for _, nome := range ramosMacro {
		categoria := domain.Categoria{
			BaseEntity:      domain.BaseEntity{Nome: nome, Ativo: true},
			IdTipoCategoria: tipoRamos.ID,
		}
		var existing domain.Categoria
		if err := db.Where("nome = ? AND id_tipo_categoria = ?", nome, tipoRamos.ID).First(&existing).Error; err != nil {
			db.Create(&categoria)
		}
	}
}

func categoriasSubRamo(db *gorm.DB) {
	var tipoSubRamos domain.TipoCategoria
	db.Where("nome = ?", "SubRamos").First(&tipoSubRamos)

	// Buscar todos os ramos macro
	var ramosMacro []domain.Categoria
	db.Where("id_tipo_categoria = ?", tipoSubRamos.ID-1).Find(&ramosMacro) // tipoRamos.ID = tipoSubRamos.ID-1

	// Subramos para cada ramo macro
	subramosPorRamo := map[string][]string{
		"Varejo":         {"Supermercado", "Loja de Roupas", "Papelaria", "Perfumaria", "Livraria", "Magazine", "E-commerce"},
		"Atacado":        {"Distribuidora", "Atacado de Alimentos", "Atacado de Bebidas", "Atacado de Limpeza", "Atacado de Eletrônicos"},
		"Veículos":       {"Concessionária", "Autopeças", "Oficina", "Locadora", "Revenda"},
		"Alimentação":    {"Restaurante", "Padaria", "Lanchonete", "Pizzaria", "Cafeteria", "Delivery", "Fast Food"},
		"Saúde":          {"Farmácia", "Clínica", "Laboratório", "Ótica", "Drogaria"},
		"Tecnologia":     {"Software", "Cloud", "Streaming", "Redes Sociais", "Mensageria", "Banco Digital", "Marketplace", "Aplicativo Móvel"},
		"Entretenimento": {"Cinema", "Música", "TV", "Streaming", "Jogos", "Eventos", "Shows"},
		"Educação":       {"Cursos", "Plataforma EAD", "Idiomas", "Faculdade", "Aulas Online"},
		"Finanças":       {"Banco Digital", "Carteira Digital", "Pagamentos", "Investimentos", "Seguros"},
		"Imóveis":        {"Aluguel", "Compra", "Venda", "Condomínio", "Imobiliária"},
		"Transporte":     {"Táxi", "Carona", "Mobilidade", "Ônibus", "Avião", "Locadora"},
		"Comunicação":    {"Mensageria", "Redes Sociais", "E-mail", "Chamadas", "Videoconferência"},
		"Utilidades":     {"Armazenamento", "Backup", "Organização", "Notas", "Calendário"},
		"Compras":        {"Marketplace", "E-commerce", "Lojas", "Ofertas", "Comparador de Preços"},
		"Serviços":       {"Delivery", "Assistência", "Consultoria", "Manutenção", "Limpeza"},
	}

	for _, ramo := range ramosMacro {
		subramos := subramosPorRamo[ramo.Nome]
		for _, nome := range subramos {
			categoria := domain.Categoria{
				BaseEntity:      domain.BaseEntity{Nome: nome, Ativo: true},
				IdTipoCategoria: tipoSubRamos.ID,
				IdPai:           &ramo.ID,
			}
			var existing domain.Categoria
			if err := db.Where("nome = ? AND id_tipo_categoria = ? AND id_pai = ?", nome, tipoSubRamos.ID, ramo.ID).First(&existing).Error; err != nil {
				db.Create(&categoria)
			}
		}
	}
}

func modelosTerminal(db *gorm.DB) {
	modelos := []domain.ModeloTerminal{
		{BaseEntity: domain.BaseEntity{Nome: "L400", Ativo: true}},
		{BaseEntity: domain.BaseEntity{Nome: "S350", Ativo: true}},
		{BaseEntity: domain.BaseEntity{Nome: "GPOS700", Ativo: true}},
	}
	for _, modelo := range modelos {
		var existing domain.ModeloTerminal
		if err := db.Where("nome = ?", modelo.Nome).First(&existing).Error; err != nil {
			db.Create(&modelo)
		}
	}
}

func tiposIntegracao(db *gorm.DB) {
	integracoes := []domain.TipoIntegracao{
		{BaseEntity: domain.BaseEntity{Nome: "ADQ", Ativo: true}},
		{BaseEntity: domain.BaseEntity{Nome: "TEF", Ativo: true}},
	}
	for _, integracao := range integracoes {
		var existing domain.TipoIntegracao
		if err := db.Where("nome = ?", integracao.Nome).First(&existing).Error; err != nil {
			db.Create(&integracao)
		}
	}
}

func aplicativos(db *gorm.DB) {
	appNames := []string{
		"WhatsApp", "Instagram", "Facebook", "Twitter", "Telegram", "Spotify", "Netflix", "YouTube", "Gmail", "Google Maps",
		"Uber", "Airbnb", "Dropbox", "Slack", "Zoom", "Teams", "Outlook", "Pinterest", "LinkedIn", "Reddit",
		"Snapchat", "TikTok", "Shazam", "Duolingo", "Evernote", "Tinder", "Discord", "Notion", "Canva", "PayPal",
		"Mercado Livre", "iFood", "Rappi", "99", "PicPay", "Banco Inter", "Nubank", "Santander", "Bradesco", "Banco do Brasil",
		"Magazine Luiza", "Amazon", "Submarino", "Americanas", "Casas Bahia", "Carrefour", "Extra", "Pão de Açúcar", "Netshoes", "Centauro",
		"OLX", "VivaReal", "Zap Imóveis", "QuintoAndar", "Catho", "InfoJobs", "Buscapé", "Peixe Urbano", "Booking", "Decolar",
		"Skyscanner", "TripAdvisor", "Moovit", "Waze", "CVC", "Hurb", "Shell Box", "Petrobras Premmia", "Google Drive", "OneDrive",
		"Adobe Reader", "Photoshop Express", "Lightroom", "VSCO", "PicsArt", "Camera360", "Remini", "InShot", "CapCut", "Kwai",
		"GloboPlay", "Prime Video", "Disney+", "HBO Max", "Paramount+", "Star+", "Crunchyroll", "Twitch", "Steam", "Epic Games",
		"Duolingo", "Coursera", "Udemy", "Khan Academy", "Alura", "Rocketseat", "Senai", "Sebrae", "Enem", "Estuda.com",
	}
	for _, name := range appNames {
		var existing domain.Aplicativo
		if err := db.Where("nome = ?", name).First(&existing).Error; err != nil {
			app := domain.Aplicativo{BaseEntity: domain.BaseEntity{Nome: name, Ativo: true}}
			db.Create(&app)
		}
	}
}

func configuraces(db *gorm.DB) {
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
				var existing domain.Configuracao
				db.Where("id_app = ? AND id_modelo_terminal = ? AND id_tipo_integracao = ?", app.ID, modelo.ID, integracao.ID).First(&existing)
				if existing.ID == 0 {
					config := domain.Configuracao{
						IdApp:            app.ID,
						IdModeloTerminal: modelo.ID,
						IdTipoIntegracao: integracao.ID,
					}
					db.Create(&config)
				}
			}
		}
	}
}

func cadastroPorConfiguracao(db *gorm.DB) {
	// 1. Buscar todos os aplicativos
	var apps []domain.Aplicativo
	db.Find(&apps)

	for _, app := range apps {
		// 2. Criar cadastro para o app
		cadastro := domain.Cadastro{}
		db.Create(&cadastro)

		// 3. Buscar todas as configurações ativas do app
		var configuracoes []domain.Configuracao
		db.Where("id_app = ?", app.ID).Find(&configuracoes)

		// 4. Relacionar o cadastro às configurações via ConfiguracaoCadastro
		for _, config := range configuracoes {
			configCadastro := domain.ConfiguracaoCadastro{
				IdCadastro:     cadastro.ID,
				IdConfiguracao: config.ID,
			}
			db.Create(&configCadastro)
		}
	}
}

func cadastrarVersaoApp(db *gorm.DB) {
	// 1. Buscar todos os aplicativos
	var apps []domain.Aplicativo
	db.Find(&apps)

	for _, app := range apps {
		// 2. Buscar todas as configurações do app
		var configuracoes []domain.Configuracao
		db.Where("id_app = ?", app.ID).Find(&configuracoes)

		for _, config := range configuracoes {
			// 3. Localizar o cadastro mais recente da configuração
			var cfgCad domain.ConfiguracaoCadastro
			db.Where("id_configuracao = ?", config.ID).Order("created_at desc").First(&cfgCad)
			if cfgCad.ID == 0 {
				continue // Nenhum cadastro para esta configuração
			}

			// 4. Cadastrar a versão
			// Gerar nome da versão: YYDDD + segundos do dia (5 dígitos)
			now := time.Now()
			year := now.Year() % 100
			dayOfYear := now.YearDay()
			seconds := now.Hour()*3600 + now.Minute()*60 + now.Second()
			nomeVersao := fmt.Sprintf("%02d%03d%05d", year, dayOfYear, seconds)

			versao := domain.VersaoAplicativo{
				IdCadastro:     cfgCad.IdCadastro,
				IdConfiguracao: config.ID,
				BaseEntity:     domain.BaseEntity{Nome: app.Nome}, // Corrected to match the structure
				Tamanho:        123456,                            // valor exemplo, pode ser ajustado
				NomeVersao:     nomeVersao,
			}
			result := db.Create(&versao)
			if result.Error != nil {
				fmt.Printf("Erro ao criar versao_app: %v\n", result.Error)
			}
		}
	}
}

func associarRegioesAosApps(db *gorm.DB) {
	// Buscar todas as regiões
	var tipoRegiao domain.TipoCategoria
	db.Where("nome = ?", "Região").First(&tipoRegiao)
	var regioes []domain.Categoria
	db.Where("id_tipo_categoria = ?", tipoRegiao.ID).Find(&regioes)

	// Buscar todos os aplicativos
	var apps []domain.Aplicativo
	db.Find(&apps)

	for _, app := range apps {
		// Selecionar 2 regiões aleatórias
		n := len(regioes)
		if n < 2 {
			continue // Não há regiões suficientes
		}
		idxs := rand.Perm(n)[:2]
		for _, idx := range idxs {
			appCat := domain.AppCategoria{
				IdApp:       app.ID,
				IdCategoria: regioes[idx].ID,
			}
			db.Create(&appCat)
		}
	}
}

func associarRamosESubramosAosApps(db *gorm.DB) {
	// Mapeamento simplificado: appName -> ramo, subramo
	appToCategoria := map[string][2]string{
		"WhatsApp":          {"Comunicação", "Mensageria"},
		"Instagram":         {"Comunicação", "Redes Sociais"},
		"Facebook":          {"Comunicação", "Redes Sociais"},
		"Twitter":           {"Comunicação", "Redes Sociais"},
		"Telegram":          {"Comunicação", "Mensageria"},
		"Spotify":           {"Entretenimento", "Música"},
		"Netflix":           {"Entretenimento", "Streaming"},
		"YouTube":           {"Entretenimento", "Streaming"},
		"Gmail":             {"Comunicação", "E-mail"},
		"Google Maps":       {"Transporte", "Mobilidade"},
		"Uber":              {"Transporte", "Mobilidade"},
		"Airbnb":            {"Imóveis", "Aluguel"},
		"Dropbox":           {"Utilidades", "Armazenamento"},
		"Slack":             {"Comunicação", "Mensageria"},
		"Zoom":              {"Comunicação", "Videoconferência"},
		"Teams":             {"Comunicação", "Videoconferência"},
		"Outlook":           {"Comunicação", "E-mail"},
		"Pinterest":         {"Entretenimento", "Redes Sociais"},
		"LinkedIn":          {"Comunicação", "Redes Sociais"},
		"Reddit":            {"Entretenimento", "Redes Sociais"},
		"Snapchat":          {"Comunicação", "Redes Sociais"},
		"TikTok":            {"Entretenimento", "Streaming"},
		"Shazam":            {"Entretenimento", "Música"},
		"Duolingo":          {"Educação", "Idiomas"},
		"Evernote":          {"Utilidades", "Notas"},
		"Tinder":            {"Entretenimento", "Eventos"},
		"Discord":           {"Comunicação", "Mensageria"},
		"Notion":            {"Utilidades", "Organização"},
		"Canva":             {"Tecnologia", "Software"},
		"PayPal":            {"Finanças", "Pagamentos"},
		"Mercado Livre":     {"Compras", "Marketplace"},
		"iFood":             {"Alimentação", "Delivery"},
		"Rappi":             {"Alimentação", "Delivery"},
		"99":                {"Transporte", "Mobilidade"},
		"PicPay":            {"Finanças", "Carteira Digital"},
		"Banco Inter":       {"Finanças", "Banco Digital"},
		"Nubank":            {"Finanças", "Banco Digital"},
		"Santander":         {"Finanças", "Banco Digital"},
		"Bradesco":          {"Finanças", "Banco Digital"},
		"Banco do Brasil":   {"Finanças", "Banco Digital"},
		"Magazine Luiza":    {"Varejo", "Magazine"},
		"Amazon":            {"Compras", "Marketplace"},
		"Submarino":         {"Compras", "Marketplace"},
		"Americanas":        {"Compras", "Marketplace"},
		"Casas Bahia":       {"Varejo", "Magazine"},
		"Carrefour":         {"Varejo", "Supermercado"},
		"Extra":             {"Varejo", "Supermercado"},
		"Pão de Açúcar":     {"Varejo", "Supermercado"},
		"Netshoes":          {"Varejo", "Lojas"},
		"Centauro":          {"Varejo", "Lojas"},
		"OLX":               {"Imóveis", "Venda"},
		"VivaReal":          {"Imóveis", "Aluguel"},
		"Zap Imóveis":       {"Imóveis", "Aluguel"},
		"QuintoAndar":       {"Imóveis", "Aluguel"},
		"Catho":             {"Serviços", "Consultoria"},
		"InfoJobs":          {"Serviços", "Consultoria"},
		"Buscapé":           {"Compras", "Comparador de Preços"},
		"Peixe Urbano":      {"Serviços", "Ofertas"},
		"Booking":           {"Imóveis", "Aluguel"},
		"Decolar":           {"Transporte", "Avião"},
		"Skyscanner":        {"Transporte", "Avião"},
		"TripAdvisor":       {"Serviços", "Consultoria"},
		"Moovit":            {"Transporte", "Mobilidade"},
		"Waze":              {"Transporte", "Mobilidade"},
		"CVC":               {"Serviços", "Consultoria"},
		"Hurb":              {"Serviços", "Consultoria"},
		"Shell Box":         {"Serviços", "Assistência"},
		"Petrobras Premmia": {"Serviços", "Assistência"},
		"Google Drive":      {"Utilidades", "Armazenamento"},
		"OneDrive":          {"Utilidades", "Armazenamento"},
		"Adobe Reader":      {"Tecnologia", "Software"},
		"Photoshop Express": {"Tecnologia", "Software"},
		"Lightroom":         {"Tecnologia", "Software"},
		"VSCO":              {"Tecnologia", "Software"},
		"PicsArt":           {"Tecnologia", "Software"},
		"Camera360":         {"Tecnologia", "Aplicativo Móvel"},
		"Remini":            {"Tecnologia", "Aplicativo Móvel"},
		"InShot":            {"Tecnologia", "Aplicativo Móvel"},
		"CapCut":            {"Tecnologia", "Aplicativo Móvel"},
		"Kwai":              {"Entretenimento", "Streaming"},
		"GloboPlay":         {"Entretenimento", "Streaming"},
		"Prime Video":       {"Entretenimento", "Streaming"},
		"Disney+":           {"Entretenimento", "Streaming"},
		"HBO Max":           {"Entretenimento", "Streaming"},
		"Paramount+":        {"Entretenimento", "Streaming"},
		"Star+":             {"Entretenimento", "Streaming"},
		"Crunchyroll":       {"Entretenimento", "Streaming"},
		"Twitch":            {"Entretenimento", "Streaming"},
		"Steam":             {"Entretenimento", "Jogos"},
		"Epic Games":        {"Entretenimento", "Jogos"},
		"Coursera":          {"Educação", "Cursos"},
		"Udemy":             {"Educação", "Cursos"},
		"Khan Academy":      {"Educação", "Cursos"},
		"Alura":             {"Educação", "Cursos"},
		"Rocketseat":        {"Educação", "Cursos"},
		"Senai":             {"Educação", "Cursos"},
		"Sebrae":            {"Educação", "Cursos"},
		"Enem":              {"Educação", "Aulas Online"},
		"Estuda.com":        {"Educação", "Aulas Online"},
	}

	// Buscar todos os aplicativos
	var apps []domain.Aplicativo
	db.Find(&apps)

	// Buscar todos os ramos
	var tipoRamos domain.TipoCategoria
	db.Where("nome = ?", "Ramos").First(&tipoRamos)
	var ramos []domain.Categoria
	db.Where("id_tipo_categoria = ?", tipoRamos.ID).Find(&ramos)

	// Buscar todos os subramos
	var tipoSubRamos domain.TipoCategoria
	db.Where("nome = ?", "SubRamos").First(&tipoSubRamos)
	var subramos []domain.Categoria
	db.Where("id_tipo_categoria = ?", tipoSubRamos.ID).Find(&subramos)

	// Indexar por nome para busca rápida
	ramoMap := make(map[string]int64)
	for _, r := range ramos {
		ramoMap[r.Nome] = r.ID
	}
	subramoMap := make(map[string]int64)
	for _, s := range subramos {
		subramoMap[s.Nome] = s.ID
	}

	for _, app := range apps {
		cat, ok := appToCategoria[app.Nome]
		if !ok {
			continue // app não mapeado
		}
		ramoID, ramoOk := ramoMap[cat[0]]
		subramoID, subramoOk := subramoMap[cat[1]]
		if ramoOk {
			appCat := domain.AppCategoria{
				IdApp:       app.ID,
				IdCategoria: ramoID,
			}
			db.Create(&appCat)
		}
		if subramoOk {
			appCat := domain.AppCategoria{
				IdApp:       app.ID,
				IdCategoria: subramoID,
			}
			db.Create(&appCat)
		}
	}
}

// Popula os estágios do aplicativo no banco
func popularEstagiosCatalogo(db *gorm.DB) {
	estagios := []domain.EstagioCatalogo{
		{BaseEntity: domain.BaseEntity{Nome: "Certificação", Descricao: "App em processo de certificação técnica e funcional", Ativo: true}},
		{BaseEntity: domain.BaseEntity{Nome: "Revisão", Descricao: "App em revisão de homologação e documentação", Ativo: true}},
		{BaseEntity: domain.BaseEntity{Nome: "Piloto", Descricao: "App liberado para piloto controlado em campo", Ativo: true}},
		{BaseEntity: domain.BaseEntity{Nome: "Produção", Descricao: "App disponível para uso em produção", Ativo: true}},
	}
	for _, estagio := range estagios {
		var existing domain.EstagioCatalogo
		if err := db.Where("nome = ?", estagio.Nome).First(&existing).Error; err != nil {
			db.Create(&estagio)
		}
	}
}

func popularCatalogoPorEsteira(db *gorm.DB) {
	// Buscar estágios
	var estagios []domain.EstagioCatalogo
	db.Order("id").Find(&estagios)
	if len(estagios) < 4 {
		fmt.Println("Estágios insuficientes para esteira (precisa de 4)")
		return
	}

	// Buscar todos os aplicativos
	var apps []domain.Aplicativo
	db.Find(&apps)

	for _, app := range apps {
		// Buscar todas as configurações do app
		var configuracoes []domain.Configuracao
		db.Where("id_app = ?", app.ID).Find(&configuracoes)

		for _, cfg := range configuracoes {
			// Buscar as 4 versões mais recentes da configuração
			var versoes []domain.VersaoAplicativo
			db.Where("id_configuracao = ?", cfg.ID).Order("created_at DESC").Limit(4).Find(&versoes)

			for i, versao := range versoes {
				if i >= 4 {
					break
				}
				estagio := estagios[i]
				// Criar registro no catálogo
				cat := domain.CatalogoAplicativo{
					IdApp:              app.ID,
					IdTipoIntegracao:   cfg.IdTipoIntegracao,
					IdModeloTerminal:   cfg.IdModeloTerminal,
					IdEstagio:          estagio.ID,
					IdVersaoAplicativo: versao.ID,
					Ativo:              true,
				}
				db.Create(&cat)
			}
		}
	}
}

func gerarJsonAppsPorRegiao(db *gorm.DB) error {
	consultaVersaoApp := func(db *gorm.DB, regiaoID, modeloID, integracaoID int64) ([]map[string]interface{}, error) {
		var result []map[string]interface{}
		err := db.Table("versao_app").
			Select(`versao_app.nome as nome_app, versao_app.tamanho, mdl_trml.nome as modelo_terminal_nome, tip_int.nome as tipo_integracao_nome, cat.nome as categoria_nome, versao_app.id as id_versao`).
			Joins("JOIN cat_app ON cat_app.id_versao_aplicativo = versao_app.id").
			Joins("JOIN cfg ON cfg.id = versao_app.id_configuracao").
			Joins("JOIN app ON app.id = cfg.id_app").
			Joins("JOIN app_cat ON app_cat.id_app = app.id").
			Joins("JOIN cat ON cat.id = app_cat.id_categoria").
			Joins("JOIN mdl_trml ON mdl_trml.id = cfg.id_modelo_terminal").
			Joins("JOIN tip_int ON tip_int.id = cfg.id_tipo_integracao").
			Where("cfg.id_modelo_terminal = ?", modeloID).
			Where("cfg.id_tipo_integracao = ?", integracaoID).
			Where("app_cat.id_categoria = ?", regiaoID).
			Order("versao_app.created_at DESC").
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
		err := db.Table("cat_app").
			Select(`versao_app.nome as nome_app, versao_app.tamanho, mdl_trml.nome as modelo_terminal_nome, tip_int.nome as tipo_integracao_nome, cat.nome as categoria_nome, versao_app.id as id_versao`).
			Joins("JOIN versao_app ON versao_app.id = cat_app.id_versao_aplicativo").
			Joins("JOIN app ON app.id = cat_app.id_app").
			Joins("JOIN app_cat ON app_cat.id_app = cat_app.id_app").
			Joins("JOIN cat ON cat.id = app_cat.id_categoria").
			Joins("JOIN mdl_trml ON mdl_trml.id = cat_app.id_modelo_terminal").
			Joins("JOIN tip_int ON tip_int.id = cat_app.id_tipo_integracao").
			Joins("JOIN est_cat ON est_cat.id = cat_app.id_estagio").
			Where("cat_app.id_modelo_terminal = ?", modeloID).
			Where("cat_app.id_tipo_integracao = ?", integracaoID).
			Where("app_cat.id_categoria = ?", regiaoID).
			Order("cat_app.id_estagio, versao_app.created_at DESC").
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
		Joins("JOIN tip_cat ON tip_cat.id = cat.id_tipo_categoria").
		Where("tip_cat.nome = ?", "Região").
		Pluck("cat.id", &idsRegioes)

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
				fmt.Printf("%s %s, Modelo %s, Integracao %s: %d itens, consulta e exportacao demorou %d ms\n", logPrefix, regiao.Nome, modelo.Nome, integracao.Nome, len(apps), duracao)
			}
		}
	}
	if metricas != nil {
		metricas.Registrar(times, totalItems, totalQueries)
	}
	return nil
}

func TestEntrypoint(t *testing.T) {
	db, err := internal.InitDB()
	if err != nil {
		t.Fatalf("Erro ao inicializar o banco: %v", err)
	}

	// for i := 0; i < 200; i++ {
	// 	cadastrarVersaoApp(db)
	// 	time.Sleep(100 * time.Millisecond)
	// }

	//associarRegioesAosApps(db)

	// categoriasRamo(db)
	// categoriasSubRamo(db)
	//associarRamosESubramosAosApps(db)
	//gerarJsonAppsPorRegiao(db)
	//popularEstagiosCatalogo(db)
	//gerarJsonAppsPorRegiaoCatalogo(db)

	var wg sync.WaitGroup
	metricasGlobal = &MetricasExecucao{} // resetar métricas
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				gerarJsonAppsPorRegiaoCatalogo(db)
				time.Sleep(100 * time.Millisecond)
			}
		}()
	}
	wg.Wait()
	metricasGlobal.CalcularMetricas("gerarJsonAppsPorRegiao")
}
