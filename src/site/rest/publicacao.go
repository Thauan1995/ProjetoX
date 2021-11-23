package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"site/autenticacao"
	"site/publicacao"
	"site/utils"
	"site/utils/log"
	"strconv"

	"github.com/gorilla/mux"
)

func PublicacaoHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodPost {
		CriarPublicacao(w, r)
		return
	}
	if r.Method == http.MethodGet {
		BuscarPublicacao(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
	return
}

//Adiciona uma nova publicação
func CriarPublicacao(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	usuarioID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		log.Warningf(c, "Erro ao extrair usuarioID do token %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao extrair usuarioID do token")
		return
	}

	corpoRequisicao, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warningf(c, "Erro ao receber body da requisição %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao receber body da requisição")
		return
	}

	var public publicacao.Publicacao
	if err = json.Unmarshal(corpoRequisicao, &public); err != nil {
		log.Warningf(c, "Falha ao realizar unmarshal da requisição %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao realizar unmarshal da requisição")
		return
	}

	if novaPublic := publicacao.CriarPublic(c, usuarioID, &public); novaPublic != nil {
		log.Warningf(c, "Erro na criação da publicação %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro na criação da publicação")
		return
	}

	log.Debugf(c, "Publicação criada com sucesso")
	utils.RespondWithJSON(w, http.StatusCreated, public)

}

//Traz as publicações que apareceriam no feed do usuario
func BuscarPublicacoes(w http.ResponseWriter, r *http.Request) {

}

//Traz uma unica publicação pelo nick
func BuscarPublicacao(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Warningf(c, "Falha ao converter id %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao converter id")
		return
	}

	public := publicacao.GetPublicacao(c, id)

	log.Debugf(c, "Busca realizada com sucesso")
	utils.RespondWithJSON(w, http.StatusOK, public)
	return

}

//Altera os dados de uma publicação
func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {

}

//Exclui uma publicação
func DeletarPublicacao(w http.ResponseWriter, r *http.Request) {

}
