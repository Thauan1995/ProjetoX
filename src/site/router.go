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
	r.HandleFunc("/usuario/buscar", rest.BuscaUsuarioHandler)
	r.HandleFunc("/usuario/registrar", rest.RegistraUsuarioHandler)
	r.HandleFunc("/usuario/login", rest.LoginHandler)
	r.HandleFunc("/usuario/atualizar", rest.AtualizaUsuarioHandler)
	r.HandleFunc("/usuario/deletar", rest.DeletaUsuarioHandler)

	//Config
	r.HandleFunc("/config", rest.ConfigHandler)

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
