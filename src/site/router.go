package main

import (
	"log"
	"net/http"
	"os"
	"site/rest"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()

	//Usuario
	r.HandleFunc("/usuario", rest.UsuarioHandler)
	r.HandleFunc("/usuario/login", rest.LoginHandler)

	http.Handle("/", router)

	var port = os.Getenv("PORT")
	if port == "" {
		port = "5000"
		log.Printf("Padronizando para porta %s", port)
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
