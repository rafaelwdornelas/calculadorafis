package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// DividendoService gerencia a busca de dividendos dos FIIs
type DividendoService struct {
	HTTPClient *http.Client
}

// NewDividendoService cria um novo serviço de dividendos
func NewDividendoService() *DividendoService {
	return &DividendoService{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FIISearchResult representa o resultado da busca de FII
type FIISearchResult struct {
	Name        string `json:"name"`
	Ticker      string `json:"ticker"`
	ID          int    `json:"id"`
	URL         string `json:"url"`
	TickerID    int    `json:"ticker_id"`
	CompanyName string `json:"company_name"`
}

// DividendoHistorico representa um pagamento de dividendo
type DividendoHistorico struct {
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
}

// ObterTickerID busca o ticker_id do FII
func (s *DividendoService) ObterTickerID(ticker string) (int, error) {
	url := fmt.Sprintf("https://investidor10.com.br/api/fii/searchquery/%s/", ticker)

	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("erro ao buscar ticker_id para %s: %w", ticker, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API retornou status %d para ticker %s", resp.StatusCode, ticker)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("erro ao ler resposta para %s: %w", ticker, err)
	}

	var results []FIISearchResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		return 0, fmt.Errorf("erro ao decodificar JSON para %s: %w", ticker, err)
	}

	if len(results) == 0 {
		return 0, fmt.Errorf("nenhum resultado encontrado para ticker %s", ticker)
	}

	return results[0].TickerID, nil
}

// ObterUltimoDividendo busca o último dividendo pago pelo FII
func (s *DividendoService) ObterUltimoDividendo(tickerID int) (float64, error) {
	url := fmt.Sprintf("https://investidor10.com.br/api/fii/dividendos/chart/%d/360/mes", tickerID)

	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("erro ao buscar dividendos para ticker_id %d: %w", tickerID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API retornou status %d para ticker_id %d", resp.StatusCode, tickerID)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("erro ao ler resposta de dividendos: %w", err)
	}

	var dividendos []DividendoHistorico
	err = json.Unmarshal(body, &dividendos)
	if err != nil {
		return 0, fmt.Errorf("erro ao decodificar JSON de dividendos: %w", err)
	}

	if len(dividendos) == 0 {
		return 0, fmt.Errorf("nenhum dividendo encontrado para ticker_id %d", tickerID)
	}

	// Retornar o último dividendo (último elemento do array)
	ultimoDividendo := dividendos[len(dividendos)-1].Price

	log.Printf("Último dividendo para ticker_id %d: R$ %.2f", tickerID, ultimoDividendo)

	return ultimoDividendo, nil
}

// ObterDividendoFII busca o último dividendo de um FII pelo ticker
func (s *DividendoService) ObterDividendoFII(ticker string) (float64, error) {
	// Primeiro obter o ticker_id
	tickerID, err := s.ObterTickerID(ticker)
	if err != nil {
		log.Printf("Erro ao obter ticker_id para %s: %v", ticker, err)
		return 0, err
	}

	log.Printf("Ticker ID obtido para %s: %d", ticker, tickerID)

	// Depois obter o último dividendo
	dividendo, err := s.ObterUltimoDividendo(tickerID)
	if err != nil {
		log.Printf("Erro ao obter dividendo para %s (id: %d): %v", ticker, tickerID, err)
		return 0, err
	}

	return dividendo, nil
}

// CalcularRendimentosCarteira calcula os rendimentos mensais de todos os FIIs da carteira
func (s *DividendoService) CalcularRendimentosCarteira(fiis []FIIComRendimento) []FIIComRendimento {
	for i := range fiis {
		// Buscar o último dividendo para cada FII
		dividendo, err := s.ObterDividendoFII(fiis[i].Ticker)
		if err != nil {
			log.Printf("Erro ao obter dividendo para %s: %v. Usando valor 0", fiis[i].Ticker, err)
			fiis[i].UltimoDividendo = 0
			fiis[i].RendimentoMensal = 0
		} else {
			fiis[i].UltimoDividendo = dividendo
			fiis[i].RendimentoMensal = dividendo * float64(fiis[i].Quantidade)
			log.Printf("FII %s: %d cotas x R$ %.2f = R$ %.2f de rendimento mensal",
				fiis[i].Ticker, fiis[i].Quantidade, dividendo, fiis[i].RendimentoMensal)
		}

		// Pequena pausa para não sobrecarregar a API
		time.Sleep(100 * time.Millisecond)
	}

	return fiis
}

// FIIComRendimento representa um FII com informações de rendimento
type FIIComRendimento struct {
	Ticker           string
	Nome             string
	Segmento         string
	Tipo             string
	Preco            float64
	Quantidade       int
	ValorTotal       float64
	Peso             float64
	UltimoDividendo  float64
	RendimentoMensal float64
}
