package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"
)

// DataComService gerencia a análise de datas com (ex-dividendo)
type DataComService struct {
	HTTPClient *http.Client
}

// NewDataComService cria um novo serviço de análise de data com
func NewDataComService() *DataComService {
	return &DataComService{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Dividendo representa um pagamento de dividendo histórico
type Dividendo struct {
	Tipo          string
	DataCom       string
	DataPagamento string
	Valor         string
	DataParsed    time.Time
}

// AnaliseDataCom contém a análise de quando comprar um ativo
type AnaliseDataCom struct {
	ProximaDataCom    time.Time
	DiasAteDataCom    int
	StatusCompra      string // "SEGURO", "ALERTA", "EVITAR", "NAO_COMPRAR"
	MensagemStatus    string
	UltimosDividendos []Dividendo
	PadraoMensal      bool
	DiaPagamentoComum int
}

// AnalisarDataComTicker analisa a data com de um ticker específico
func (s *DataComService) AnalisarDataComTicker(ticker string, tipoAtivo string) (*AnaliseDataCom, error) {
	// Buscar histórico de dividendos
	dividendos, err := s.extrairDividendos(ticker, tipoAtivo)
	if err != nil {
		return nil, fmt.Errorf("erro ao extrair dividendos: %w", err)
	}

	if len(dividendos) == 0 {
		return &AnaliseDataCom{
			StatusCompra:   "SEGURO",
			MensagemStatus: "Sem histórico de dividendos disponível",
		}, nil
	}

	// Parsear datas
	for i := range dividendos {
		if data, err := time.Parse("02/01/2006", dividendos[i].DataCom); err == nil {
			dividendos[i].DataParsed = data
		}
	}

	// Ordenar por data (mais recente primeiro)
	sort.Slice(dividendos, func(i, j int) bool {
		return dividendos[i].DataParsed.After(dividendos[j].DataParsed)
	})

	hoje := time.Now()
	analise := &AnaliseDataCom{
		UltimosDividendos: dividendos[:min(5, len(dividendos))],
	}

	// Análise diferente para FIIs (pagamento mensal) e Ações (trimestral/semestral)
	if tipoAtivo == "FII" {
		analise.PadraoMensal = true
		s.analisarPadraoMensal(dividendos, analise, hoje)
	} else {
		analise.PadraoMensal = false
		s.analisarPadraoTrimestral(dividendos, analise, hoje)
	}

	return analise, nil
}

// extrairDividendos busca o histórico de dividendos do ativo
func (s *DataComService) extrairDividendos(ticker, tipoAtivo string) ([]Dividendo, error) {
	var url string

	// Construir URL baseado no tipo
	if tipoAtivo == "FII" {
		url = fmt.Sprintf("https://investidor10.com.br/fiis/%s/", strings.ToLower(ticker))
	} else if tipoAtivo == "ETF" {
		// ETFs geralmente não pagam dividendos mensais
		return []Dividendo{}, nil
	} else {
		// Ações
		url = fmt.Sprintf("https://investidor10.com.br/acoes/%s/", strings.ToLower(ticker))
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	content := string(body)
	var dividendos []Dividendo

	// Buscar dividendos com padrão flexível
	blockPattern := `(?s)(Dividendos|JCP|Rendimento)[^0-9]{0,50}(\d{2}/\d{2}/\d{4})[^0-9]{0,50}(\d{2}/\d{2}/\d{4})[^0-9]{0,200}(0,\d+)`
	re := regexp.MustCompile(blockPattern)
	matches := re.FindAllStringSubmatch(content, -1)

	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) >= 5 {
			key := match[2] + match[3]
			if !seen[key] {
				seen[key] = true
				div := Dividendo{
					Tipo:          match[1],
					DataCom:       match[2],
					DataPagamento: match[3],
					Valor:         match[4],
				}
				dividendos = append(dividendos, div)
			}
		}
	}

	return dividendos, nil
}

// analisarPadraoMensal analisa FIIs que pagam mensalmente
func (s *DataComService) analisarPadraoMensal(dividendos []Dividendo, analise *AnaliseDataCom, hoje time.Time) {
	// Analisar últimos 12 meses
	umAnoAtras := hoje.AddDate(-1, 0, 0)
	diasDataCom := make(map[int]int)
	var dividendosUltimoAno []Dividendo

	for _, div := range dividendos {
		if div.DataParsed.After(umAnoAtras) {
			dividendosUltimoAno = append(dividendosUltimoAno, div)
			dia := div.DataParsed.Day()
			diasDataCom[dia]++
		}
	}

	// Encontrar dia mais comum
	var diaComum int
	maxFreq := 0
	for dia, freq := range diasDataCom {
		if freq > maxFreq {
			diaComum = dia
			maxFreq = freq
		}
	}

	analise.DiaPagamentoComum = diaComum

	// Calcular próxima data com
	proximaDataCom := s.calcularProximaDataComMensal(dividendos[0].DataParsed, diaComum, hoje)
	proximaDataComUtil := s.ajustarParaDiaUtil(proximaDataCom)

	analise.ProximaDataCom = proximaDataComUtil
	analise.DiasAteDataCom = int(proximaDataComUtil.Sub(hoje).Hours() / 24)

	// Definir status de compra
	s.definirStatusCompra(analise)
}

// analisarPadraoTrimestral analisa ações que pagam trimestral/semestralmente
func (s *DataComService) analisarPadraoTrimestral(dividendos []Dividendo, analise *AnaliseDataCom, hoje time.Time) {
	// Mapear padrão de pagamento por mês
	padraoMes := make(map[time.Month][]int)

	for _, div := range dividendos {
		mes := div.DataParsed.Month()
		dia := div.DataParsed.Day()
		padraoMes[mes] = append(padraoMes[mes], dia)
	}

	// Encontrar próximo mês provável de pagamento
	var proximaDataCom *time.Time

	// Verificar próximos 6 meses
	for i := 0; i < 6; i++ {
		mesTestado := hoje.AddDate(0, i, 0).Month()
		anoTestado := hoje.AddDate(0, i, 0).Year()

		if dias, existe := padraoMes[mesTestado]; existe && len(dias) >= 2 {
			// Calcular dia médio
			soma := 0
			for _, d := range dias {
				soma += d
			}
			diaMedia := soma / len(dias)

			dataTestada := time.Date(anoTestado, mesTestado, diaMedia, 0, 0, 0, 0, hoje.Location())
			if dataTestada.After(hoje) {
				proximaDataCom = &dataTestada
				break
			}
		}
	}

	if proximaDataCom != nil {
		proximaDataComUtil := s.ajustarParaDiaUtil(*proximaDataCom)
		analise.ProximaDataCom = proximaDataComUtil
		analise.DiasAteDataCom = int(proximaDataComUtil.Sub(hoje).Hours() / 24)
	} else {
		// Se não encontrou padrão, usar última data + 3 meses
		proximaDataCom := dividendos[0].DataParsed.AddDate(0, 3, 0)
		proximaDataComUtil := s.ajustarParaDiaUtil(proximaDataCom)
		analise.ProximaDataCom = proximaDataComUtil
		analise.DiasAteDataCom = int(proximaDataComUtil.Sub(hoje).Hours() / 24)
	}

	s.definirStatusCompra(analise)
}

// calcularProximaDataComMensal calcula a próxima data com para pagamento mensal
func (s *DataComService) calcularProximaDataComMensal(ultimaData time.Time, diaComum int, hoje time.Time) time.Time {
	// Começar com o próximo mês da última data
	proximaData := ultimaData.AddDate(0, 1, 0)
	proximaData = time.Date(
		proximaData.Year(),
		proximaData.Month(),
		diaComum,
		0, 0, 0, 0,
		proximaData.Location(),
	)

	// Continuar adicionando meses até encontrar data futura
	for proximaData.Before(hoje) || proximaData.Equal(hoje) {
		proximaData = proximaData.AddDate(0, 1, 0)
	}

	return proximaData
}

// ajustarParaDiaUtil ajusta a data para dia útil
func (s *DataComService) ajustarParaDiaUtil(data time.Time) time.Time {
	// Se cai no sábado, voltar para sexta
	if data.Weekday() == time.Saturday {
		return data.AddDate(0, 0, -1)
	}
	// Se cai no domingo, voltar para sexta
	if data.Weekday() == time.Sunday {
		return data.AddDate(0, 0, -2)
	}
	return data
}

// definirStatusCompra define o status e mensagem baseado nos dias até a data com
func (s *DataComService) definirStatusCompra(analise *AnaliseDataCom) {
	dias := analise.DiasAteDataCom

	if dias < 0 {
		analise.StatusCompra = "SEGURO"
		analise.MensagemStatus = "✅ Data com já passou - Compra liberada"
	} else if dias == 0 {
		analise.StatusCompra = "NAO_COMPRAR"
		analise.MensagemStatus = "🚫 HOJE É A DATA COM - NÃO COMPRAR!"
	} else if dias <= 2 {
		analise.StatusCompra = "EVITAR"
		analise.MensagemStatus = fmt.Sprintf("🚫 EVITE COMPRAR - Data com em %d dias (%s)",
			dias, analise.ProximaDataCom.Format("02/01/2006"))
	} else if dias <= 5 {
		analise.StatusCompra = "ALERTA"
		analise.MensagemStatus = fmt.Sprintf("⚠️ ALERTA - Data com em %d dias (%s). Considere aguardar",
			dias, analise.ProximaDataCom.Format("02/01/2006"))
	} else {
		analise.StatusCompra = "SEGURO"
		analise.MensagemStatus = fmt.Sprintf("✅ SEGURO - %d dias até a data com (%s)",
			dias, analise.ProximaDataCom.Format("02/01/2006"))
	}
}

// min retorna o menor entre dois inteiros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
