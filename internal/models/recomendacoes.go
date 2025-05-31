package models

// Estrutura para o FII recomendado
type FIIRecomendado struct {
	Ticker    string
	Nome      string
	Segmento  string
	Tipo      string
	PesoIdeal float64
	Preco     float64
}

// Estrutura para a ação recomendada
type AcaoRecomendada struct {
	Nome      string
	Ticker    string
	PesoIdeal float64
	Preco     float64
}

// Estrutura para ETF recomendado
type ETFRecomendado struct {
	Ticker    string
	Nome      string
	PesoIdeal float64
	Preco     float64
}

// Estrutura para recomendação de compra de FII
type RecomendacaoCompraFII struct {
	Ticker         string
	Nome           string
	Segmento       string
	Tipo           string
	Preco          float64
	PesoAtual      float64
	PesoIdeal      float64
	Diferenca      float64
	Quantidade     int
	ValorCompra    float64
	PesoAposCompra float64
}

// Estrutura para recomendação de compra de ação
type RecomendacaoCompraAcao struct {
	Ticker         string
	Nome           string
	Preco          float64
	PL             float64
	PVP            float64
	DY             float64
	PesoAtual      float64
	PesoIdeal      float64
	Diferenca      float64
	Quantidade     int
	ValorCompra    float64
	PesoAposCompra float64
}

// Estrutura para recomendação de compra de ETF
type RecomendacaoCompraETF struct {
	Ticker         string
	Nome           string
	Preco          float64
	PesoAtual      float64
	PesoIdeal      float64
	Diferenca      float64
	Quantidade     int
	ValorCompra    float64
	PesoAposCompra float64
}

// Estrutura para FII na carteira final
type FIICarteiraFinal struct {
	Ticker            string
	Nome              string
	Segmento          string
	Tipo              string
	Preco             float64
	DY                float64
	PVP               float64
	Quantidade        int
	ValorTotal        float64
	Peso              float64
	DividendosMensais float64
	PesoIdeal         float64
}

// Estrutura para ETF na carteira final
type ETFCarteiraFinal struct {
	Ticker     string
	Nome       string
	Preco      float64
	Quantidade int
	ValorTotal float64
	Peso       float64
	PesoIdeal  float64
}

// Estrutura para ação na carteira final
type AcaoCarteiraFinal struct {
	Ticker     string
	Nome       string
	Preco      float64
	PL         float64
	PVP        float64
	DY         float64
	Quantidade int
	ValorTotal float64
	Peso       float64
	PesoIdeal  float64
}
