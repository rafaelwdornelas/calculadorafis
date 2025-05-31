package handlers

import (
	"net/http"
)

// IndexHandler manipula requisições para a página inicial
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Criar uma nova instância de Handlers
	handlers := NewHandlers()

	// Renderizar a página inicial
	err := handlers.RenderizarTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Erro ao carregar o template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// StaticHandler manipula requisições para arquivos estáticos
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	// Servir o arquivo estático
	http.ServeFile(w, r, "."+r.URL.Path)
}
