package services

import (
	"calculadora-investimentos/internal/models"
	"log"
	"sort"
	"strconv"
	"strings"
)

// RecomendadoraService gerencia as recomendações de investimentos
type RecomendadoraService struct {
	dataComService *DataComService
}

// Atualize o construtor
func NewRecomendadoraService() *RecomendadoraService {
	return &RecomendadoraService{
		dataComService: NewDataComService(),
	}
}

// GerarRecomendacoesFII gera recomendações de compra para FIIs
func (s *RecomendadoraService) GerarRecomendacoesFII(
	carteira *models.CarteiraDados,
	recomendados []models.FIIRecomendado,
	valorInvestimento, valorTotalCarteira float64,
) []models.RecomendacaoCompraFII {
	var recomendacoes []models.RecomendacaoCompraFII

	// Valor total futuro da carteira
	valorTotalFuturo := valorTotalCarteira + valorInvestimento

	// Mapa para facilitar a busca de ativos na carteira
	ativosNaCarteira := make(map[string]models.AtivoFII)
	for _, ativo := range carteira.Data {
		ativosNaCarteira[ativo.TickerName] = ativo
	}

	// Calcular o valor atual alocado em cada FII
	valorAtualFII := make(map[string]float64)
	for _, ativo := range carteira.Data {
		valorAtualFII[ativo.TickerName] = ativo.CurrentPrice * float64(ativo.Quantity)
	}

	// Calcular o valor futuro ideal para cada FII com base nos pesos
	valorIdealFII := make(map[string]float64)
	for _, rec := range recomendados {
		// Valor ideal baseado no peso
		valorIdealFII[rec.Ticker] = (rec.PesoIdeal / 100) * valorTotalFuturo
	}

	// Calcular quanto precisa ser comprado de cada FII
	valorCompraFII := make(map[string]float64)
	for _, rec := range recomendados {
		atual := valorAtualFII[rec.Ticker]
		ideal := valorIdealFII[rec.Ticker]

		// Só recomenda compra se o valor ideal for maior que o atual
		if ideal > atual {
			valorCompraFII[rec.Ticker] = ideal - atual
		}
	}

	// Verificar se o valor total de compra excede o disponível
	valorTotalCompra := 0.0
	for _, valor := range valorCompraFII {
		valorTotalCompra += valor
	}

	// Se exceder, ajustar proporcionalmente
	if valorTotalCompra > valorInvestimento {
		fatorAjuste := valorInvestimento / valorTotalCompra
		for ticker := range valorCompraFII {
			valorCompraFII[ticker] *= fatorAjuste
		}
	}

	// Calcular a quantidade a ser comprada de cada FII
	for _, rec := range recomendados {
		valorCompra := valorCompraFII[rec.Ticker]

		if valorCompra > 0 {
			// Calcula quantidade a ser comprada (arredonda para baixo)
			quantidadeCompra := int(valorCompra / rec.Preco)

			// Recalcula o valor da compra com a quantidade ajustada
			valorCompraAjustado := float64(quantidadeCompra) * rec.Preco

			// Só recomenda se a quantidade for maior que zero
			if quantidadeCompra > 0 {
				// Calcula o peso atual
				var pesoAtual float64
				if ativo, existe := ativosNaCarteira[rec.Ticker]; existe {
					pesoAtual = ativo.PercentWallet
				}

				// Calcula o peso após a compra
				valorAtual := valorAtualFII[rec.Ticker]
				pesoAposCompra := ((valorAtual + valorCompraAjustado) / valorTotalFuturo) * 100

				recomendacao := models.RecomendacaoCompraFII{
					Ticker:         rec.Ticker,
					Nome:           rec.Nome,
					Segmento:       rec.Segmento,
					Tipo:           rec.Tipo,
					Preco:          rec.Preco,
					PesoAtual:      pesoAtual,
					PesoIdeal:      rec.PesoIdeal,
					Diferenca:      rec.PesoIdeal - pesoAtual,
					Quantidade:     quantidadeCompra,
					ValorCompra:    valorCompraAjustado,
					PesoAposCompra: pesoAposCompra,
				}

				// ADICIONE ESTE CÓDIGO PARA ANÁLISE DE DATA COM
				log.Printf("Analisando data com para FII: %s", rec.Ticker)
				if analiseDataCom, err := s.dataComService.AnalisarDataComTicker(rec.Ticker, "FII"); err == nil {
					if analiseDataCom != nil {
						recomendacao.ProximaDataCom = analiseDataCom.ProximaDataCom.Format("02/01/2006")
						recomendacao.DiasAteDataCom = analiseDataCom.DiasAteDataCom
						recomendacao.StatusCompra = analiseDataCom.StatusCompra
						recomendacao.MensagemStatus = analiseDataCom.MensagemStatus
						log.Printf("Data com para %s: %s (Status: %s)", rec.Ticker, recomendacao.ProximaDataCom, recomendacao.StatusCompra)
					}
				} else {
					log.Printf("Erro ao analisar data com para %s: %v", rec.Ticker, err)
					recomendacao.StatusCompra = "INDISPONIVEL"
					recomendacao.MensagemStatus = "Dados não disponíveis"
				}

				recomendacoes = append(recomendacoes, recomendacao)
			}
		}
	}

	return recomendacoes
}

