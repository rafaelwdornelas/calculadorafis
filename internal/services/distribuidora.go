package services

import (
	"calculadora-investimentos/internal/config"
	"calculadora-investimentos/internal/models"
	"math"
)

// DistribuidoraService gerencia a distribuição de ativos na carteira
type DistribuidoraService struct {
	Config *config.Config
}

// NewDistribuidoraService cria um novo serviço de distribuição
func NewDistribuidoraService(cfg *config.Config) *DistribuidoraService {
	return &DistribuidoraService{
		Config: cfg,
	}
}

// CalcularDistribuicaoAtual calcula a distribuição atual da carteira
func (s *DistribuidoraService) CalcularDistribuicaoAtual(
	valorFII, valorAcao, valorETF, valorRendaFixa, valorTotal float64,
) map[string]float64 {
	distribuicao := make(map[string]float64)

	if valorTotal > 0 {
		distribuicao["FIIs"] = (valorFII / valorTotal) * 100
		distribuicao["Ações"] = (valorAcao / valorTotal) * 100
		distribuicao["ETFs"] = (valorETF / valorTotal) * 100
		distribuicao["RendaFixa"] = (valorRendaFixa / valorTotal) * 100
	} else {
		distribuicao["FIIs"] = 0
		distribuicao["Ações"] = 0
		distribuicao["ETFs"] = 0
		distribuicao["RendaFixa"] = 0
	}

	return distribuicao
}

// CalcularDistribuicaoIdeal calcula a distribuição ideal com base nos tipos selecionados
func (s *DistribuidoraService) CalcularDistribuicaoIdeal(tiposInvestimento models.TiposInvestimento) map[string]float64 {
	distribuicaoIdeal := map[string]float64{
		"FIIs":      0.0,
		"Ações":     0.0,
		"ETFs":      0.0,
		"RendaFixa": 0.0,
	}

	// Contar quantos tipos foram selecionados
	tiposSelecionados := 0
	if tiposInvestimento.FIIs {
		tiposSelecionados++
	}
	if tiposInvestimento.Acoes {
		tiposSelecionados++
	}
	if tiposInvestimento.ETFs {
		tiposSelecionados++
	}
	if tiposInvestimento.RendaFixa {
		tiposSelecionados++
	}

	// Distribuir os percentuais proporcionalmente
	if tiposSelecionados > 0 {
		// Usar a distribuição ideal da configuração
		percentuaisPadrao := s.Config.DistribuicaoIdeal

		// Se todos os tipos estão selecionados, usar os percentuais padrão
		if tiposSelecionados == 4 {
			distribuicaoIdeal = percentuaisPadrao
		} else {
			// Calcular o total dos percentuais padrão dos tipos selecionados
			totalSelecionado := 0.0
			if tiposInvestimento.FIIs {
				totalSelecionado += percentuaisPadrao["FIIs"]
			}
			if tiposInvestimento.Acoes {
				totalSelecionado += percentuaisPadrao["Ações"]
			}
			if tiposInvestimento.ETFs {
				totalSelecionado += percentuaisPadrao["ETFs"]
			}
			if tiposInvestimento.RendaFixa {
				totalSelecionado += percentuaisPadrao["RendaFixa"]
			}

			// Distribuir proporcionalmente
			if tiposInvestimento.FIIs {
				distribuicaoIdeal["FIIs"] = (percentuaisPadrao["FIIs"] / totalSelecionado) * 100.0
			}
			if tiposInvestimento.Acoes {
				distribuicaoIdeal["Ações"] = (percentuaisPadrao["Ações"] / totalSelecionado) * 100.0
			}
			if tiposInvestimento.ETFs {
				distribuicaoIdeal["ETFs"] = (percentuaisPadrao["ETFs"] / totalSelecionado) * 100.0
			}
			if tiposInvestimento.RendaFixa {
				distribuicaoIdeal["RendaFixa"] = (percentuaisPadrao["RendaFixa"] / totalSelecionado) * 100.0
			}
		}
	}

	return distribuicaoIdeal
}

// DistribuirInvestimento distribui o valor do investimento entre as classes de ativos
func (s *DistribuidoraService) DistribuirInvestimento(
	valorInvestimento, valorTotalFalta, valorFaltaFII, valorFaltaAcao, valorFaltaETF, valorFaltaRendaFixa float64,
	distribuicaoIdeal map[string]float64,
) (float64, float64, float64, float64) {
	// Distribuir o investimento
	valorParaFII := 0.0
	valorParaAcao := 0.0
	valorParaETF := 0.0
	valorParaRendaFixa := 0.0

	// Estratégia de alocação
	if valorTotalFalta > 0 {
		// Se o investimento for suficiente para cobrir tudo que falta
		if valorTotalFalta <= valorInvestimento {
			valorParaFII = valorFaltaFII
			valorParaAcao = valorFaltaAcao
			valorParaETF = valorFaltaETF
			valorParaRendaFixa = valorFaltaRendaFixa
		} else {
			// Se o investimento não for suficiente, distribuir proporcionalmente
			proporcaoFII := valorFaltaFII / valorTotalFalta
			proporcaoAcao := valorFaltaAcao / valorTotalFalta
			proporcaoETF := valorFaltaETF / valorTotalFalta
			proporcaoRendaFixa := valorFaltaRendaFixa / valorTotalFalta

			valorParaFII = valorInvestimento * proporcaoFII
			valorParaAcao = valorInvestimento * proporcaoAcao
			valorParaETF = valorInvestimento * proporcaoETF
			valorParaRendaFixa = valorInvestimento * proporcaoRendaFixa
		}
	} else {
		// Se não falta nada para nenhuma classe, distribuir conforme percentual ideal
		valorParaFII = valorInvestimento * distribuicaoIdeal["FIIs"] / 100
		valorParaAcao = valorInvestimento * distribuicaoIdeal["Ações"] / 100
		valorParaETF = valorInvestimento * distribuicaoIdeal["ETFs"] / 100
		valorParaRendaFixa = valorInvestimento * distribuicaoIdeal["RendaFixa"] / 100
	}

	// Garantir que valores não sejam negativos
	valorParaFII = math.Max(0, valorParaFII)
	valorParaAcao = math.Max(0, valorParaAcao)
	valorParaETF = math.Max(0, valorParaETF)
	valorParaRendaFixa = math.Max(0, valorParaRendaFixa)

	// Ajustar para garantir que a soma seja exatamente o valor do investimento
	somaValores := valorParaFII + valorParaAcao + valorParaETF + valorParaRendaFixa
	if somaValores > 0 && math.Abs(somaValores-valorInvestimento) > 0.01 {
		fatorAjuste := valorInvestimento / somaValores
		valorParaFII *= fatorAjuste
		valorParaAcao *= fatorAjuste
		valorParaETF *= fatorAjuste
		valorParaRendaFixa *= fatorAjuste
	}

	return valorParaFII, valorParaAcao, valorParaETF, valorParaRendaFixa
}
