package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"site/estabelecimento"
	"site/utils"
	"site/utils/log"
	"strconv"
)

func EstabelecimentoHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodGet {
		BuscaEstabelecimento(w, r)
		return
	}

	if r.Method == http.MethodPost {
		InsereEstabelecimento(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
	return

}

func BuscaEstabelecimento(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.FormValue("ID") != "" {
		id, err := strconv.ParseInt(r.FormValue("ID"), 10, 64)
		if err != nil {
			log.Warningf(c, "Erro ao converter o ID: %v", err)
			utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao converter ID")
			return
		}

		estab := estabelecimento.GetEstabelecimento(c, id)
		if estab == nil {
			log.Warningf(c, "Estabelecimento não encontrado: %v", id)
			utils.RespondWithError(w, http.StatusBadRequest, 0, "Estabelecimento não encontrado")
		}
	}
	filtros := estabelecimento.Estabelecimento{
		Nome:  r.FormValue("Nome"),
		CNPJ:  r.FormValue("CNPJ"),
		Email: r.FormValue("Email"),
	}

	estabelecimentos, err := estabelecimento.FiltrarEstabelecimento(c, filtros)
	if err != nil {
		log.Warningf(c, "Erro ao buscar Estabelecimento: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao buscar Estabelecimento")
		return
	}
	log.Debugf(c, "Busca realizada com sucesso")
	utils.RespondWithJSON(w, http.StatusOK, estabelecimentos)
}

func InsereEstabelecimento(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	estabelecimentos := &estabelecimento.Estabelecimento{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warningf(c, "Erro ao receber body de Estabelecimento: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao receber body de Estabelecimento")
		return
	}

	err = json.Unmarshal(body, &estabelecimentos)
	if err != nil {
		log.Warningf(c, "Erro ao realizar unmarshal de Estabelecimento: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao realizar Unmarshal")
		return
	}

	err = estabelecimento.InserirEstabelecimento(c, estabelecimentos)
	if err != nil {
		log.Warningf(c, "Falha ao inserir Estabelecimento: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao inserir Estabelecimento")
		return
	}

	log.Debugf(c, "Estabelecimento inserido com sucesso")
	utils.RespondWithJSON(w, http.StatusOK, "Estabelecimento inserido com sucesso")
}
