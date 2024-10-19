package main

import (
	"kvStore/internal"
	"kvStore/internal/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Print("Starting server")

	internal.InitializeTransactionLog()

	kvService := handler.NewHandlerService()

	r := mux.NewRouter()

	r.HandleFunc("/v1/{key}", kvService.KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", kvService.KeyValueGetHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
