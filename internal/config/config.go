package config

import "time"

// Config representa a configuração da aplicação
type Config struct {
	Port              int
	APIToken          string
	APIBaseURL        string
	DataDir           string
	TemplatesDir      string
	StaticDir         string
	DefaultTimeout    int
	IDInvestidor10    string
	DistribuicaoIdeal map[string]float64
	CacheDuracao      time.Duration
	CacheLimpeza      time.Duration
}

// Load carrega a configuração da aplicação
func Load() *Config {
	return &Config{
		Port:           5000,
		APIToken:       "dGubyGPMakfrACS1qoSTye",
		APIBaseURL:     "https://brapi.dev/api",
		DataDir:        "./data",
		TemplatesDir:   "./templates",
		StaticDir:      "./static",
		DefaultTimeout: 10, // segundos
		IDInvestidor10: "1399345",
		DistribuicaoIdeal: map[string]float64{
			"FIIs":      30.0,
			"Ações":     30.0,
			"ETFs":      20.0,
			"RendaFixa": 20.0,
		},
		CacheDuracao: 30 * time.Minute, // Duração do cache (30 minutos)
		CacheLimpeza: 10 * time.Minute, // Intervalo de limpeza (10 minutos)
	}
}
