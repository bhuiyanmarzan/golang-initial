package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/bhuiyanmarzan/olx-api/internal/config"
	"github.com/bhuiyanmarzan/olx-api/internal/db"
	"github.com/bhuiyanmarzan/olx-api/internal/handlers"
)

func main() {
	cfg := config.MustLoad()

	db, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	// Global logger
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo, // default
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	fmt.Println("database connected")
	fmt.Println("starting olx server...")

	lh := handlers.NewListingHandler(db, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("server is listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
