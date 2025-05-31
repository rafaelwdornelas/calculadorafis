package handlers

import (
	"bytes"
	"calculadora-investimentos/internal/api"
	"calculadora-investimentos/internal/config"
	"calculadora-investimentos/internal/services"
	"calculadora-investimentos/internal/utils"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// Handlers contém todos os manipuladores HTTP
type Handlers struct {
	Config             *config.Config
	BrapiClient        *api.BrapiClient
	CalculadoraService *services.Calculadora
	DataService        *services.DataService
	DividendoService   *services.DividendoService
}

// NewHandlers cria uma nova instância de Handlers
func NewHandlers() *Handlers {
	// Carregar configuração
	cfg := config.Load()

	// Obter a instância única do cliente BrAPI
	brapiClient := api.GetInstance(
		cfg.APIBaseURL,
		cfg.APIToken,
		cfg.DefaultTimeout,
		cfg.CacheDuracao,
		cfg.CacheLimpeza,
	)

	// Criar serviços
	distribuidoraService := services.NewDistribuidoraService(cfg)
	recomendadoraService := services.NewRecomendadoraService()
	otimizadoraService := services.NewOtimizadoraService()
	dataService := services.NewDataService(cfg, brapiClient)
	dividendoService := services.NewDividendoService()

	// Criar serviço de calculadora com dividendoService
	calculadoraService := services.NewCalculadora(
		distribuidoraService,
		recomendadoraService,
		otimizadoraService,
		dividendoService,
	)

	return &Handlers{
		Config:             cfg,
		BrapiClient:        brapiClient,
		CalculadoraService: calculadoraService,
		DataService:        dataService,
		DividendoService:   dividendoService,
	}
}

// StatusCacheHandler manipula requisições para visualizar o status do cache
func StatusCacheHandler(w http.ResponseWriter, r *http.Request) {
	handlers := NewHandlers()

	// Obter estatísticas do cache
	status := handlers.BrapiClient.StatusCache()

	// Enviar resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// RenderizarTemplate renderiza um template HTML
func (h *Handlers) RenderizarTemplate(w http.ResponseWriter, nomeTemplate string, dados interface{}) error {
	// Criar um novo template com as funções utilitárias
	tmpl, err := template.New(nomeTemplate).Funcs(utils.GetTemplateFuncs()).ParseFiles(h.Config.TemplatesDir + "/" + nomeTemplate)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, dados)
}

// RenderizarTemplateParaString renderiza um template HTML para uma string
func (h *Handlers) RenderizarTemplateParaString(nomeTemplate string, dados interface{}) (template.HTML, error) {
	// Criar um novo template com as funções utilitárias
	tmpl, err := template.New(nomeTemplate).Funcs(utils.GetTemplateFuncs()).ParseFiles(h.Config.TemplatesDir + "/" + nomeTemplate)
	if err != nil {
		log.Printf("Erro ao carregar o template %s: %v", nomeTemplate, err)
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, dados); err != nil {
		log.Printf("Erro ao renderizar o template %s: %v", nomeTemplate, err)
		return "", err
	}

	return template.HTML(buf.String()), nil
}
