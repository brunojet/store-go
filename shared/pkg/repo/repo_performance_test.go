// Para cada configuração de cada app, resgata as 4 versões mais recentes e insere no catálogo com os estágios
package repo

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/brunojet/store-go/shared/pkg/domain"

	"gorm.io/driver/postgres"
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

func cadastroTipoCategoria(db *gorm.DB) {

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

func cadastroCategoriasRegiao(db *gorm.DB) {
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

func cadastroCategoriasRamo(db *gorm.DB) {
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

func cadastroCategoriasSubRamo(db *gorm.DB) {
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

func casdastroModelosTerminal(db *gorm.DB) {
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

func cadastroTiposIntegracao(db *gorm.DB) {
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

func cadastroAplicativos(db *gorm.DB) {
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

func cadastroConfiguraces(db *gorm.DB) {
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

func updateCadastroCategoria(db *gorm.DB, categoria domain.Categoria, baseModel domain.BaseModel) error {
	// Associa categoria ao cadastro via GORM M2M
	cadastro := domain.Cadastro{BaseModel: baseModel}
	if err := db.Model(&cadastro).Association("Categorias").Append(&categoria); err != nil {
		return err
	}
	return nil
}

func associarCategoriasAoCadastro(db *gorm.DB, categorias []domain.Categoria, baseModel domain.BaseModel) error {
	n := len(categorias)
	if n < 2 {
		return nil
	}
	idxs := rand.Perm(n)[:2]

	for _, idx := range idxs {
		if err := updateCadastroCategoria(db, categorias[idx], baseModel); err != nil {
			return err
		}
	}

	return nil
}

func associarRamoESubRamoAoCadastro(db *gorm.DB, nome string, ramos, subramos []domain.Categoria, baseModel domain.BaseModel) error {
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

	ramoSubRamoApp := appToCategoria[nome]
	for _, r := range ramos {
		if r.Nome == ramoSubRamoApp[0] {
			if err := updateCadastroCategoria(db, r, baseModel); err != nil {
				return err
			}
			break
		}
	}

	for _, sr := range subramos {
		if sr.Nome == ramoSubRamoApp[1] {
			if err := updateCadastroCategoria(db, sr, baseModel); err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func associarConfiguracaoAoCadastro(db *gorm.DB, config domain.Configuracao, baseModel domain.BaseModel) error {
	// Associa configuração ao cadastro via GORM M2M
	cadastro := domain.Cadastro{BaseModel: baseModel}
	if err := db.Model(&cadastro).Association("Configuracoes").Append(&config); err != nil {
		return err
	}
	return nil
}

func solicitarCadastroAplicativo(db *gorm.DB) error {
	var tipoRegiao domain.TipoCategoria
	db.Where("nome = ?", "Região").First(&tipoRegiao)
	var regioes []domain.Categoria
	db.Where("id_tipo_categoria = ?", tipoRegiao.ID).Find(&regioes)
	var tipoRamos domain.TipoCategoria
	db.Where("nome = ?", "Ramos").First(&tipoRamos)
	var ramos []domain.Categoria
	db.Where("id_tipo_categoria = ?", tipoRamos.ID).Find(&ramos)
	var tipoSubRamos domain.TipoCategoria
	db.Where("nome = ?", "SubRamos").First(&tipoSubRamos)
	var subramos []domain.Categoria
	db.Where("id_tipo_categoria = ?", tipoSubRamos.ID).Find(&subramos)

	var apps []domain.Aplicativo
	db.Find(&apps)

	for _, app := range apps {
		// Buscar todas as configurações do app
		var configuracoes []domain.Configuracao
		db.Preload("ModeloTerminal", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nome")
		}).Preload("TipoIntegracao", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nome")
		}).Where("id_app = ?", app.ID).Find(&configuracoes)

		cadastroConfig := domain.Cadastro{
			IdApp: app.ID,
			Contato: domain.Contato{
				BaseEntity: domain.BaseEntity{
					Nome: "Contato " + app.Nome,
				},
				RazaoSocial: app.Nome + " S.A.",
				Email:       "contato@" + app.Nome + ".com",
				Telefone:    "1234-5678",
				Site:        "www." + app.Nome + ".com",
			},
			DetalhesAplicativo: domain.DetalhesAplicativo{
				Descricao: fmt.Sprintf(
					"Cadastro para %s",
					app.Nome),
			},
		}
		if tx := db.Create(&cadastroConfig); tx.Error != nil {
			return tx.Error
		}
		if err := associarCategoriasAoCadastro(db, regioes, cadastroConfig.BaseModel); err != nil {
			return err
		}
		if err := associarRamoESubRamoAoCadastro(db, app.Nome, ramos, subramos, cadastroConfig.BaseModel); err != nil {
			return err
		}

		for _, config := range configuracoes {
			if err := associarConfiguracaoAoCadastro(db, config, cadastroConfig.BaseModel); err != nil {
				return err
			}
		}
	}

	return nil
}

func solicitarCadastroVersaoAplicativo(db *gorm.DB) error {
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
		for _, config := range cadastro.Configuracoes {
			now := time.Now()
			year := now.Year() % 100
			dayOfYear := now.YearDay()
			seconds := now.Hour()*3600 + now.Minute()*60 + now.Second()
			nomeVersao := fmt.Sprintf("%02d%03d%05d", year, dayOfYear, seconds)
			tamanho := int64(rand.Intn(190)+10) * 1024 * 1024 // 10MB a 200MB

			versao := domain.VersaoAplicativo{
				IdConfiguracao: config.ID,
				BaseEntity:     domain.BaseEntity{Nome: app.Nome},
				Tamanho:        tamanho,
				NomeVersao:     nomeVersao,
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

func getCadastroAplicativo(db *gorm.DB, appID int64) (domain.Cadastro, error) {
	var cadastro domain.Cadastro
	if err := db.Where("id_app = ?", appID).
		Order("created_at desc").
		Preload("Configuracoes").
		First(&cadastro).Error; err != nil {
		return domain.Cadastro{}, err
	}
	return cadastro, nil
}

func solicitarPublicacaoVersaoAplicativo(db *gorm.DB) error {
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
		for _, config := range cadastro.Configuracoes {
			var versaoAtual domain.VersaoAplicativo
			if err := db.Where("id_configuracao = ?", config.ID).Order("created_at desc").First(&versaoAtual).Error; err != nil || versaoAtual.ID == 0 {
				continue
			}

			var catalogos []domain.CatalogoAplicativo

			db.Preload("Estagio").Where("id_configuracao = ?", config.ID).Order("id_estagio DESC").Find(&catalogos)

			lenCatalogos := len(catalogos)
			lenEstagios := len(estagios)

			// Cria novo catálogo se ainda não percorreu todos os estágios
			if lenCatalogos < lenEstagios {
				novoCatalogo := domain.CatalogoAplicativo{
					IdCadastro:     cadastro.ID,
					IdConfiguracao: config.ID,
					Ativo:          true,
					IdEstagio:      estagios[lenCatalogos].ID,
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

// Popula os estágios do aplicativo no banco
func popularEstagiosCatalogo(db *gorm.DB) {
	estagios := []domain.Estagio{
		{BaseEntity: domain.BaseEntity{Nome: "Certificação", Descricao: "App em processo de certificação técnica e funcional", Ativo: true}},
		{BaseEntity: domain.BaseEntity{Nome: "Revisão", Descricao: "App em revisão de homologação e documentação", Ativo: true}},
		{BaseEntity: domain.BaseEntity{Nome: "Piloto", Descricao: "App liberado para piloto controlado em campo", Ativo: true}},
		{BaseEntity: domain.BaseEntity{Nome: "Produção", Descricao: "App disponível para uso em produção", Ativo: true}},
	}
	for _, estagio := range estagios {
		var existing domain.Estagio
		if err := db.Where("nome = ?", estagio.Nome).First(&existing).Error; err != nil {
			db.Create(&estagio)
		}
	}
}

func gerarJsonAppsPorRegiao(db *gorm.DB) error {
	consultaVersaoApp := func(db *gorm.DB, regiaoID, modeloID, integracaoID int64) ([]map[string]interface{}, error) {
		var result []map[string]interface{}
		err := db.Debug().Table("versao_app").
			Select(`versao_app.nome as nome_app, versao_app.tamanho, mdl_trml.nome as modelo_terminal_nome, tip_int.nome as tipo_integracao_nome, cat.nome as categoria_nome, versao_app.id as id_versao`).
			Joins("JOIN cat_app ON cat_app.id_versao_aplicativo = versao_app.id"). //Localiza a versao no catalogo
			Joins("JOIN cfg ON cfg.id = versao_app.id_configuracao").              //Localiza a configuracao
			Joins("JOIN cad ON cad.id = cat_app.id_cadastro").                     //Localiza o cadastro
			Joins("JOIN cat_cad ON cat_cad.cadastro_id = cad.id").                 //M2M cadastro-categoria
			Joins("JOIN cat ON cat.id = cat_cad.categoria_id").                    //Localiza a categoria
			Joins("JOIN mdl_trml ON mdl_trml.id = cfg.id_modelo_terminal").        //Localiza o modelo do terminal
			Joins("JOIN tip_int ON tip_int.id = cfg.id_tipo_integracao").          //Localiza o tipo de integração
			Where("cfg.id_modelo_terminal = ?", modeloID).
			Where("cfg.id_tipo_integracao = ?", integracaoID).
			Where("cat.id = ?", regiaoID).
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
			Joins("JOIN versao_app ON versao_app.id = cat_app.id_versao_aplicativo"). //Localiza a versão no catálogo
			Joins("JOIN cfg ON cfg.id = versao_app.id_configuracao").                 //Localiza a configuração
			Joins("JOIN cad ON cad.id = cat_app.id_cadastro").                        //Localiza o cadastro
			Joins("JOIN cat_cad ON cat_cad.cadastro_id = cad.id").                    //M2M cadastro-categoria
			Joins("JOIN cat ON cat.id = cat_cad.categoria_id").                       //Localiza a categoria
			Joins("JOIN mdl_trml ON mdl_trml.id = cfg.id_modelo_terminal").           //Localiza o modelo do terminal
			Joins("JOIN tip_int ON tip_int.id = cfg.id_tipo_integracao").             //Localiza o tipo de integração
			Where("cfg.id_modelo_terminal = ?", modeloID).
			Where("cfg.id_tipo_integracao = ?", integracaoID).
			Where("cat.id = ?", regiaoID).
			Order("versao_app.created_at DESC").
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
		cadastroTipoCategoria(db)
		cadastroCategoriasRegiao(db)
		popularEstagiosCatalogo(db)
		solicitarCadastroAplicativo(db)
		cadastroCategoriasRamo(db)
		cadastroCategoriasSubRamo(db)
		casdastroModelosTerminal(db)
		cadastroTiposIntegracao(db)
		cadastroAplicativos(db)
		cadastroConfiguraces(db)
	}

	loadVersao := false

	if loadVersao {
		for i := 0; i < 100; i++ {
			solicitarCadastroVersaoAplicativo(db)
			solicitarPublicacaoVersaoAplicativo(db)
		}
	}

	//	gerarJsonAppsPorRegiao(db)
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
