package main

import (
	"log"
	"net/http"
	"os"
	"webapp/src/rest"
	"webapp/src/utils"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	utils.CarregarTemplates()

	r := router.PathPrefix("/web").Subrouter()

	//Login
	r.HandleFunc("/", rest.LoginHandle)
	r.HandleFunc("/login", rest.LoginHandle)

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
