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

	//Configuração da pasta assets
	fileServer := http.FileServer(http.Dir("./assets/"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer))

	//Login
	r.HandleFunc("/", rest.LoginHandle)
	r.HandleFunc("/login", rest.LoginHandle)

	//Cadastro
	r.HandleFunc("/criar-usuario", rest.CarregarTelaCadastroUsuario)

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
