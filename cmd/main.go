package main

import (
	"block9_practice/internal/domain"
	"block9_practice/internal/handler"
	"log"
	"net/http"
)

func main() {
	store := domain.NewStore()
	h := handler.New(store)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", h.List)
	mux.HandleFunc("GET /products/{id}", h.GetByID)
	mux.HandleFunc("POST /products", h.Create)
	mux.HandleFunc("PATCH /products/{id}", h.Update)
	mux.HandleFunc("DELETE /products/{id}", h.Delete)

	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
	log.Println("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