// GerarRecomendacoesAcao gera recomendações de compra para ações
func (s *RecomendadoraService) GerarRecomendacoesAcao(
	carteira *models.CarteiraAcoes,
	recomendados []models.AcaoRecomendada,
	valorInvestimento, valorTotalCarteira float64,
) []models.RecomendacaoCompraAcao {
	var recomendacoes []models.RecomendacaoCompraAcao

	// Valor total futuro da carteira
	valorTotalFuturo := valorTotalCarteira + valorInvestimento

	// Mapa para facilitar a busca de ativos na carteira
	ativosNaCarteira := make(map[string]models.AtivoAcao)
	for _, ativo := range carteira.Data {
		ativosNaCarteira[ativo.TickerName] = ativo
	}

	// Calcular o valor atual alocado em cada ação
	valorAtualAcao := make(map[string]float64)
	for _, ativo := range carteira.Data {
		valorAtualAcao[ativo.TickerName] = ativo.CurrentPrice * float64(ativo.Quantity)
	}

	// Calcular o valor futuro ideal para cada ação com base nos pesos
	valorIdealAcao := make(map[string]float64)
	for _, rec := range recomendados {
		// Valor ideal baseado no peso
		valorIdealAcao[rec.Ticker] = (rec.PesoIdeal / 100) * valorTotalFuturo
	}

	// Calcular quanto precisa ser comprado de cada ação
	valorCompraAcao := make(map[string]float64)
	for _, rec := range recomendados {
		atual := valorAtualAcao[rec.Ticker]
		ideal := valorIdealAcao[rec.Ticker]

		// Só recomenda compra se o valor ideal for maior que o atual
		if ideal > atual {
			valorCompraAcao[rec.Ticker] = ideal - atual
		}
	}

	// Verificar se o valor total de compra excede o disponível
	valorTotalCompra := 0.0
	for _, valor := range valorCompraAcao {
		valorTotalCompra += valor
	}

	// Se exceder, ajustar proporcionalmente
	if valorTotalCompra > valorInvestimento {
		fatorAjuste := valorInvestimento / valorTotalCompra
		for ticker := range valorCompraAcao {
			valorCompraAcao[ticker] *= fatorAjuste
		}
	}

	// Calcular a quantidade a ser comprada de cada ação
	for _, rec := range recomendados {
		valorCompra := valorCompraAcao[rec.Ticker]

		if valorCompra > 0 {
			// Calcula quantidade a ser comprada (arredonda para baixo)
			quantidadeCompra := int(valorCompra / rec.Preco)

			// Recalcula o valor da compra com a quantidade ajustada
			valorCompraAjustado := float64(quantidadeCompra) * rec.Preco

			// Só recomenda se a quantidade for maior que zero
			if quantidadeCompra > 0 {
				// Calcula o peso atual
				var pesoAtual float64
				if ativo, existe := ativosNaCarteira[rec.Ticker]; existe {
					pesoAtual = ativo.PercentWallet
				}

				// Calcula o peso após a compra
				valorAtual := valorAtualAcao[rec.Ticker]
				pesoAposCompra := ((valorAtual + valorCompraAjustado) / valorTotalFuturo) * 100

				recomendacao := models.RecomendacaoCompraAcao{
					Ticker:         rec.Ticker,
					Nome:           rec.Nome,
					Preco:          rec.Preco,
					PesoAtual:      pesoAtual,
					PesoIdeal:      rec.PesoIdeal,
					Diferenca:      rec.PesoIdeal - pesoAtual,
					Quantidade:     quantidadeCompra,
					ValorCompra:    valorCompraAjustado,
					PesoAposCompra: pesoAposCompra,
				}

				// ADICIONE ESTE CÓDIGO PARA ANÁLISE DE DATA COM
				log.Printf("Analisando data com para Ação: %s", rec.Ticker)
				if analiseDataCom, err := s.dataComService.AnalisarDataComTicker(rec.Ticker, "ACAO"); err == nil {
					if analiseDataCom != nil {
						recomendacao.ProximaDataCom = analiseDataCom.ProximaDataCom.Format("02/01/2006")
						recomendacao.DiasAteDataCom = analiseDataCom.DiasAteDataCom
						recomendacao.StatusCompra = analiseDataCom.StatusCompra
						recomendacao.MensagemStatus = analiseDataCom.MensagemStatus
						log.Printf("Data com para %s: %s (Status: %s)", rec.Ticker, recomendacao.ProximaDataCom, recomendacao.StatusCompra)
					}
				} else {
					log.Printf("Erro ao analisar data com para %s: %v", rec.Ticker, err)
					recomendacao.StatusCompra = "INDISPONIVEL"
					recomendacao.MensagemStatus = "Dados não disponíveis"
				}

				recomendacoes = append(recomendacoes, recomendacao)
			}
		}
	}

	return recomendacoes
}

