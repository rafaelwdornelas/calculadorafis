package services

import (
	"calculadora-investimentos/internal/models"
	"log"
	"math"
	"strconv"
)

// Calculadora é o serviço que gerencia os cálculos de investimentos
type Calculadora struct {
	distribuicaoService *DistribuidoraService
	recomendacaoService *RecomendadoraService
	otimizadoraService  *OtimizadoraService
	dividendoService    *DividendoService
}

// NewCalculadora cria uma nova instância do serviço de calculadora
func NewCalculadora(
	distribuicaoService *DistribuidoraService,
	recomendacaoService *RecomendadoraService,
	otimizadoraService *OtimizadoraService,
	dividendoService *DividendoService,
) *Calculadora {
	return &Calculadora{
		distribuicaoService: distribuicaoService,
		recomendacaoService: recomendacaoService,
		otimizadoraService:  otimizadoraService,
		dividendoService:    dividendoService,
	}
}

// CalcularRecomendacoes calcula as recomendações de investimento
func (c *Calculadora) CalcularRecomendacoes(
	valorInvestimento float64,
	tiposInvestimento models.TiposInvestimento,
	carteiraFII *models.CarteiraDados,
	carteiraAcao *models.CarteiraAcoes,
	carteiraETF *models.CarteiraETFs,
	carteiraRendaFixa *models.CarteiraRendaFixa,
	recomendadosFII []models.FIIRecomendado,
	recomendadosAcao []models.AcaoRecomendada,
	recomendadosETF []models.ETFRecomendado,
) (*models.TemplateDados, error) {
	// Calcular valor total das carteiras
	valorTotalCarteiraFII := c.calcularValorTotalCarteiraFII(carteiraFII)
	valorTotalCarteiraAcao := c.calcularValorTotalCarteiraAcao(carteiraAcao)
	valorTotalCarteiraETF := c.calcularValorTotalCarteiraETF(carteiraETF)
	valorTotalCarteiraRendaFixa := c.calcularValorTotalCarteiraRendaFixa(carteiraRendaFixa)
	valorTotalCarteira := valorTotalCarteiraFII + valorTotalCarteiraAcao + valorTotalCarteiraETF + valorTotalCarteiraRendaFixa

	// Calcular a distribuição atual
	distribuicaoAtual := c.distribuicaoService.CalcularDistribuicaoAtual(
		valorTotalCarteiraFII,
		valorTotalCarteiraAcao,
		valorTotalCarteiraETF,
		valorTotalCarteiraRendaFixa,
		valorTotalCarteira,
	)

	// Calcular a distribuição ideal (nova distribuição)
	distribuicaoIdeal := c.distribuicaoService.CalcularDistribuicaoIdeal(tiposInvestimento)

	// Calcular a distribuição do novo investimento
	valorTotalAposInvestimento := valorTotalCarteira + valorInvestimento

	// Valores ideais de cada classe de ativo após o investimento
	valorIdealFII := valorTotalAposInvestimento * distribuicaoIdeal["FIIs"] / 100
	valorIdealAcao := valorTotalAposInvestimento * distribuicaoIdeal["Ações"] / 100
	valorIdealETF := valorTotalAposInvestimento * distribuicaoIdeal["ETFs"] / 100
	valorIdealRendaFixa := valorTotalAposInvestimento * distribuicaoIdeal["RendaFixa"] / 100

	// Calcular quanto falta para atingir o valor ideal de cada classe
	valorFaltaFII := math.Max(0, valorIdealFII-valorTotalCarteiraFII)
	valorFaltaAcao := math.Max(0, valorIdealAcao-valorTotalCarteiraAcao)
	valorFaltaETF := math.Max(0, valorIdealETF-valorTotalCarteiraETF)
	valorFaltaRendaFixa := math.Max(0, valorIdealRendaFixa-valorTotalCarteiraRendaFixa)

	// Calcular o valor total que falta para atingir a distribuição ideal
	valorTotalFalta := valorFaltaFII + valorFaltaAcao + valorFaltaETF + valorFaltaRendaFixa

	// Distribuir o investimento
	valorParaFII, valorParaAcao, valorParaETF, valorParaRendaFixa := c.distribuicaoService.DistribuirInvestimento(
		valorInvestimento,
		valorTotalFalta,
		valorFaltaFII,
		valorFaltaAcao,
		valorFaltaETF,
		valorFaltaRendaFixa,
		distribuicaoIdeal,
	)

	// Gerar recomendações
	recomendacoesFII := c.recomendacaoService.GerarRecomendacoesFII(carteiraFII, recomendadosFII, valorParaFII, valorTotalCarteiraFII)
	recomendacoesAcao := c.recomendacaoService.GerarRecomendacoesAcao(carteiraAcao, recomendadosAcao, valorParaAcao, valorTotalCarteiraAcao)
	recomendacoesETF := c.recomendacaoService.GerarRecomendacoesETF(carteiraETF, recomendadosETF, valorParaETF, valorTotalCarteiraETF)

	// Calcular valores totais para o template
	valorTotalRecomendadoFII := 0.0
	for _, rec := range recomendacoesFII {
		valorTotalRecomendadoFII += rec.ValorCompra
	}

	valorTotalRecomendadoAcao := 0.0
	for _, rec := range recomendacoesAcao {
		valorTotalRecomendadoAcao += rec.ValorCompra
	}

	valorTotalRecomendadoETF := 0.0
	for _, rec := range recomendacoesETF {
		valorTotalRecomendadoETF += rec.ValorCompra
	}

	valorTotalRecomendadoFixa := valorParaRendaFixa

	// Calcular o valor total recomendado
	valorTotalRecomendado := valorTotalRecomendadoFII + valorTotalRecomendadoAcao + valorTotalRecomendadoETF + valorTotalRecomendadoFixa
	valorSobra := valorInvestimento - valorTotalRecomendado

	// Otimizar as sobras
	valorRestante := c.otimizadoraService.OtimizarSobras(
		valorSobra,
		recomendadosFII,
		recomendadosAcao,
		recomendadosETF,
		&recomendacoesFII,
		&recomendacoesAcao,
		&recomendacoesETF,
		&valorTotalRecomendadoFII,
		&valorTotalRecomendadoAcao,
		&valorTotalRecomendadoETF,
	)

	// Recalcular os totais após otimização
	valorTotalRecomendado = valorTotalRecomendadoFII + valorTotalRecomendadoAcao + valorTotalRecomendadoETF + valorTotalRecomendadoFixa

	// Calcula os percentuais de cada classe com base no valor total do investimento
	percentualRecomendadoFII := (valorTotalRecomendadoFII / valorInvestimento) * 100
	percentualRecomendadoAcao := (valorTotalRecomendadoAcao / valorInvestimento) * 100
	percentualRecomendadoETF := (valorTotalRecomendadoETF / valorInvestimento) * 100
	percentualRecomendadoFixa := (valorTotalRecomendadoFixa / valorInvestimento) * 100

	// Obter carteira final de FIIs
	carteiraFinalFII := c.recomendacaoService.ObterCarteiraFinalFII(carteiraFII, recomendacoesFII, recomendadosFII, valorTotalCarteiraFII, valorTotalRecomendadoFII)

	// Obter carteira final de ações
	carteiraFinalAcao := c.recomendacaoService.ObterCarteiraFinalAcao(carteiraAcao, recomendacoesAcao, recomendadosAcao, valorTotalCarteiraAcao, valorTotalRecomendadoAcao)

	// Obter carteira final de ETFs
	carteiraFinalETF := c.recomendacaoService.ObterCarteiraFinalETF(carteiraETF, recomendacoesETF, recomendadosETF, valorTotalCarteiraETF, valorTotalRecomendadoETF)

	valorTotalFinalFII := valorTotalCarteiraFII + valorTotalRecomendadoFII
	valorTotalFinalAcao := valorTotalCarteiraAcao + valorTotalRecomendadoAcao
	valorTotalFinalETF := valorTotalCarteiraETF + valorTotalRecomendadoETF
	valorTotalFinalFixa := valorTotalCarteiraRendaFixa + valorTotalRecomendadoFixa
	valorTotalFinal := valorTotalFinalFII + valorTotalFinalAcao + valorTotalFinalETF + valorTotalFinalFixa

	// Calcular a distribuição final
	distribuicaoFinal := map[string]float64{
		"FIIs":      (valorTotalFinalFII / valorTotalFinal) * 100,
		"Ações":     (valorTotalFinalAcao / valorTotalFinal) * 100,
		"ETFs":      (valorTotalFinalETF / valorTotalFinal) * 100,
		"RendaFixa": (valorTotalFinalFixa / valorTotalFinal) * 100,
	}

	// Calcular DY médio ponderado e dividendos mensais totais para FIIs
	dyPonderadoFII := 0.0
	dividendosMensaisTotaisFII := 0.0

	for _, fii := range carteiraFinalFII {
		dyPonderadoFII += fii.DY * (fii.ValorTotal / valorTotalFinalFII)
		dividendosMensaisTotaisFII += fii.DividendosMensais
	}

	// Calcular DY médio ponderado para ações
	dyPonderadoAcao := 0.0
	dividendosAnuaisTotalAcao := 0.0

	for _, acao := range carteiraFinalAcao {
		dyPonderadoAcao += acao.DY * (acao.ValorTotal / valorTotalFinalAcao)
		dividendosAnuaisTotalAcao += (acao.DY / 100.0) * acao.ValorTotal
	}

	// Calcular resumos por tipo e segmento de FII
	tipoFII := make(map[string]map[string]float64)
	segmentoFII := make(map[string]map[string]float64)

	for _, fii := range carteiraFinalFII {
		// Inicializa o mapa para o tipo se não existir
		if _, existe := tipoFII[fii.Tipo]; !existe {
			tipoFII[fii.Tipo] = make(map[string]float64)
		}

		// Inicializa o mapa para o segmento se não existir
		if _, existe := segmentoFII[fii.Segmento]; !existe {
			segmentoFII[fii.Segmento] = make(map[string]float64)
		}

		// Adiciona os valores
		tipoFII[fii.Tipo]["valor"] += fii.ValorTotal
		tipoFII[fii.Tipo]["percentual"] = (tipoFII[fii.Tipo]["valor"] / valorTotalFinalFII) * 100
		tipoFII[fii.Tipo]["dividendos"] += fii.DividendosMensais

		segmentoFII[fii.Segmento]["valor"] += fii.ValorTotal
		segmentoFII[fii.Segmento]["percentual"] = (segmentoFII[fii.Segmento]["valor"] / valorTotalFinalFII) * 100
		segmentoFII[fii.Segmento]["dividendos"] += fii.DividendosMensais
	}

	// Calcular rendimentos dos FIIs usando o novo serviço
	carteiraFinalFIIComRendimento := c.calcularRendimentosFIIs(carteiraFinalFII)

	// Calcular totais de rendimentos
	totalRendimentosMensaisFII := 0.0
	for _, fii := range carteiraFinalFIIComRendimento {
		totalRendimentosMensaisFII += fii.RendimentoMensal
	}
	totalRendimentosAnuaisFII := totalRendimentosMensaisFII * 12

	// Calcular yield médio da carteira
	yieldMedioCarteiraFII := 0.0
	if valorTotalFinalFII > 0 {
		yieldMedioCarteiraFII = (totalRendimentosAnuaisFII / valorTotalFinalFII) * 100
	}

	// Preparar dados para o template
	dados := &models.TemplateDados{
		ValorInvestimento:                 valorInvestimento,
		ValorTotalCarteira:                valorTotalCarteira,
		ValorTotalCarteiraFII:             valorTotalCarteiraFII,
		ValorTotalCarteiraAcao:            valorTotalCarteiraAcao,
		ValorTotalCarteiraETF:             valorTotalCarteiraETF,
		ValorTotalCarteiraRendaFixa:       valorTotalCarteiraRendaFixa,
		ValorFuturoCarteira:               valorTotalCarteira + valorInvestimento,
		RecomendacoesFII:                  recomendacoesFII,
		RecomendacoesAcao:                 recomendacoesAcao,
		RecomendacoesETF:                  recomendacoesETF,
		ValorTotalRecomendadoFII:          valorTotalRecomendadoFII,
		ValorTotalRecomendadoAcao:         valorTotalRecomendadoAcao,
		ValorTotalRecomendadoETF:          valorTotalRecomendadoETF,
		ValorTotalRecomendadoFixa:         valorTotalRecomendadoFixa,
		PercentualRecomendadoFII:          percentualRecomendadoFII,
		PercentualRecomendadoAcao:         percentualRecomendadoAcao,
		PercentualRecomendadoETF:          percentualRecomendadoETF,
		PercentualRecomendadoFixa:         percentualRecomendadoFixa,
		ValorRestante:                     valorRestante,
		CarteiraFinalFII:                  carteiraFinalFII,
		CarteiraFinalAcao:                 carteiraFinalAcao,
		CarteiraFinalETF:                  carteiraFinalETF,
		ValorTotalFinalFII:                valorTotalFinalFII,
		ValorTotalFinalAcao:               valorTotalFinalAcao,
		ValorTotalFinalETF:                valorTotalFinalETF,
		ValorTotalFinalFixa:               valorTotalFinalFixa,
		ValorTotalFinal:                   valorTotalFinal,
		DYMedioPonderadoFII:               dyPonderadoFII,
		DYMedioPonderadoAcao:              dyPonderadoAcao,
		DividendosMensaisTotaisFII:        dividendosMensaisTotaisFII,
		DividendosAnuaisTotalFII:          dividendosMensaisTotaisFII * 12,
		DividendosAnuaisTotalAcao:         dividendosAnuaisTotalAcao,
		TipoFII:                           tipoFII,
		SegmentoFII:                       segmentoFII,
		DistribuicaoAtual:                 distribuicaoAtual,
		DistribuicaoFinal:                 distribuicaoFinal,
		DistribuicaoIdeal:                 distribuicaoIdeal,
		AtivosRendaFixa:                   carteiraRendaFixa.Data,
		PercentualRendaFixaNoInvestimento: (valorTotalRecomendadoFixa / valorInvestimento) * 100,
		// Novos campos
		CarteiraFinalFIIComRendimento: carteiraFinalFIIComRendimento,
		TotalRendimentosMensaisFII:    totalRendimentosMensaisFII,
		TotalRendimentosAnuaisFII:     totalRendimentosAnuaisFII,
		YieldMedioCarteiraFII:         yieldMedioCarteiraFII,
	}

	return dados, nil
}

