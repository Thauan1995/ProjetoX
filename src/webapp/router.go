package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	//r := router.PathPrefix("/web").Subrouter()

	http.Handle("/", router)

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
		log.Printf("Padronizando para porta %s", port)
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
