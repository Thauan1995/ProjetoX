package rest

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
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
		log.Fatal(err)
	}

	urlApi := "https://estudos-312813.rj.r.appspot.com/api/usuario/registrar"
	req, err := http.NewRequest(http.MethodPost, urlApi, bytes.NewBuffer(usuario))
	if err != nil {
		log.Fatal(err)
		return
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
}
