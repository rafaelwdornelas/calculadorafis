package models

// Estruturas para os dados da carteira atual de FIIs
type CarteiraDados struct {
	Total    int        `json:"total"`
	Data     []AtivoFII `json:"data"`
	Draw     int        `json:"draw"`
	Weighted float64    `json:"weighted"`
}

type AtivoFII struct {
	ID             int     `json:"id"`
	Quantity       int     `json:"quantity"`
	TickerName     string  `json:"ticker_name"`
	AvgPrice       float64 `json:"avg_price"`
	CurrentPrice   float64 `json:"current_price"`
	EquityBRL      string  `json:"equity_brl"`
	Appreciation   float64 `json:"appreciation"`
	WeightedReturn float64 `json:"weighted_return"`
	PercentWallet  float64 `json:"percent_wallet"`
	PercentIdeal   float64 `json:"percent_ideal"`
	Buy            string  `json:"buy"`
	Segment        string  `json:"segment"`
	FiiType        string  `json:"fii_type"`
	PVP            string  `json:"p_vp"`
	DY             string  `json:"dy"`
	YOC            float64 `json:"yoc"`
}

// Estrutura para a carteira de ETFs
type CarteiraETFs struct {
	Total    int        `json:"total"`
	Data     []AtivoETF `json:"data"`
	Draw     int        `json:"draw"`
	Weighted float64    `json:"weighted"`
}

type AtivoETF struct {
	ID             int     `json:"id"`
	Quantity       int     `json:"quantity"`
	TickerName     string  `json:"ticker_name"`
	AvgPrice       float64 `json:"avg_price"`
	CurrentPrice   string  `json:"current_price"`
	EquityBRL      string  `json:"equity_brl"`
	Appreciation   float64 `json:"appreciation"`
	WeightedReturn float64 `json:"weighted_return"`
	PercentWallet  float64 `json:"percent_wallet"`
	PercentIdeal   float64 `json:"percent_ideal"`
	Buy            string  `json:"buy"`
	TickerType     string  `json:"ticker_type"`
	Rating         string  `json:"rating"`
	RawRating      string  `json:"raw_rating"`
	EquityTotal    float64 `json:"equity_total"`
}

// Estrutura para a carteira de ações
type CarteiraAcoes struct {
	Total int         `json:"total"`
	Data  []AtivoAcao `json:"data"`
	Draw  int         `json:"draw"`
}

type AtivoAcao struct {
	ID             int     `json:"id"`
	Quantity       int     `json:"quantity"`
	TickerName     string  `json:"ticker_name"`
	AvgPrice       float64 `json:"avg_price"`
	CurrentPrice   float64 `json:"current_price"`
	EquityBRL      string  `json:"equity_brl"`
	Appreciation   float64 `json:"appreciation"`
	WeightedReturn float64 `json:"weighted_return"`
	PercentWallet  float64 `json:"percent_wallet"`
	PercentIdeal   float64 `json:"percent_ideal"`
	Buy            string  `json:"buy"`
	PL             string  `json:"p_l"`
	PVP            string  `json:"p_vp"`
	DY             string  `json:"dy"`
}

// Estrutura para a carteira de renda fixa
type CarteiraRendaFixa struct {
	Total int              `json:"total"`
	Data  []AtivoRendaFixa `json:"data"`
	Draw  int              `json:"draw"`
}

type AtivoRendaFixa struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Ticker         string  `json:"ticker"`
	Quantity       string  `json:"quantity"`
	EquityBRL      string  `json:"equity_brl"`
	Appreciation   float64 `json:"appreciation"`
	WeightedReturn float64 `json:"weighted_return"`
	PercentWallet  string  `json:"percent_wallet"`
	PercentIdeal   string  `json:"percent_ideal"`
	EquityTotal    string  `json:"equity_total"`
	Applied        string  `json:"applied"`
	Buy            string  `json:"buy"`
	Indexer        string  `json:"indexer"`
	Emitter        string  `json:"emitter"`
	InvestmentType string  `json:"investment_type"`
	RateType       string  `json:"rate_type"`
	PercentageCDI  string  `json:"percentage_cdi"`
	PercentageYear string  `json:"percentage_year"`
	DueDate        string  `json:"due_date"`
	DailyLiquidity int     `json:"dailyLiquidity"`
	TickerType     string  `json:"ticker_type"`
}

// TiposInvestimento representa os tipos de investimento selecionados pelo usuário
type TiposInvestimento struct {
	FIIs      bool
	Acoes     bool
	ETFs      bool
	RendaFixa bool
}
