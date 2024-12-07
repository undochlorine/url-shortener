package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	shortenerPb "url-shortener/pb/shortener"
	"url-shortener/services/gateway/config"
	urlClientPkg "url-shortener/services/gateway/internal/client/url"
	urlCachePkg "url-shortener/services/gateway/internal/client/url/cache"
	"url-shortener/services/gateway/internal/server"
	urlServicePkg "url-shortener/services/gateway/internal/service/url"
)

func main() {
	runServices()
}

func runServices() {
	cfgPath := flag.String("config", "../config/config.yaml", "config file")

	flag.Parse()

	cfg, err := config.New(*cfgPath)
	if err != nil {
		log.Fatalf("config readong error: %s", err)
	}

	urlRpc, urlCloser, err := shortenerPb.NewShortenerService(cfg.Shortener.ConnString)
	if err != nil {
		log.Fatalf("ulr grpc initialization error: %s", err)
	}
	defer urlCloser()

	urlCache := urlCachePkg.New(urlRpc, 10, time.Hour, time.Millisecond)
	urlClient := urlClientPkg.New(urlCache)

	urlService := urlServicePkg.New(urlClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := server.New(urlService)
	srv.Run(ctx, cfg.Server.Host, cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