// GerarRecomendacoesETF gera recomendações de compra para ETFs
func (s *RecomendadoraService) GerarRecomendacoesETF(
	carteira *models.CarteiraETFs,
	recomendados []models.ETFRecomendado,
	valorInvestimento, valorTotalCarteira float64,
) []models.RecomendacaoCompraETF {
	var recomendacoes []models.RecomendacaoCompraETF

	// Valor total futuro da carteira
	valorTotalFuturo := valorTotalCarteira + valorInvestimento

	// Mapa para facilitar a busca de ativos na carteira
	ativosNaCarteira := make(map[string]models.AtivoETF)
	for _, ativo := range carteira.Data {
		ativosNaCarteira[ativo.TickerName] = ativo
	}

	// Calcular o valor atual alocado em cada ETF
	valorAtualETF := make(map[string]float64)
	for _, ativo := range carteira.Data {
		currentPrice, err := strconv.ParseFloat(ativo.CurrentPrice, 64)
		if err == nil {
			valorAtualETF[ativo.TickerName] = currentPrice * float64(ativo.Quantity)
		}
	}

	// Calcular o valor futuro ideal para cada ETF com base nos pesos
	valorIdealETF := make(map[string]float64)
	for _, rec := range recomendados {
		// Valor ideal baseado no peso
		valorIdealETF[rec.Ticker] = (rec.PesoIdeal / 100) * valorTotalFuturo
	}

	// Calcular quanto precisa ser comprado de cada ETF
	valorCompraETF := make(map[string]float64)
	for _, rec := range recomendados {
		atual := valorAtualETF[rec.Ticker]
		ideal := valorIdealETF[rec.Ticker]

		// Só recomenda compra se o valor ideal for maior que o atual
		if ideal > atual {
			valorCompraETF[rec.Ticker] = ideal - atual
		}
	}

	// Verificar se o valor total de compra excede o disponível
	valorTotalCompra := 0.0
	for _, valor := range valorCompraETF {
		valorTotalCompra += valor
	}

	// Se exceder, ajustar proporcionalmente
	// Se exceder, ajustar proporcionalmente
	if valorTotalCompra > valorInvestimento {
		fatorAjuste := valorInvestimento / valorTotalCompra
		for ticker := range valorCompraETF {
			valorCompraETF[ticker] *= fatorAjuste
		}
	}

	// Calcular a quantidade a ser comprada de cada ETF
	for _, rec := range recomendados {
		valorCompra := valorCompraETF[rec.Ticker]

		if valorCompra > 0 {
			// Calcula quantidade a ser comprada (arredonda para baixo)
			quantidadeCompra := int(valorCompra / rec.Preco)

			// Recalcula o valor da compra com a quantidade ajustada
			valorCompraAjustado := float64(quantidadeCompra) * rec.Preco

			// Só recomenda se a quantidade for maior que zero
			if quantidadeCompra > 0 {
				// Calcula o peso atual
				var pesoAtual float64
				if ativo, existe := ativosNaCarteira[rec.Ticker]; existe {
					pesoAtual = ativo.PercentWallet
				}

				// Calcula o peso após a compra
				valorAtual := valorAtualETF[rec.Ticker]
				pesoAposCompra := ((valorAtual + valorCompraAjustado) / valorTotalFuturo) * 100

				recomendacao := models.RecomendacaoCompraETF{
					Ticker:         rec.Ticker,
					Nome:           rec.Nome,
					Preco:          rec.Preco,
					PesoAtual:      pesoAtual,
					PesoIdeal:      rec.PesoIdeal,
					Diferenca:      rec.PesoIdeal - pesoAtual,
					Quantidade:     quantidadeCompra,
					ValorCompra:    valorCompraAjustado,
					PesoAposCompra: pesoAposCompra,
				}

				recomendacoes = append(recomendacoes, recomendacao)
			}
		}
	}

	return recomendacoes
}

