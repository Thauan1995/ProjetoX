package rest

import (
	"fmt"
	"net/http"
	"webapp/src/config"
	"webapp/src/requisicoes"
	"webapp/src/utils"
)

func LoginHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		CarregarTelaLogin(w, r)
		return
	}

	if r.Method == http.MethodPost {
		FazerLogin(w, r)
		return
	}

}

func CadastroHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		CarregarTelaCadastroUsuario(w, r)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		CarregarHome(w, r)
		return
	}
}

//Renderiza a tela de login
func CarregarTelaLogin(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "login.html", nil)
}

//Renderiza a tela de cadastro de usuario
func CarregarTelaCadastroUsuario(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "cadastro.html", nil)
}

//Renderiza a pagina princial com as publicações
func CarregarHome(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/publicacoes", config.ApiUrl)
	resp, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	fmt.Println(resp.StatusCode, err)

	utils.ExecutarTemplate(w, "home.html", nil)
}