// calcularRendimentosFIIs calcula os rendimentos mensais dos FIIs
func (c *Calculadora) calcularRendimentosFIIs(carteiraFinal []models.FIICarteiraFinal) []models.FIICarteiraFinalComRendimento {
	var carteiraComRendimento []models.FIICarteiraFinalComRendimento

	log.Println("Iniciando cálculo de rendimentos dos FIIs...")

	for _, fii := range carteiraFinal {
		fiiComRendimento := models.FIICarteiraFinalComRendimento{
			Ticker:     fii.Ticker,
			Nome:       fii.Nome,
			Segmento:   fii.Segmento,
			Tipo:       fii.Tipo,
			Preco:      fii.Preco,
			DY:         fii.DY,
			PVP:        fii.PVP,
			Quantidade: fii.Quantidade,
			ValorTotal: fii.ValorTotal,
			Peso:       fii.Peso,
			PesoIdeal:  fii.PesoIdeal,
		}

		// Buscar o último dividendo
		dividendo, err := c.dividendoService.ObterDividendoFII(fii.Ticker)
		if err != nil {
			log.Printf("Erro ao obter dividendo para %s: %v. Usando valor 0", fii.Ticker, err)
			fiiComRendimento.UltimoDividendo = 0
			fiiComRendimento.RendimentoMensal = 0
			fiiComRendimento.YieldOnCost = 0
		} else {
			fiiComRendimento.UltimoDividendo = dividendo
			fiiComRendimento.RendimentoMensal = dividendo * float64(fii.Quantidade)

			// Calcular Yield on Cost
			if fii.Preco > 0 {
				fiiComRendimento.YieldOnCost = (dividendo * 12 / fii.Preco) * 100
			}

			log.Printf("FII %s: %d cotas x R$ %.2f = R$ %.2f de rendimento mensal (YoC: %.2f%%)",
				fii.Ticker, fii.Quantidade, dividendo, fiiComRendimento.RendimentoMensal, fiiComRendimento.YieldOnCost)
		}

		carteiraComRendimento = append(carteiraComRendimento, fiiComRendimento)
	}

	return carteiraComRendimento
}

