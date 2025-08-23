package pkg

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"store-go/shared/internal"
	"store-go/shared/internal/domain"
	"testing"
	"time"

	"gorm.io/gorm"
)

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

func gerarJsonAppsPorRegiao(db *gorm.DB) error {
	// Buscar todos os IDs das regiões
	var idsRegioes []int64
	db.Model(&domain.Categoria{}).
		Joins("JOIN tip_cat ON tip_cat.id = cat.id_tipo_categoria").
		Where("tip_cat.nome = ?", "Região").
		Pluck("cat.id", &idsRegioes)

	// Buscar o par modelo + integração desejado
	var modelo domain.ModeloTerminal
	db.Where("nome = ?", "L400").First(&modelo)
	var integracao domain.TipoIntegracao
	db.Where("nome = ?", "ADQ").First(&integracao)

	for _, regiaoID := range idsRegioes {
		start := time.Now()
		// Buscar nome da região
		var regiao domain.Categoria
		db.First(&regiao, regiaoID)

		// Buscar versões dos apps relacionados à região, modelo e integração usando apenas uma chamada GORM
		var appsVersao []map[string]interface{}
		db.Table("versao_app").
			Select(`versao_app.*, cfg.id_modelo_terminal, cfg.id_tipo_integracao, mdl_trml.nome as modelo_terminal_nome, tip_int.nome as tipo_integracao_nome`).
			Joins("JOIN cfg ON cfg.id = versao_app.id_configuracao").
			Joins("JOIN app_cat ON app_cat.id_app = cfg.id_app").
			Joins("JOIN mdl_trml ON mdl_trml.id = cfg.id_modelo_terminal").
			Joins("JOIN tip_int ON tip_int.id = cfg.id_tipo_integracao").
			Where("app_cat.id_categoria = ?", regiaoID).
			Where("cfg.id_modelo_terminal = ?", modelo.ID).
			Where("cfg.id_tipo_integracao = ?", integracao.ID).
			Order("versao_app.created_at DESC").
			Limit(1000).
			Scan(&appsVersao)

		// Gerar arquivo JSON
		fileName := fmt.Sprintf("apps_regiao_%s.json", regiao.Nome)
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(appsVersao); err != nil {
			return err
		}
		duracao := time.Since(start).Milliseconds()
		fmt.Printf("Regiao %s: consulta e exportacao demorou %d ms\n", regiao.Nome, duracao)
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
	gerarJsonAppsPorRegiao(db)
}