// ObterCarteiraFinalFII obtém a carteira final de FIIs após as compras sugeridas
func (s *RecomendadoraService) ObterCarteiraFinalFII(
	carteira *models.CarteiraDados,
	recomendacoes []models.RecomendacaoCompraFII,
	recomendados []models.FIIRecomendado,
	valorTotalCarteira, valorTotalRecomendado float64,
) []models.FIICarteiraFinal {
	// Mapa para facilitar a busca de ativos na carteira
	ativosNaCarteira := make(map[string]*models.AtivoFII)
	for i := range carteira.Data {
		ativosNaCarteira[carteira.Data[i].TickerName] = &carteira.Data[i]
	}

	// Mapa para facilitar a busca de informações de recomendações
	infoRecomendados := make(map[string]models.FIIRecomendado)
	for _, rec := range recomendados {
		infoRecomendados[rec.Ticker] = rec
	}

	// Lista para armazenar os FIIs da carteira final
	var carteiraFinal []models.FIICarteiraFinal

	// Primeiro, adiciona os FIIs que já estão na carteira
	for _, ativo := range carteira.Data {
		dy, _ := strconv.ParseFloat(strings.TrimSuffix(ativo.DY, "%"), 64)
		pvp, _ := strconv.ParseFloat(strings.Replace(ativo.PVP, ",", ".", -1), 64)

		// Verifica se há compra adicional deste FII
		quantidade := ativo.Quantity
		for _, rec := range recomendacoes {
			if rec.Ticker == ativo.TickerName {
				quantidade += rec.Quantidade
				break
			}
		}

		valorTotal := float64(quantidade) * ativo.CurrentPrice

		// Calcula dividendos mensais (DY anual / 12 * valor total)
		dividendosMensais := (dy / 100 / 12) * valorTotal

		// Verifica se tem peso ideal nas recomendações
		var pesoIdeal float64
		if rec, existe := infoRecomendados[ativo.TickerName]; existe {
			pesoIdeal = rec.PesoIdeal
		}

		fii := models.FIICarteiraFinal{
			Ticker:            ativo.TickerName,
			Nome:              ativo.TickerName, // Poderia buscar o nome completo em recomendados
			Segmento:          ativo.Segment,
			Tipo:              ativo.FiiType,
			Preco:             ativo.CurrentPrice,
			DY:                dy,
			PVP:               pvp,
			Quantidade:        quantidade,
			ValorTotal:        valorTotal,
			Peso:              0, // Será calculado depois
			DividendosMensais: dividendosMensais,
			PesoIdeal:         pesoIdeal,
		}

		carteiraFinal = append(carteiraFinal, fii)
	}

	// Adiciona os novos FIIs que não estavam na carteira original
	for _, rec := range recomendacoes {
		// Verifica se já existe na carteira
		if _, existe := ativosNaCarteira[rec.Ticker]; !existe {
			valorTotal := float64(rec.Quantidade) * rec.Preco

			fii := models.FIICarteiraFinal{
				Ticker:            rec.Ticker,
				Nome:              rec.Nome,
				Segmento:          rec.Segmento,
				Tipo:              rec.Tipo,
				Preco:             rec.Preco,
				Quantidade:        rec.Quantidade,
				ValorTotal:        valorTotal,
				Peso:              0, // Será calculado depois
				DividendosMensais: 0,
				PesoIdeal:         rec.PesoIdeal,
			}

			carteiraFinal = append(carteiraFinal, fii)
		}
	}

	// Calcula o valor total da carteira final
	valorTotalFinal := valorTotalCarteira + valorTotalRecomendado

	// Calcula os pesos na carteira final
	for i := range carteiraFinal {
		carteiraFinal[i].Peso = (carteiraFinal[i].ValorTotal / valorTotalFinal) * 100
	}

	// Ordena a carteira final por peso (do maior para o menor)
	sort.Slice(carteiraFinal, func(i, j int) bool {
		return carteiraFinal[i].Peso > carteiraFinal[j].Peso
	})

	return carteiraFinal
}

