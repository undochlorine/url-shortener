package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"os/signal"
	"url-shortener/services/shortener/config"
	"url-shortener/services/shortener/internal/server"
	"url-shortener/services/shortener/internal/service"
	"url-shortener/services/shortener/internal/storage"
)

func main() {
	cfgPath := flag.String("config", "../config/config.yaml", "config file")
	flag.Parse()

	cfg, err := config.New(*cfgPath)
	if err != nil {
		log.Fatalf("config reading error: %s", err)
	}

	db, err := sqlx.Open("pgx", cfg.Database)
	if err != nil {
		log.Fatalf("opening database connection: %s", err)
	}
	defer db.Close()

	db.SetConnMaxIdleTime(cfg.DBMaxIdleTime)

	storage := storage.New(db)

	svc := service.New(storage)

	builder := server.GrpcServerBuilder{}
	builder.SetService(svc)
	grpcServer, err := builder.Build(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = grpcServer.Start(cfg.Host)
	if err != nil {
		log.Fatal(err)
	}
	defer grpcServer.Cleanup()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	fmt.Println("waiting")
	<-quit
	fmt.Println("shutting down")
}
