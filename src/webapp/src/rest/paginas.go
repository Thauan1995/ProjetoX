package rest

import (
	"net/http"
	"webapp/src/utils"
)

func LoginHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		CarregarTelaLogin(w, r)
		return
	}

}

func CadastroHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		CarregarTelaCadastroUsuario(w, r)
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