// calcularValorTotalCarteiraFII calcula o valor total da carteira de FIIs
func (c *Calculadora) calcularValorTotalCarteiraFII(carteira *models.CarteiraDados) float64 {
	var total float64
	for _, ativo := range carteira.Data {
		total += ativo.CurrentPrice * float64(ativo.Quantity)
	}
	return total
}

// calcularValorTotalCarteiraAcao calcula o valor total da carteira de ações
func (c *Calculadora) calcularValorTotalCarteiraAcao(carteira *models.CarteiraAcoes) float64 {
	var total float64
	for _, ativo := range carteira.Data {
		total += ativo.CurrentPrice * float64(ativo.Quantity)
	}
	return total
}

// calcularValorTotalCarteiraETF calcula o valor total da carteira de ETFs
func (c *Calculadora) calcularValorTotalCarteiraETF(carteira *models.CarteiraETFs) float64 {
	var total float64
	for _, ativo := range carteira.Data {
		currentPrice, err := strconv.ParseFloat(ativo.CurrentPrice, 64)
		if err == nil {
			total += currentPrice * float64(ativo.Quantity)
		}
	}
	return total
}

// calcularValorTotalCarteiraRendaFixa calcula o valor total da carteira de renda fixa
func (c *Calculadora) calcularValorTotalCarteiraRendaFixa(carteira *models.CarteiraRendaFixa) float64 {
	var total float64
	for _, ativo := range carteira.Data {
		valorEquity, err := strconv.ParseFloat(ativo.EquityTotal, 64)
		if err == nil {
			total += valorEquity
		}
	}
	return total
}
