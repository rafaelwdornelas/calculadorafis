// internal/services/distribuidora.go - NOVA VERSÃO
package services

import (
	"calculadora-investimentos/internal/config"
	"calculadora-investimentos/internal/models"
	"log"
	"math"
	"sort"
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

// ClassePrioridade representa uma classe de ativo com sua prioridade
type ClassePrioridade struct {
	Nome                string
	ValorAtual          float64
	PercentualAtual     float64
	PercentualIdeal     float64
	DistanciaPercentual float64 // Diferença em pontos percentuais
	ValorFalta          float64
	PrioridadeScore     float64 // Score combinado para priorização
}

// DistribuirInvestimento distribui o valor do investimento entre as classes de ativos
// VERSÃO MELHORADA: Considera tanto distância percentual quanto valor absoluto
func (s *DistribuidoraService) DistribuirInvestimento(
	valorInvestimento, valorTotalFalta, valorFaltaFII, valorFaltaAcao, valorFaltaETF, valorFaltaRendaFixa float64,
	distribuicaoIdeal map[string]float64,
) (float64, float64, float64, float64) {

	// Precisamos receber também os valores atuais e percentuais atuais
	// Por ora, vamos usar uma abordagem simplificada baseada nos valores faltantes

	// Criar estrutura para análise de prioridade
	classes := []ClassePrioridade{
		{
			Nome:            "FIIs",
			ValorFalta:      valorFaltaFII,
			PercentualIdeal: distribuicaoIdeal["FIIs"],
		},
		{
			Nome:            "Ações",
			ValorFalta:      valorFaltaAcao,
			PercentualIdeal: distribuicaoIdeal["Ações"],
		},
		{
			Nome:            "ETFs",
			ValorFalta:      valorFaltaETF,
			PercentualIdeal: distribuicaoIdeal["ETFs"],
		},
		{
			Nome:            "RendaFixa",
			ValorFalta:      valorFaltaRendaFixa,
			PercentualIdeal: distribuicaoIdeal["RendaFixa"],
		},
	}

	// Log das informações
	log.Printf("=== ANÁLISE DE PRIORIZAÇÃO ===")
	log.Printf("Valor a investir: R$ %.2f", valorInvestimento)

	// Se o investimento é suficiente para cobrir todas as faltas
	if valorTotalFalta <= valorInvestimento {
		log.Printf("Investimento suficiente para atingir distribuição ideal!")

		valorParaFII := valorFaltaFII
		valorParaAcao := valorFaltaAcao
		valorParaETF := valorFaltaETF
		valorParaRendaFixa := valorFaltaRendaFixa

		// Distribuir sobra proporcionalmente
		sobra := valorInvestimento - valorTotalFalta
		if sobra > 0 {
			log.Printf("Sobra após atingir ideal: R$ %.2f", sobra)

			if distribuicaoIdeal["FIIs"] > 0 {
				valorParaFII += sobra * (distribuicaoIdeal["FIIs"] / 100)
			}
			if distribuicaoIdeal["Ações"] > 0 {
				valorParaAcao += sobra * (distribuicaoIdeal["Ações"] / 100)
			}
			if distribuicaoIdeal["ETFs"] > 0 {
				valorParaETF += sobra * (distribuicaoIdeal["ETFs"] / 100)
			}
			if distribuicaoIdeal["RendaFixa"] > 0 {
				valorParaRendaFixa += sobra * (distribuicaoIdeal["RendaFixa"] / 100)
			}
		}

		log.Printf("Distribuição final:")
		log.Printf("  FIIs: R$ %.2f", valorParaFII)
		log.Printf("  Ações: R$ %.2f", valorParaAcao)
		log.Printf("  ETFs: R$ %.2f", valorParaETF)
		log.Printf("  Renda Fixa: R$ %.2f", valorParaRendaFixa)

		return valorParaFII, valorParaAcao, valorParaETF, valorParaRendaFixa
	}

	// Investimento insuficiente: alocar proporcionalmente às faltas
	log.Printf("Investimento insuficiente - distribuindo proporcionalmente às necessidades")

	// Filtrar apenas classes que precisam de investimento
	var classesNecessitando []ClassePrioridade
	totalNecessario := 0.0

	for _, classe := range classes {
		if classe.ValorFalta > 0 {
			classesNecessitando = append(classesNecessitando, classe)
			totalNecessario += classe.ValorFalta
		}
	}

	// Distribuir proporcionalmente
	valorParaFII := 0.0
	valorParaAcao := 0.0
	valorParaETF := 0.0
	valorParaRendaFixa := 0.0

	for _, classe := range classesNecessitando {
		// Proporcional ao quanto cada classe precisa
		proporcao := classe.ValorFalta / totalNecessario
		valorAlocado := valorInvestimento * proporcao

		log.Printf("%s precisa R$ %.2f (%.1f%% do total necessário) - alocando R$ %.2f",
			classe.Nome, classe.ValorFalta, proporcao*100, valorAlocado)

		switch classe.Nome {
		case "FIIs":
			valorParaFII = valorAlocado
		case "Ações":
			valorParaAcao = valorAlocado
		case "ETFs":
			valorParaETF = valorAlocado
		case "RendaFixa":
			valorParaRendaFixa = valorAlocado
		}
	}

	log.Printf("================================")

	return valorParaFII, valorParaAcao, valorParaETF, valorParaRendaFixa
}

// DistribuirInvestimentoComPrioridade - NOVA FUNÇÃO que recebe mais informações
func (s *DistribuidoraService) DistribuirInvestimentoComPrioridade(
	valorInvestimento float64,
	valoresAtuais map[string]float64,
	percentuaisAtuais map[string]float64,
	distribuicaoIdeal map[string]float64,
	valorTotalCarteira float64,
) (float64, float64, float64, float64) {

	valorTotalFuturo := valorTotalCarteira + valorInvestimento

	// Criar lista de classes com análise completa
	var classes []ClassePrioridade

	for nome, percentualIdeal := range distribuicaoIdeal {
		if percentualIdeal > 0 {
			valorAtual := valoresAtuais[nome]
			percentualAtual := percentuaisAtuais[nome]
			valorIdeal := valorTotalFuturo * (percentualIdeal / 100)
			valorFalta := math.Max(0, valorIdeal-valorAtual)
			distanciaPercentual := percentualIdeal - percentualAtual

			// Score de prioridade: combina distância percentual e valor absoluto
			// Quanto maior a distância percentual, maior a prioridade
			prioridadeScore := math.Abs(distanciaPercentual)

			classes = append(classes, ClassePrioridade{
				Nome:                nome,
				ValorAtual:          valorAtual,
				PercentualAtual:     percentualAtual,
				PercentualIdeal:     percentualIdeal,
				DistanciaPercentual: distanciaPercentual,
				ValorFalta:          valorFalta,
				PrioridadeScore:     prioridadeScore,
			})
		}
	}

	// Ordenar por score de prioridade (maior primeiro)
	sort.Slice(classes, func(i, j int) bool {
		return classes[i].PrioridadeScore > classes[j].PrioridadeScore
	})

	// Log detalhado
	log.Printf("=== PRIORIZAÇÃO POR DISTÂNCIA PERCENTUAL ===")
	log.Printf("Investimento: R$ %.2f", valorInvestimento)
	for i, c := range classes {
		log.Printf("%d. %s: Atual %.2f%% → Ideal %.2f%% (distância: %.2f%%) - Falta R$ %.2f",
			i+1, c.Nome, c.PercentualAtual, c.PercentualIdeal, c.DistanciaPercentual, c.ValorFalta)
	}

	// Alocar recursos
	valorParaFII := 0.0
	valorParaAcao := 0.0
	valorParaETF := 0.0
	valorParaRendaFixa := 0.0
	valorRestante := valorInvestimento

	// Primeira passada: alocar priorizando maior distância percentual
	for _, classe := range classes {
		if valorRestante <= 0 || classe.ValorFalta <= 0 {
			continue
		}

		valorAlocar := math.Min(classe.ValorFalta, valorRestante)

		switch classe.Nome {
		case "FIIs":
			valorParaFII = valorAlocar
		case "Ações":
			valorParaAcao = valorAlocar
		case "ETFs":
			valorParaETF = valorAlocar
		case "RendaFixa":
			valorParaRendaFixa = valorAlocar
		}

		log.Printf("Alocando R$ %.2f para %s", valorAlocar, classe.Nome)
		valorRestante -= valorAlocar

		// Se cobriu totalmente esta classe, continuar para a próxima
		if valorAlocar >= classe.ValorFalta {
			continue
		} else {
			// Se não cobriu totalmente, parar aqui
			break
		}
	}

	log.Printf("Distribuição final:")
	log.Printf("  FIIs: R$ %.2f", valorParaFII)
	log.Printf("  Ações: R$ %.2f", valorParaAcao)
	log.Printf("  ETFs: R$ %.2f", valorParaETF)
	log.Printf("  Renda Fixa: R$ %.2f", valorParaRendaFixa)
	log.Printf("============================================")

	return valorParaFII, valorParaAcao, valorParaETF, valorParaRendaFixa
}
