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

	//Config
	r.HandleFunc("/config", rest.ConfigHandler)

	//Estabelecimento
	r.HandleFunc("/estabelecimento", middlewares.Autenticar(rest.EstabelecimentoHandler))

	//Usuario
	r.HandleFunc("/usuario/registrar", rest.RegistraUsuarioHandler)                                   //Registra um usuario
	r.HandleFunc("/usuario/login", rest.LoginHandler)                                                 //Efetua login do usuario
	r.HandleFunc("/usuario/buscar", middlewares.Autenticar(rest.BuscaUsuarioHandler))                 //Busca um usuario
	r.HandleFunc("/usuario/atualizar", middlewares.Autenticar(rest.AtualizaUsuarioHandler))           //Atualiza dados do usuario
	r.HandleFunc("/usuario/{id}/atualizarSenha", middlewares.Autenticar(rest.AtualizaUsuarioHandler)) //Atualiza senha do usuario
	r.HandleFunc("/usuario/deletar", middlewares.Autenticar(rest.DeletaUsuarioHandler))               //Exclui um usuario
	r.HandleFunc("/usuario/seguir", middlewares.Autenticar(rest.SeguidorHandler))                     //Segue um usuario
	r.HandleFunc("/usuario/unseguir", middlewares.Autenticar(rest.SeguidorHandler))                   //Para de seguir um usuario
	r.HandleFunc("/usuario/seguidos", middlewares.Autenticar(rest.BuscaUsuariosSeguidosHandler))      //Busca todos os usuarios que determinado usuario segue
	r.HandleFunc("/usuario/seguidores", middlewares.Autenticar(rest.BuscaSeguidoresHandler))          //Busca todos os usuarios que seguem determinado usuario

	//Publicação
	r.HandleFunc("/publicacao", middlewares.Autenticar(rest.PublicacaoHandler))
	r.HandleFunc("/publicacao/{id}", middlewares.Autenticar(rest.PublicacaoHandler))
	r.HandleFunc("/publicacoes", middlewares.Autenticar(rest.PublicacoesHandler))
	r.HandleFunc("/publicacoes/{idpublic}", middlewares.Autenticar(rest.PublicacaoHandler))
	r.HandleFunc("/publicacoes/{idpublic}/deletar", middlewares.Autenticar(rest.PublicacaoHandler))
	r.HandleFunc("/publicacoes/{idpublic}/curtir", middlewares.Autenticar(rest.PublicacoesHandler))
	r.HandleFunc("/publicacoes/{idpublic}/descurtir", middlewares.Autenticar(rest.PublicacoesHandler))
	r.HandleFunc("/usuario/{usuarioId}/publicacoes", middlewares.Autenticar(rest.PublicacoesUsuarioHandler))

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
