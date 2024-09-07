package server

import (
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	pb "url-shortener/pb/shortener"
)

func (s *Server) UrlHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /", s.GetUrl)
	handler.HandleFunc("POST /", s.SetUrl)

	return handler
}

func (s *Server) GetUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	shortUrl := r.URL.Query().Get("short_url")
	if shortUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"shortUrl is required"}`))
		return
	}

	// todo: call svc -> get fullLink, err
	fullUrl, err := s.urlSvc.Get(r.Context(), &pb.ShortUrl{ShortUrl: shortUrl})

	u, err := protojson.Marshal(fullUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"internal server error"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(u)
}

func (s *Server) SetUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fullUrl := r.URL.Query().Get("full_url")
	if fullUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"fullUrl is required"}`))
		return
	}

	// todo: call svc -> get shortUrl, err
	shortUrl, err := s.urlSvc.Set(r.Context(), &pb.FullUrl{FullUrl: fullUrl})

	u, err := protojson.Marshal(shortUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"internal server error"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(u)
}
