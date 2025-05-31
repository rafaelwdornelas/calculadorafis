package models

import "html/template"

// TemplateDados representa os dados enviados ao template HTML
type TemplateDados struct {
	ValorInvestimento                 float64
	ValorTotalCarteira                float64
	ValorTotalCarteiraFII             float64
	ValorTotalCarteiraAcao            float64
	ValorTotalCarteiraETF             float64
	ValorTotalCarteiraRendaFixa       float64
	ValorFuturoCarteira               float64
	RecomendacoesFII                  []RecomendacaoCompraFII
	RecomendacoesAcao                 []RecomendacaoCompraAcao
	RecomendacoesETF                  []RecomendacaoCompraETF
	ValorTotalRecomendadoFII          float64
	ValorTotalRecomendadoAcao         float64
	ValorTotalRecomendadoETF          float64
	ValorTotalRecomendadoFixa         float64
	PercentualRecomendadoFII          float64
	PercentualRecomendadoAcao         float64
	PercentualRecomendadoETF          float64
	PercentualRecomendadoFixa         float64
	ValorRestante                     float64
	CarteiraFinalFII                  []FIICarteiraFinal
	CarteiraFinalAcao                 []AcaoCarteiraFinal
	CarteiraFinalETF                  []ETFCarteiraFinal
	ValorTotalFinalFII                float64
	ValorTotalFinalAcao               float64
	ValorTotalFinalETF                float64
	ValorTotalFinalFixa               float64
	ValorTotalFinal                   float64
	DYMedioPonderadoFII               float64
	DYMedioPonderadoAcao              float64
	DividendosMensaisTotaisFII        float64
	DividendosAnuaisTotalFII          float64
	DividendosAnuaisTotalAcao         float64
	TipoFII                           map[string]map[string]float64
	SegmentoFII                       map[string]map[string]float64
	DistribuicaoAtual                 map[string]float64
	DistribuicaoFinal                 map[string]float64
	DistribuicaoIdeal                 map[string]float64
	AtivosRendaFixa                   []AtivoRendaFixa
	PercentualRendaFixaNoInvestimento float64
	// Novos campos para rendimentos
	CarteiraFinalFIIComRendimento []FIICarteiraFinalComRendimento
	TotalRendimentosMensaisFII    float64
	TotalRendimentosAnuaisFII     float64
	YieldMedioCarteiraFII         float64
}

// FIICarteiraFinalComRendimento representa um FII com informações de rendimento
type FIICarteiraFinalComRendimento struct {
	Ticker           string
	Nome             string
	Segmento         string
	Tipo             string
	Preco            float64
	DY               float64
	PVP              float64
	Quantidade       int
	ValorTotal       float64
	Peso             float64
	PesoIdeal        float64
	UltimoDividendo  float64
	RendimentoMensal float64
	YieldOnCost      float64 // (UltimoDividendo * 12 / Preco) * 100
}

// RespostaCalculadora representa o formato de resposta JSON para o frontend
type RespostaCalculadora struct {
	Status    string        `json:"status"`
	Message   string        `json:"message"`
	DadosHtml template.HTML `json:"dados_html"`
}
