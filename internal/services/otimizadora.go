package services

import (
	"calculadora-investimentos/internal/models"
	"sort"
)

// OtimizadoraService otimiza o investimento das sobras
type OtimizadoraService struct{}

// NewOtimizadoraService cria um novo serviço de otimização
func NewOtimizadoraService() *OtimizadoraService {
	return &OtimizadoraService{}
}

// AtivoCandidate representa um ativo candidato a receber investimento adicional
type AtivoCandidate struct {
	Tipo      string  // "FII", "ACAO" ou "ETF"
	Indice    int     // Índice no array original
	Preco     float64 // Preço unitário
	Ticker    string  // Código do ativo
	PesoIdeal float64 // Peso ideal para desempate
}

// OtimizarSobras otimiza o investimento das sobras
func (s *OtimizadoraService) OtimizarSobras(
	valorSobra float64,
	recomendadosFII []models.FIIRecomendado,
	recomendadosAcao []models.AcaoRecomendada,
	recomendadosETF []models.ETFRecomendado,
	recomendacoesFII *[]models.RecomendacaoCompraFII,
	recomendacoesAcao *[]models.RecomendacaoCompraAcao,
	recomendacoesETF *[]models.RecomendacaoCompraETF,
	valorTotalRecomendadoFII *float64,
	valorTotalRecomendadoAcao *float64,
	valorTotalRecomendadoETF *float64,
) float64 {
	// Cria uma estrutura para armazenar os ativos e seus preços
	var candidatos []AtivoCandidate

	// Adiciona FIIs como candidatos (se não estiver já em excesso)
	if len(*recomendacoesFII) > 0 {
		for i, rec := range recomendadosFII {
			candidatos = append(candidatos, AtivoCandidate{
				Tipo:      "FII",
				Indice:    i,
				Preco:     rec.Preco,
				Ticker:    rec.Ticker,
				PesoIdeal: rec.PesoIdeal,
			})
		}
	}

	// Adiciona Ações como candidatos
	for i, rec := range recomendadosAcao {
		candidatos = append(candidatos, AtivoCandidate{
			Tipo:      "ACAO",
			Indice:    i,
			Preco:     rec.Preco,
			Ticker:    rec.Ticker,
			PesoIdeal: rec.PesoIdeal,
		})
	}

	// Adiciona ETFs como candidatos
	for i, rec := range recomendadosETF {
		candidatos = append(candidatos, AtivoCandidate{
			Tipo:      "ETF",
			Indice:    i,
			Preco:     rec.Preco,
			Ticker:    rec.Ticker,
			PesoIdeal: rec.PesoIdeal,
		})
	}

	// Ordena os candidatos por preço (do mais barato para o mais caro)
	sort.Slice(candidatos, func(i, j int) bool {
		return candidatos[i].Preco < candidatos[j].Preco
	})

	// Enquanto houver sobra e candidatos disponíveis
	sobraFinal := valorSobra
	for sobraFinal > 0 && len(candidatos) > 0 {
		// Encontra o candidato mais barato que cabe na sobra
		candidatoEncontrado := false

		for i, candidato := range candidatos {
			if candidato.Preco <= sobraFinal {
				// Podemos comprar este ativo
				switch candidato.Tipo {
				case "FII":
					// Verifica se o ticker já existe nas recomendações
					encontrado := false
					for j, rec := range *recomendacoesFII {
						if rec.Ticker == candidato.Ticker {
							// Aumenta a quantidade
							(*recomendacoesFII)[j].Quantidade++
							(*recomendacoesFII)[j].ValorCompra += candidato.Preco
							*valorTotalRecomendadoFII += candidato.Preco
							encontrado = true
							break
						}
					}

					// Se não encontrou, cria uma nova recomendação
					if !encontrado && len(*recomendacoesFII) > 0 {
						fii := recomendadosFII[candidato.Indice]
						novaRec := models.RecomendacaoCompraFII{
							Ticker:         fii.Ticker,
							Nome:           fii.Nome,
							Segmento:       fii.Segmento,
							Tipo:           fii.Tipo,
							Preco:          fii.Preco,
							PesoAtual:      0,
							PesoIdeal:      fii.PesoIdeal,
							Diferenca:      fii.PesoIdeal,
							Quantidade:     1,
							ValorCompra:    fii.Preco,
							PesoAposCompra: 0, // Será recalculado depois
						}
						*recomendacoesFII = append(*recomendacoesFII, novaRec)
						*valorTotalRecomendadoFII += fii.Preco
					}

				case "ACAO":
					// Verifica se o ticker já existe nas recomendações
					encontrado := false
					for j, rec := range *recomendacoesAcao {
						if rec.Ticker == candidato.Ticker {
							// Aumenta a quantidade
							(*recomendacoesAcao)[j].Quantidade++
							(*recomendacoesAcao)[j].ValorCompra += candidato.Preco
							*valorTotalRecomendadoAcao += candidato.Preco
							encontrado = true
							break
						}
					}

					// Se não encontrou, cria uma nova recomendação
					if !encontrado {
						acao := recomendadosAcao[candidato.Indice]
						novaRec := models.RecomendacaoCompraAcao{
							Ticker:         acao.Ticker,
							Nome:           acao.Nome,
							Preco:          acao.Preco,
							PesoAtual:      0,
							PesoIdeal:      acao.PesoIdeal,
							Diferenca:      acao.PesoIdeal,
							Quantidade:     1,
							ValorCompra:    acao.Preco,
							PesoAposCompra: 0, // Será recalculado depois
						}
						*recomendacoesAcao = append(*recomendacoesAcao, novaRec)
						*valorTotalRecomendadoAcao += acao.Preco
					}

				case "ETF":
					// Verifica se o ticker já existe nas recomendações
					encontrado := false
					for j, rec := range *recomendacoesETF {
						if rec.Ticker == candidato.Ticker {
							// Aumenta a quantidade
							(*recomendacoesETF)[j].Quantidade++
							(*recomendacoesETF)[j].ValorCompra += candidato.Preco
							*valorTotalRecomendadoETF += candidato.Preco
							encontrado = true
							break
						}
					}

					// Se não encontrou, cria uma nova recomendação
					if !encontrado {
						etf := recomendadosETF[candidato.Indice]
						novaRec := models.RecomendacaoCompraETF{
							Ticker:         etf.Ticker,
							Nome:           etf.Nome,
							Preco:          etf.Preco,
							PesoAtual:      0,
							PesoIdeal:      etf.PesoIdeal,
							Diferenca:      etf.PesoIdeal,
							Quantidade:     1,
							ValorCompra:    etf.Preco,
							PesoAposCompra: 0, // Será recalculado depois
						}
						*recomendacoesETF = append(*recomendacoesETF, novaRec)
						*valorTotalRecomendadoETF += etf.Preco
					}
				}

				// Atualiza a sobra
				sobraFinal -= candidato.Preco
				candidatoEncontrado = true

				// Remove este candidato da lista e recomeça
				candidatos = append(candidatos[:i], candidatos[i+1:]...)
				break
			}
		}

		// Se não encontrou nenhum candidato que caiba na sobra, encerra o loop
		if !candidatoEncontrado {
			break
		}
	}

	// Retorna a sobra que não foi possível investir
	return sobraFinal
}
