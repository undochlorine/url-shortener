package server

import (
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	pb "url-shortener/pb/shortener"
)

func (s *Server) UrlHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /", s.GetUrl)
	handler.HandleFunc("POST /", s.AddUrl)

	return handler
}

// GetUrl
// @Summary Get website
// @Description Get redirected to origin website, using shortened link
// @Tags URL
// @Produce json
// @Param short_url query string true "Short url"
// @Success 200 {object} handbook.CustomersHandbookMsg
// @Failure 500 {string} string "Internal error"
// @Router /url/ [get]
func (s *Server) GetUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	shortUrl := r.URL.Query().Get("short_url")
	if shortUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"shortUrl is required"}`))
		return
	}

	// todo: call svc -> get fullLink, err
	fullUrl, err := s.urlSvc.GetUrl(r.Context(), &pb.ShortUrl{ShortUrl: shortUrl})

	u, err := protojson.Marshal(fullUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"internal server error"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(u)
}

// AddUrl
// @Summary Add URL
// @Description Create alias for the link
// @Tags URL
// @Produce json
// @Param full_url query string true "URL to be shortened"
// @Success 200 {object} handbook.CustomersHandbookMsg
// @Failure 500 {string} string "Internal error"
// @Router /url/ [post]
func (s *Server) AddUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fullUrl := r.URL.Query().Get("full_url")
	if fullUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"fullUrl is required"}`))
		return
	}

	// todo: call svc -> get shortUrl, err
	shortUrl, err := s.urlSvc.Add(r.Context(), &pb.FullUrl{FullUrl: fullUrl})

	u, err := protojson.Marshal(shortUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"internal server error"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(u)
}
