package server

import (
	"net/http"
	"os"
)

func (s *Server) staticHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /swagger", s.GetSwagger)

	return handler
}

func (s *Server) GetSwagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	filePath := "./docs/swagger3.0.yaml"

	f, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	w.Write(f)
}