// ObterCarteiraFinalAcao obtém a carteira final de ações após as compras sugeridas
func (s *RecomendadoraService) ObterCarteiraFinalAcao(
	carteira *models.CarteiraAcoes,
	recomendacoes []models.RecomendacaoCompraAcao,
	recomendados []models.AcaoRecomendada,
	valorTotalCarteira, valorTotalRecomendado float64,
) []models.AcaoCarteiraFinal {
	// Mapa para facilitar a busca de ativos na carteira
	ativosNaCarteira := make(map[string]*models.AtivoAcao)
	for i := range carteira.Data {
		ativosNaCarteira[carteira.Data[i].TickerName] = &carteira.Data[i]
	}

	// Mapa para facilitar a busca de informações de recomendações
	infoRecomendados := make(map[string]models.AcaoRecomendada)
	for _, rec := range recomendados {
		infoRecomendados[rec.Ticker] = rec
	}

	// Lista para armazenar as ações da carteira final
	var carteiraFinal []models.AcaoCarteiraFinal

	// Primeiro, adiciona as ações que já estão na carteira
	for _, ativo := range carteira.Data {
		dy, _ := strconv.ParseFloat(strings.TrimSuffix(ativo.DY, "%"), 64)
		pvp, _ := strconv.ParseFloat(strings.Replace(ativo.PVP, ",", ".", -1), 64)
		pl, _ := strconv.ParseFloat(strings.Replace(ativo.PL, ",", ".", -1), 64)

		// Verifica se há compra adicional desta ação
		quantidade := ativo.Quantity
		for _, rec := range recomendacoes {
			if rec.Ticker == ativo.TickerName {
				quantidade += rec.Quantidade
				break
			}
		}

		valorTotal := float64(quantidade) * ativo.CurrentPrice

		// Verifica se tem peso ideal nas recomendações
		var pesoIdeal float64
		if rec, existe := infoRecomendados[ativo.TickerName]; existe {
			pesoIdeal = rec.PesoIdeal
		}

		acao := models.AcaoCarteiraFinal{
			Ticker:     ativo.TickerName,
			Nome:       ativo.TickerName, // Poderia buscar o nome completo em recomendados
			Preco:      ativo.CurrentPrice,
			PL:         pl,
			PVP:        pvp,
			DY:         dy,
			Quantidade: quantidade,
			ValorTotal: valorTotal,
			Peso:       0, // Será calculado depois
			PesoIdeal:  pesoIdeal,
		}

		carteiraFinal = append(carteiraFinal, acao)
	}

	// Adiciona as novas ações que não estavam na carteira original
	for _, rec := range recomendacoes {
		// Verifica se já existe na carteira
		if _, existe := ativosNaCarteira[rec.Ticker]; !existe {
			valorTotal := float64(rec.Quantidade) * rec.Preco

			acao := models.AcaoCarteiraFinal{
				Ticker:     rec.Ticker,
				Nome:       rec.Nome,
				Preco:      rec.Preco,
				PL:         rec.PL,
				PVP:        rec.PVP,
				DY:         rec.DY,
				Quantidade: rec.Quantidade,
				ValorTotal: valorTotal,
				Peso:       0, // Será calculado depois
				PesoIdeal:  rec.PesoIdeal,
			}

			carteiraFinal = append(carteiraFinal, acao)
		}
	}

	// Calcula o valor total da carteira final
	valorTotalFinal := valorTotalCarteira + valorTotalRecomendado

	// Calcula os pesos na carteira final
	for i := range carteiraFinal {
		carteiraFinal[i].Peso = (carteiraFinal[i].ValorTotal / valorTotalFinal) * 100
	}

	// Ordena a carteira final por peso (do maior para o menor)
	sort.Slice(carteiraFinal, func(i, j int) bool {
		return carteiraFinal[i].Peso > carteiraFinal[j].Peso
	})

	return carteiraFinal
}

