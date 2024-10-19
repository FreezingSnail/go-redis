package handler

import (
	"io"
	"kvStore/internal"
	"log"
	"net/http"

	tsLog "kvStore/internal/transactionlog"

	"github.com/gorilla/mux"
)

type HandlerService struct {
	logger tsLog.TransactionLogger
}

func (h *HandlerService) KeyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Handling put key")
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Printf("error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = internal.Put(key, string(value))
	if err != nil {
		log.Printf("error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.WritePut(key, string(value))

	log.Print("done Handling put key")
	w.WriteHeader(http.StatusCreated)
}

func (h *HandlerService) KeyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Handling get key")
	vars := mux.Vars(r)
	key := vars["key"]

	defer r.Body.Close()

	value, err := internal.Get(key)
	if err != nil {
		log.Printf("error %s", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Write([]byte(value))
}

func NewHandlerService() HandlerService {
	logger, err := tsLog.NewFileTransactionLogger("transaction.log")
	if err != nil {
		log.Fatal("failed to create logger: %w", err)
	}
	logger.Run()

	return HandlerService{logger: logger}
}
