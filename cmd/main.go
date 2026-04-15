package main

import (
	"block9_practice/config"
	"block9_practice/internal/domain"
	"block9_practice/internal/handler"
	"block9_practice/internal/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	cfg, err := config.Load("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	logger, closeFile, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFile()

	store := domain.NewStore()
	h := handler.New(store, logger)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", h.List)
	mux.HandleFunc("GET /products/{id}", h.GetByID)
	mux.HandleFunc("POST /products", h.Create)
	mux.HandleFunc("PATCH /products/{id}", h.Update)
	mux.HandleFunc("DELETE /products/{id}", h.Delete)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		logger.Info("Сервер запущен")
		http.ListenAndServe(":8080", mux)
	}()

	<-stop

	logger.Info("Сервер выключился")

}
