package api

import (
	"calculadora-investimentos/internal/cache"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Instância global do cliente BrAPI (padrão Singleton)
var (
	brapiClientInstance *BrapiClient
	clientOnce          sync.Once
)

// QuoteResponse é a estrutura de resposta da API BrAPI para cotações
type QuoteResponse struct {
	Results []struct {
		Symbol             string  `json:"symbol"`
		ShortName          string  `json:"shortName"`
		RegularMarketPrice float64 `json:"regularMarketPrice"`
		Currency           string  `json:"currency"`
	} `json:"results"`
}

// BrapiClient representa um cliente para a API BrAPI
type BrapiClient struct {
	BaseURL      string
	Token        string
	HTTPClient   *http.Client
	Cache        *cache.Cache
	CacheDuracao time.Duration
}

// GetInstance retorna a instância única do cliente BrAPI
func GetInstance(baseURL, token string, timeout int, cacheDuracao, cacheLimpeza time.Duration) *BrapiClient {
	clientOnce.Do(func() {
		log.Println("Criando nova instância do BrapiClient com cache")
		// Obter a instância do cache
		cacheInstance := cache.GetInstance(cacheDuracao, cacheLimpeza)

		brapiClientInstance = &BrapiClient{
			BaseURL: baseURL,
			Token:   token,
			HTTPClient: &http.Client{
				Timeout: time.Duration(timeout) * time.Second,
			},
			Cache:        cacheInstance,
			CacheDuracao: cacheDuracao,
		}
	})

	status := brapiClientInstance.StatusCache()
	log.Printf("Status atual do cache: %d itens totais, %d itens de cotações",
		status["total_itens"],
		status["itens_por_prefixo"].(map[string]int)["quote_"])

	return brapiClientInstance
}

// GetQuote obtém a cotação de um ativo (com cache)
func (c *BrapiClient) GetQuote(ticker string) (float64, error) {
	// Verificar se existe no cache
	cacheKey := fmt.Sprintf("quote_%s", ticker)
	if cachedPrice, found := c.Cache.Get(cacheKey); found {
		log.Printf("Cache HIT para ticker %s: Usando preço em cache", ticker)
		return cachedPrice.(float64), nil
	}

	log.Printf("Cache MISS para ticker %s: Buscando preço na API", ticker)

	// Se não estiver no cache, fazer a requisição à API
	url := fmt.Sprintf("%s/quote/%s?token=%s&range=1d&interval=1d", c.BaseURL, ticker, c.Token)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API retornou status %d", resp.StatusCode)
	}

	var data QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	if len(data.Results) == 0 {
		return 0, fmt.Errorf("nenhum resultado encontrado para o ticker: %s", ticker)
	}

	price := data.Results[0].RegularMarketPrice

	// Adicionar ao cache
	c.Cache.Set(cacheKey, price, c.CacheDuracao)
	log.Printf("Preço do ticker %s armazenado em cache por %v", ticker, c.CacheDuracao)

	return price, nil
}

// GetQuotesBatch obtém cotações para múltiplos ativos (com cache)
func (c *BrapiClient) GetQuotesBatch(tickers []string) (map[string]float64, error) {
	result := make(map[string]float64)

	for _, ticker := range tickers {
		price, err := c.GetQuote(ticker)
		if err != nil {
			return nil, err
		}
		result[ticker] = price
	}

	return result, nil
}

// StatusCache retorna estatísticas sobre o cache
func (c *BrapiClient) StatusCache() map[string]interface{} {
	return c.Cache.StatusCache()
}

// LimparCache limpa todo o cache de cotações
func (c *BrapiClient) LimparCache() {
	c.Cache.Limpar()
	log.Println("Cache de cotações foi completamente limpo")
}

// InvalidarCache invalida o cache para um ticker específico
func (c *BrapiClient) InvalidarCache(ticker string) {
	cacheKey := fmt.Sprintf("quote_%s", ticker)
	c.Cache.Delete(cacheKey)
	log.Printf("Cache invalidado para ticker %s", ticker)
}

// AtualizarCotacao força a atualização da cotação de um ticker no cache
func (c *BrapiClient) AtualizarCotacao(ticker string) (float64, error) {
	// Invalidar cache existente
	c.InvalidarCache(ticker)

	// Obter nova cotação e armazenar no cache
	return c.GetQuote(ticker)
}

// AtualizarTodasCotacoes força a atualização de todas as cotações no cache
func (c *BrapiClient) AtualizarTodasCotacoes() error {
	// Obter todos os tickers atualmente no cache
	tickers := c.ListarTickersEmCache()

	// Atualizar cada ticker
	for _, ticker := range tickers {
		_, err := c.AtualizarCotacao(ticker)
		if err != nil {
			return fmt.Errorf("erro ao atualizar cotação de %s: %w", ticker, err)
		}
	}

	return nil
}

// ListarTickersEmCache retorna uma lista de todos os tickers armazenados no cache
func (c *BrapiClient) ListarTickersEmCache() []string {
	return c.Cache.ListarChavesComPrefixo("quote_")
}

// PrintStatusCache imprime o status atual do cache no log
func (c *BrapiClient) PrintStatusCache() {
	status := c.StatusCache()
	log.Printf("==== Status do Cache ====")
	log.Printf("Total de itens: %d", status["total_itens"])
	log.Printf("Itens ativos: %d", status["itens_ativos"])
	log.Printf("Itens expirados: %d", status["itens_expirados"])
	log.Printf("Cotações em cache: %d", status["itens_por_prefixo"].(map[string]int)["quote_"])
	log.Printf("========================")
}
