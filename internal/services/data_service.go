package services

import (
	"bufio"
	"calculadora-investimentos/internal/api"
	"calculadora-investimentos/internal/config"
	"calculadora-investimentos/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// DataService gerencia o acesso a dados da aplicação
type DataService struct {
	Config      *config.Config
	BrapiClient *api.BrapiClient
}

// NewDataService cria um novo serviço de dados
func NewDataService(cfg *config.Config, brapiClient *api.BrapiClient) *DataService {
	return &DataService{
		Config:      cfg,
		BrapiClient: brapiClient,
	}
}

// CarregarRecomendadosFII carrega as recomendações de FIIs do arquivo
func (s *DataService) CarregarRecomendadosFII() ([]models.FIIRecomendado, error) {
	nomeArquivo := filepath.Join(s.Config.DataDir, "recomendados_fiis.txt")
	file, err := os.Open(nomeArquivo)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var recomendados []models.FIIRecomendado
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		linha := scanner.Text()
		campos := strings.Split(linha, "\t")
		if len(campos) == 5 {
			pesoIdeal, _ := strconv.ParseFloat(strings.Replace(strings.TrimSuffix(campos[4], "%"), ",", ".", -1), 64)

			// Obter preço via API
			preco, err := s.BrapiClient.GetQuote(campos[0])
			if err != nil {
				fmt.Printf("Erro ao obter preço para %s: %v\n", campos[0], err)
				continue
			}

			fii := models.FIIRecomendado{
				Ticker:    campos[0],
				Nome:      campos[1],
				Segmento:  campos[2],
				Tipo:      campos[3],
				PesoIdeal: pesoIdeal,
				Preco:     preco,
			}

			recomendados = append(recomendados, fii)
		}
	}

	return recomendados, scanner.Err()
}

// CarregarRecomendadosAcao carrega as recomendações de ações do arquivo
func (s *DataService) CarregarRecomendadosAcao() ([]models.AcaoRecomendada, error) {
	nomeArquivo := filepath.Join(s.Config.DataDir, "recomendados_acoes.txt")
	file, err := os.Open(nomeArquivo)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var recomendados []models.AcaoRecomendada
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		linha := scanner.Text()
		campos := strings.Split(linha, "\t")

		if len(campos) == 3 {
			pesoIdeal, _ := strconv.ParseFloat(strings.Replace(strings.TrimSuffix(campos[2], "%"), ",", ".", -1), 64)

			// Obter preço via API
			preco, err := s.BrapiClient.GetQuote(campos[1])
			if err != nil {
				fmt.Printf("Erro ao obter preço para %s: %v\n", campos[1], err)
				continue
			}

			acao := models.AcaoRecomendada{
				Nome:      campos[0],
				Ticker:    campos[1],
				PesoIdeal: pesoIdeal,
				Preco:     preco,
			}

			recomendados = append(recomendados, acao)
		}
	}

	return recomendados, scanner.Err()
}

// CarregarRecomendadosETF carrega as recomendações de ETFs do arquivo
func (s *DataService) CarregarRecomendadosETF() ([]models.ETFRecomendado, error) {
	nomeArquivo := filepath.Join(s.Config.DataDir, "recomendados_etfs.txt")
	file, err := os.Open(nomeArquivo)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var recomendados []models.ETFRecomendado
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		linha := scanner.Text()
		campos := strings.Split(linha, "\t")

		if len(campos) == 3 {
			pesoIdeal, _ := strconv.ParseFloat(strings.Replace(strings.TrimSuffix(campos[2], "%"), ",", ".", -1), 64)

			// Obter preço via API
			preco, err := s.BrapiClient.GetQuote(campos[0])
			if err != nil {
				fmt.Printf("Erro ao obter preço para %s: %v\n", campos[0], err)
				continue
			}

			etf := models.ETFRecomendado{
				Ticker:    campos[0],
				Nome:      campos[1],
				PesoIdeal: pesoIdeal,
				Preco:     preco,
			}

			recomendados = append(recomendados, etf)
		}
	}

	return recomendados, scanner.Err()
}

// ObterCarteiraAtualFII obtém a carteira atual de FIIs via API
func (s *DataService) ObterCarteiraAtualFII() (*models.CarteiraDados, error) {
	url := fmt.Sprintf("https://investidor10.com.br/api/carteiras/datatable/ativos/%s/Fii?draw=1", s.Config.IDInvestidor10)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var carteira models.CarteiraDados
	err = json.Unmarshal(body, &carteira)
	if err != nil {
		return nil, err
	}

	return &carteira, nil
}

// ObterCarteiraAtualAcao obtém a carteira atual de Ações via API
func (s *DataService) ObterCarteiraAtualAcao() (*models.CarteiraAcoes, error) {
	url := fmt.Sprintf("https://investidor10.com.br/api/carteiras/datatable/ativos/%s/Ticker?draw=1", s.Config.IDInvestidor10)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var carteira models.CarteiraAcoes
	err = json.Unmarshal(body, &carteira)
	if err != nil {
		return nil, err
	}

	return &carteira, nil
}

// ObterCarteiraAtualETF obtém a carteira atual de ETFs via API
func (s *DataService) ObterCarteiraAtualETF() (*models.CarteiraETFs, error) {
	url := fmt.Sprintf("https://investidor10.com.br/api/carteiras/datatable/ativos/%s/Etf?draw=1", s.Config.IDInvestidor10)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var carteira models.CarteiraETFs
	err = json.Unmarshal(body, &carteira)
	if err != nil {
		return nil, err
	}

	return &carteira, nil
}

// ObterCarteiraAtualRendaFixa obtém a carteira atual de Renda Fixa via API
func (s *DataService) ObterCarteiraAtualRendaFixa() (*models.CarteiraRendaFixa, error) {
	url := fmt.Sprintf("https://investidor10.com.br/api/carteiras/datatable/outrosativos/%s/fixed?draw=1", s.Config.IDInvestidor10)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var carteira models.CarteiraRendaFixa
	err = json.Unmarshal(body, &carteira)
	if err != nil {
		return nil, err
	}

	return &carteira, nil
}