// ObterCarteiraFinalETF obtém a carteira final de ETFs após as compras sugeridas
func (s *RecomendadoraService) ObterCarteiraFinalETF(
	carteira *models.CarteiraETFs,
	recomendacoes []models.RecomendacaoCompraETF,
	recomendados []models.ETFRecomendado,
	valorTotalCarteira, valorTotalRecomendado float64,
) []models.ETFCarteiraFinal {
	// Mapa para facilitar a busca de ativos na carteira
	ativosNaCarteira := make(map[string]*models.AtivoETF)
	for i := range carteira.Data {
		ativosNaCarteira[carteira.Data[i].TickerName] = &carteira.Data[i]
	}

	// Mapa para facilitar a busca de informações de recomendações
	infoRecomendados := make(map[string]models.ETFRecomendado)
	for _, rec := range recomendados {
		infoRecomendados[rec.Ticker] = rec
	}

	// Lista para armazenar os ETFs da carteira final
	var carteiraFinal []models.ETFCarteiraFinal

	// Primeiro, adiciona os ETFs que já estão na carteira
	for _, ativo := range carteira.Data {
		currentPrice, err := strconv.ParseFloat(ativo.CurrentPrice, 64)
		if err != nil {
			continue
		}

		// Verifica se há compra adicional deste ETF
		quantidade := ativo.Quantity
		for _, rec := range recomendacoes {
			if rec.Ticker == ativo.TickerName {
				quantidade += rec.Quantidade
				break
			}
		}

		valorTotal := float64(quantidade) * currentPrice

		// Verifica se tem peso ideal nas recomendações
		var pesoIdeal float64
		if rec, existe := infoRecomendados[ativo.TickerName]; existe {
			pesoIdeal = rec.PesoIdeal
		}

		etf := models.ETFCarteiraFinal{
			Ticker:     ativo.TickerName,
			Nome:       ativo.TickerName, // Poderia buscar o nome completo em recomendados
			Preco:      currentPrice,
			Quantidade: quantidade,
			ValorTotal: valorTotal,
			Peso:       0, // Será calculado depois
			PesoIdeal:  pesoIdeal,
		}

		carteiraFinal = append(carteiraFinal, etf)
	}

	// Adiciona os novos ETFs que não estavam na carteira original
	for _, rec := range recomendacoes {
		// Verifica se já existe na carteira
		if _, existe := ativosNaCarteira[rec.Ticker]; !existe {
			valorTotal := float64(rec.Quantidade) * rec.Preco

			etf := models.ETFCarteiraFinal{
				Ticker:     rec.Ticker,
				Nome:       rec.Nome,
				Preco:      rec.Preco,
				Quantidade: rec.Quantidade,
				ValorTotal: valorTotal,
				Peso:       0, // Será calculado depois
				PesoIdeal:  rec.PesoIdeal,
			}

			carteiraFinal = append(carteiraFinal, etf)
		}
	}

	// Calcula o valor total da carteira final
	valorTotalFinal := valorTotalCarteira + valorTotalRecomendado

	// Calcula os pesos na carteira final
	for i := range carteiraFinal {
		carteiraFinal[i].Peso = (carteiraFinal[i].ValorTotal / valorTotalFinal) * 100
	}

	// Ordena a carteira final por peso (do maior para o menor)
	sort.Slice(carteiraFinal, func(i, j int) bool {
		return carteiraFinal[i].Peso > carteiraFinal[j].Peso
	})

	return carteiraFinal
}
