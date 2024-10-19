package main

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var store = make(map[string]string)

func Put(key string, value string) error {
	store[key] = value
	return nil
}

var ErrorNoSuchKey = errors.New("key not found")

func Get(key string) (string, error) {
	if val, ok := store[key]; ok {
		return val, nil
	}
	return "", ErrorNoSuchKey
}

func Delete(key string) error {
	delete(store, key)
	return nil
}

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Handling put key")
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Print("error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value))
	if err != nil {
		log.Print("error %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Handling get key")
	vars := mux.Vars(r)
	key := vars["key"]

	defer r.Body.Close()

	value, err := Get(key)
	if err != nil {
		log.Print("error %s", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Write([]byte(value))
}

func main() {
	log.Print("Starting server")
	r := mux.NewRouter()

	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
