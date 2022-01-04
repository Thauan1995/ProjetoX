package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"webapp/src/config"
	"webapp/src/cookies"
	"webapp/src/requisicoes"
	"webapp/src/utils"

	"github.com/gorilla/mux"
)

func CriarUsuarioHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		CriarUsuario(w, r)
		return
	}
}

func PararDeSeguirHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		PararDeSeguir(w, r)
		return
	}
}

func SeguirHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		Seguir(w, r)
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

	url := fmt.Sprintf("%s/usuario/registrar", config.ApiUrl)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(usuario))
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

//Chama a API para parar de seguir um usuario
func PararDeSeguir(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseInt(parametros["idusuario"], 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, utils.ErroAPI{Erro: err.Error()})
		return
	}

	url := fmt.Sprintf("%s/usuario/unfollow/%d", config.ApiUrl, usuarioID)
	resp, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPut, url, nil)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, utils.ErroAPI{Erro: err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		utils.TratarStatusCodeErro(w, resp)
		return
	}

	utils.JSON(w, resp.StatusCode, nil)
}

//Chama a API para seguir um usuario
func Seguir(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseInt(parametros["idusuario"], 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, utils.ErroAPI{Erro: err.Error()})
		return
	}

	url := fmt.Sprintf("%s/usuario/seguir/%d", config.ApiUrl, usuarioID)
	resp, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPut, url, nil)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, utils.ErroAPI{Erro: err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		utils.TratarStatusCodeErro(w, resp)
		return
	}

	utils.JSON(w, resp.StatusCode, nil)
}

//Chama a API para editar o usuario
func EditarUsuario(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	usuario, err := json.Marshal(map[string]string{
		"nome":  r.FormValue("nome"),
		"email": r.FormValue("email"),
		"nick":  r.FormValue("nick"),
	})

	if err != nil {
		utils.JSON(w, http.StatusBadRequest, utils.ErroAPI{Erro: err.Error()})
		return
	}

	cookie, _ := cookies.Ler(r)
	usuarioID, _ := strconv.ParseInt(cookie["id"], 10, 64)

	url := fmt.Sprintf("%s/usuario/atualizar/%d", config.ApiUrl, usuarioID)
	resp, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPut, url, bytes.NewBuffer(usuario))
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, utils.ErroAPI{Erro: err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		utils.TratarStatusCodeErro(w, resp)
		return
	}

	utils.JSON(w, resp.StatusCode, nil)
}
