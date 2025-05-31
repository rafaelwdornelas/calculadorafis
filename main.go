package main

import (
	"calculadora-investimentos/internal/config"
	"calculadora-investimentos/internal/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Carregar configurações
	cfg := config.Load()

	// Configurar rotas
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/static/", handlers.StaticHandler)
	mux.HandleFunc("/calcular", handlers.CalcularHandler)
	mux.HandleFunc("/status-cache", handlers.StatusCacheHandler) // Nova rota para verificar o status do cache

	// Iniciar servidor
	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Servidor iniciado em http://localhost%s\n", addr)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
