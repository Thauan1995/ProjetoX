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

	//Logout
	r.HandleFunc("/logout", middlewares.Logger(middlewares.Autenticar(rest.FazerLogout)))

	//Cadastro
	r.HandleFunc("/criar-usuario", rest.CarregarTelaCadastroUsuario)
	r.HandleFunc("/usuario/registrar", rest.CriarUsuarioHandler)

	//Home
	r.HandleFunc("/home", middlewares.Logger(middlewares.Autenticar(rest.HomeHandler)))

	//Publicacoes
	r.HandleFunc("/publicacoes", middlewares.Logger(middlewares.Autenticar(rest.PublicacaoHandler)))
	r.HandleFunc("/publicacoes/{publicacaoId}/curtir", middlewares.Logger(middlewares.Autenticar(rest.CurtirPublicHandler)))
	r.HandleFunc("/publicacoes/{publicacaoId}/descurtir", middlewares.Logger(middlewares.Autenticar(rest.DescurtirPublicHandler)))
	r.HandleFunc("/publicacoes/{publicacaoId}/editar", middlewares.Logger(middlewares.Autenticar(rest.PaginaEditPublicHandler)))
	r.HandleFunc("/publicacoes/{publicacaoId}", middlewares.Logger(middlewares.Autenticar(rest.AtualizaPublicHandler)))
	r.HandleFunc("/publicacoes/{publicacaoId}/deletar", middlewares.Logger(middlewares.Autenticar(rest.ExcluiPublicHandler)))

	http.Handle("/", router)

	fmt.Printf("Escutando na porta %d\n", config.Porta)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), nil); err != nil {
		log.Fatal(err)
	}
}
