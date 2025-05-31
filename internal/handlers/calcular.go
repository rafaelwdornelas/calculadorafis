package handlers

import (
	"calculadora-investimentos/internal/models"
	"calculadora-investimentos/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// CalcularHandler manipula requisições para o endpoint de cálculo
func CalcularHandler(w http.ResponseWriter, r *http.Request) {
	// Definir o tipo de conteúdo como JSON
	w.Header().Set("Content-Type", "application/json")

	// Criar uma nova instância de Handlers
	handlers := NewHandlers()

	// Verificar se é uma requisição POST
	if r.Method != "POST" {
		json.NewEncoder(w).Encode(models.RespostaCalculadora{
			Status:  "error",
			Message: "Método não permitido",
		})
		return
	}

	// Debugar a requisição
	log.Println("Recebida requisição POST para /calcular")
	log.Println("Content-Type:", r.Header.Get("Content-Type"))

	var valorInvestimentoStr string

	// Verificar se é um formulário multipart
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			log.Println("Erro ao processar formulário multipart:", err)
			json.NewEncoder(w).Encode(models.RespostaCalculadora{
				Status:  "error",
				Message: "Erro ao processar o formulário: " + err.Error(),
			})
			return
		}

		// Obter o valor do campo valorInvestimento
		valores := r.MultipartForm.Value["valorInvestimento"]
		if len(valores) > 0 {
			valorInvestimentoStr = valores[0]
			log.Println("Valor obtido do formulário multipart:", valorInvestimentoStr)
		} else {
			log.Println("Campo 'valorInvestimento' não encontrado no formulário multipart")
		}
	} else {
		// Parse formulário normal
		err := r.ParseForm()
		if err != nil {
			log.Println("Erro ao processar formulário padrão:", err)
			json.NewEncoder(w).Encode(models.RespostaCalculadora{
				Status:  "error",
				Message: "Erro ao processar o formulário: " + err.Error(),
			})
			return
		}

		// Obter o valor do campo valorInvestimento
		valorInvestimentoStr = r.FormValue("valorInvestimento")
		log.Println("Valor obtido do formulário padrão:", valorInvestimentoStr)
	}

	// Verificar se os valores estão vazios
	if valorInvestimentoStr == "" {
		log.Println("Valor de investimento está vazio")
		json.NewEncoder(w).Encode(models.RespostaCalculadora{
			Status:  "error",
			Message: "Por favor, digite um valor para investimento inicial",
		})
		return
	}

	// Processar o valor do investimento inicial
	valorInvestimento, err := utils.ProcessarValorMonetario(valorInvestimentoStr)
	if err != nil {
		log.Println("Erro ao converter valor inicial para float:", err)
		json.NewEncoder(w).Encode(models.RespostaCalculadora{
			Status:  "error",
			Message: "Valor de investimento inicial inválido: " + err.Error(),
		})
		return
	}

	log.Println("Valor inicial convertido para float:", valorInvestimento)

	if valorInvestimento <= 0 {
		log.Println("Valor de investimento inicial é menor ou igual a zero")
		json.NewEncoder(w).Encode(models.RespostaCalculadora{
			Status:  "error",
			Message: "O valor de investimento inicial deve ser maior que zero",
		})
		return
	}

	// Verificar se há distribuição personalizada
	distribuicaoPersonalizada := r.FormValue("distribuicaoPersonalizada") == "true"
	var tiposInvestimento models.TiposInvestimento
	tiposInvestimento.FIIs = true
	tiposInvestimento.Acoes = true
	tiposInvestimento.ETFs = true
	tiposInvestimento.RendaFixa = true

	if distribuicaoPersonalizada {
		tiposSelecionadosJSON := r.FormValue("tiposInvestimento")
		if tiposSelecionadosJSON != "" {
			var tiposSelecionados []string
			err := json.Unmarshal([]byte(tiposSelecionadosJSON), &tiposSelecionados)
			if err != nil {
				log.Println("Erro ao processar tipos de investimento selecionados:", err)
			} else {
				// Resetar todos os tipos para false
				tiposInvestimento.FIIs = false
				tiposInvestimento.Acoes = false
				tiposInvestimento.ETFs = false
				tiposInvestimento.RendaFixa = false

				// Marcar apenas os tipos selecionados
				for _, tipo := range tiposSelecionados {
					switch tipo {
					case "FIIs":
						tiposInvestimento.FIIs = true
					case "Ações":
						tiposInvestimento.Acoes = true
					case "ETFs":
						tiposInvestimento.ETFs = true
					case "RendaFixa":
						tiposInvestimento.RendaFixa = true
					}
				}
			}
		}
	}

	// Carregar dados
	recomendadosFII, err := handlers.DataService.CarregarRecomendadosFII()
	if err != nil {
		log.Println("Erro ao carregar recomendações de FIIs:", err)
		json.NewEncoder(w).Encode(models.RespostaCalculadora{
			Status:  "error",
			Message: "Erro ao carregar recomendações de FIIs: " + err.Error(),
		})
		return
	}

	recomendadosAcao, err := handlers.DataService.CarregarRecomendadosAcao()
	if err != nil {
		log.Println("Erro ao carregar recomendações de ações:", err)
		json.NewEncoder(w).Encode(models.RespostaCalculadora{
			Status:  "error",
			Message: "Erro ao carregar recomendações de ações: " + err.Error(),
		})
		return
	}

	recomendadosETF, err := handlers.DataService.CarregarRecomendadosETF()
	if err != nil {
		log.Println("Erro ao carregar recomendações de ETFs:", err)
		// Usar dados padrão mínimos
		recomendadosETF = []models.ETFRecomendado{
			{Ticker: "WRLD11", Nome: "ETF BDRs Mundo", PesoIdeal: 100, Preco: 123.68},
		}
	}

	// Carregar carteiras
	carteiraFII, err := handlers.DataService.ObterCarteiraAtualFII()
	if err != nil {
		log.Println("Erro ao obter carteira atual de FIIs:", err)
		json.NewEncoder(w).Encode(models.RespostaCalculadora{
			Status:  "error",
			Message: "Erro ao obter carteira atual de FIIs: " + err.Error(),
		})
		return
	}

	carteiraAcao, err := handlers.DataService.ObterCarteiraAtualAcao()
	if err != nil {
		log.Println("Erro ao obter carteira atual de ações:", err)
		// Cria uma carteira vazia em caso de erro
		carteiraAcao = &models.CarteiraAcoes{
			Total: 0,
			Data:  []models.AtivoAcao{},
			Draw:  1,
		}
	}

	carteiraETF, err := handlers.DataService.ObterCarteiraAtualETF()
	if err != nil {
		log.Println("Erro ao obter carteira atual de ETFs:", err)
		// Cria uma carteira vazia em caso de erro
		carteiraETF = &models.CarteiraETFs{
			Total: 0,
			Data:  []models.AtivoETF{},
			Draw:  1,
		}
	}

	carteiraRendaFixa, err := handlers.DataService.ObterCarteiraAtualRendaFixa()
	if err != nil {
		log.Println("Erro ao obter carteira atual de renda fixa:", err)
		// Cria uma carteira vazia em caso de erro
		carteiraRendaFixa = &models.CarteiraRendaFixa{
			Total: 0,
			Data:  []models.AtivoRendaFixa{},
			Draw:  1,
		}
	}

	// Calcular recomendações
	dados, err := handlers.CalculadoraService.CalcularRecomendacoes(
		valorInvestimento,
		tiposInvestimento,
		carteiraFII,
		carteiraAcao,
		carteiraETF,
		carteiraRendaFixa,
		recomendadosFII,
		recomendadosAcao,
		recomendadosETF,
	)
	if err != nil {
		log.Println("Erro ao calcular recomendações:", err)
		json.NewEncoder(w).Encode(models.RespostaCalculadora{
			Status:  "error",
			Message: "Erro ao calcular recomendações: " + err.Error(),
		})
		return
	}

	// Renderizar o template
	html, err := handlers.RenderizarTemplateParaString("resultado.html", dados)
	if err != nil {
		log.Println("Erro ao renderizar o resultado:", err)
		json.NewEncoder(w).Encode(models.RespostaCalculadora{
			Status:  "error",
			Message: "Erro ao renderizar o resultado: " + err.Error(),
		})
		return
	}

	// Enviar a resposta JSON com o HTML renderizado
	response := models.RespostaCalculadora{
		Status:    "success",
		Message:   "Cálculo realizado com sucesso",
		DadosHtml: html,
	}

	// Serializar a resposta
	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Println("Erro ao converter resposta para JSON:", err)
		json.NewEncoder(w).Encode(models.RespostaCalculadora{
			Status:  "error",
			Message: "Erro ao converter resposta para JSON: " + err.Error(),
		})
		return
	}

	// Enviar a resposta
	w.Write(jsonResponse)
}
