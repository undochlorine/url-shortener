package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"url-shortener/common/validate"
	"url-shortener/services/gateway/internal/server/middleware"
	"url-shortener/services/gateway/internal/service/url"
)

type Server struct {
	httpSrv *http.Server
	urlSvc  url.Interface
}

func New(urlSvc url.Interface) *Server {
	return &Server{
		urlSvc: urlSvc,
	}
}

func (s *Server) Run(ctx context.Context, host, port string) {
	handler := http.NewServeMux()

	validate.EnableDefaultConfig()

	s.registerRoutes(handler)

	globalChain := middleware.Chain(
		middleware.CORS,
	)

	server := http.Server{
		Addr:         fmt.Sprintf("%v:%v", host, port),
		Handler:      globalChain(handler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.httpSrv = &server

	go func() {
		log.Printf("server started at: %v", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("error listening and serving: %s\n", err)
		}
	}()

	go func(ctx context.Context) {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		s.Shutdown(ctx)
	}(ctx)
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server")
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Server) registerRoutes(handler *http.ServeMux) {
	handler.HandleFunc("GET /probe", s.Probe)

	handler.Handle("/static/", http.StripPrefix("/static", s.staticHandler()))

	handler.Handle("/url/", http.StripPrefix("/url", s.UrlHandler()))
}

func (s *Server) Probe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
