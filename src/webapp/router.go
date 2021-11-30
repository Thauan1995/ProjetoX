package main

import (
	"fmt"
	"log"
	"net/http"
	"webapp/src/config"
	"webapp/src/cookies"
	"webapp/src/middlewares"
	"webapp/src/rest"
	"webapp/src/utils"

	"github.com/gorilla/mux"
)

func main() {
	config.Carregar()
	cookies.Configurar()
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
	r.HandleFunc("/usuario/registrar", rest.CriarUsuarioHandler)

	//Home
	r.HandleFunc("/home", middlewares.Logger(middlewares.Autenticar(rest.HomeHandler)))

	http.Handle("/", router)

	fmt.Printf("Escutando na porta %d\n", config.Porta)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), nil); err != nil {
		log.Fatal(err)
	}
}
