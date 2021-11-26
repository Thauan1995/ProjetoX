package rest

import (
	"net/http"
)

// Utiliza o email e senha do usuario para autenticar na aplicação
func FazerLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	/* 	usuario, err := json.Marshal(map[string]string{
	   		"email": r.FormValue("email"),
	   		"senha": r.FormValue("senha"),
	   	})

	   	if err != nil {
	   		utils.JSON(w, http.StatusBadRequest, utils.ErroAPI{Erro: err.Error()})
	   		return
	   	}
	*/

}
