package rest

import (
	"net/http"
	"webapp/src/utils"
)

//Renderiza a tela de login
func LoginHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		CarregarTelaLogin(w, r)
		return
	}

}
func CarregarTelaLogin(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "login.html", nil)
}
