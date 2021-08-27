package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"site/usuario"
	"site/utils"
	"site/utils/log"
	"strconv"
)

func UsuarioHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodGet {
		BuscaUsuario(w, r)
		return
	}

	if r.Method == http.MethodPost {
		InsereUsuario(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
}
func BuscaUsuario(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var (
		id  int64
		err error
	)

	if r.FormValue("ID") != "" {
		id, err = strconv.ParseInt(r.FormValue("ID"), 10, 64)
		if err != nil {
			log.Warningf(c, "Erro ao converter ID: %v", err)
			utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao converter ID")
			return
		}
	}
	filtro := usuario.Usuario{
		ID:    id,
		Nome:  r.FormValue("Nome"),
		Nick:  r.FormValue("Nick"),
		Email: r.FormValue("Email"),
	}

	usuario, err := usuario.FiltrarUsuario(c, filtro)
	if err != nil {
		log.Warningf(c, "Erro ao buscar Usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao buscar Usuario")
		return
	}
	log.Debugf(c, "Busca realizada com sucesso")
	utils.RespondWithJSON(w, http.StatusOK, usuario)
}
func InsereUsuario(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	var usuarios []usuario.Usuario

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warningf(c, "Erro ao receber body de Usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao receber body de Usuario")
		return
	}

	err = json.Unmarshal(body, &usuarios)
	if err != nil {
		log.Warningf(c, "Erro ao realizar unmarshal de Usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao realizar unmarshal")
		return
	}

	for i := range usuarios {
		log.Warningf(c, "Inserindo usuario")
		err = usuario.InserirUsuario(c, &usuarios[i])
		if err != nil {
			log.Warningf(c, "Falha ao inserir Usuario: %v", err)
			utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao inserir Usuario")
			return
		}
	}
	log.Debugf(c, "Usuario inserido com sucesso")
	utils.RespondWithJSON(w, http.StatusOK, usuarios)
}
