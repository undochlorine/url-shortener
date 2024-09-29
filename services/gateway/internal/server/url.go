package server

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
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

	shortUrl := r.URL.Path[1:]
	if shortUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"shortUrl is required"}`))
		return
	}

	fullUrl, err := s.urlSvc.Get(r.Context(), &pb.ShortUrlMsg{ShortUrl: shortUrl})
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

	var fullUrl = &pb.FullUrlMsg{FullUrl: ""}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"input url is empty"}`))
		return
	}

	err = protojson.Unmarshal(body, fullUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"failed to unmarshall"}`))
	}

	shortUrl, err := s.urlSvc.Set(r.Context(), fullUrl)

	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"internal server error"}`))
		return
	}

	fmt.Printf("sg2")

	u, err := protojson.Marshal(&pb.PairMsg{
		ShortUrl: shortUrl.ShortUrl,
		FullUrl:  fullUrl.FullUrl,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"internal server error"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(u)
}
