package main

import (
	"log"
	"net/http"
	"os"
	"site/middlewares"
	"site/rest"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()

	//Usuario
	r.HandleFunc("/usuario/registrar", rest.RegistraUsuarioHandler)
	r.HandleFunc("/usuario/login", rest.LoginHandler)
	r.HandleFunc("/usuario/buscar", middlewares.Autenticar(rest.BuscaUsuarioHandler))
	r.HandleFunc("/usuario/atualizar", middlewares.Autenticar(rest.AtualizaUsuarioHandler))
	r.HandleFunc("/usuario/deletar", middlewares.Autenticar(rest.DeletaUsuarioHandler))
	r.HandleFunc("/usuario/seguir", middlewares.Autenticar(rest.SeguidorHandler))
	r.HandleFunc("/usuario/unseguir", middlewares.Autenticar(rest.SeguidorHandler))
	r.HandleFunc("/usuario/seguidos", middlewares.Autenticar(rest.BuscaUsuariosSeguidosHandler))
	r.HandleFunc("/usuario/seguidores", middlewares.Autenticar(rest.BuscaSeguidoresHandler))

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
