package rest

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"webapp/src/utils"
)

func CriarUsuarioHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		CriarUsuario(w, r)
		return
	}
}

//Chama a API para cadastrar um usuario no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	usuario, err := json.Marshal(map[string]string{
		"nome":  r.FormValue("nome"),
		"nick":  r.FormValue("nick"),
		"email": r.FormValue("email"),
		"senha": r.FormValue("senha"),
	})
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, utils.ErroAPI{Erro: err.Error()})
		return
	}

	urlApi := "https://estudos-312813.rj.r.appspot.com/api/usuario/registrar"
	req, err := http.NewRequest(http.MethodPost, urlApi, bytes.NewBuffer(usuario))
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, utils.ErroAPI{Erro: err.Error()})
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		utils.TratarStatusCodeErro(w, resp)
		return
	}

	utils.JSON(w, resp.StatusCode, nil)
}