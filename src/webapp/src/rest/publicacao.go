package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"webapp/src/config"
	"webapp/src/requisicoes"
	"webapp/src/utils"

	"github.com/gorilla/mux"
)

func PublicacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		CriarPublicacao(w, r)
		return
	}
}

func CurtirPublicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		CurtirPublicacao(w, r)
		return
	}
}

func DescurtirPublicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		DescurtirPublicacao(w, r)
		return
	}
}

func AtualizaPublicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		AtualizarPublicacao(w, r)
		return
	}
}

//Chama a API para cadastrar a publicação no banco de dados
func CriarPublicacao(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	publicacao, err := json.Marshal(map[string]string{
		"titulo":   r.FormValue("titulo"),
		"conteudo": r.FormValue("conteudo"),
	})

	if err != nil {
		utils.JSON(w, http.StatusBadRequest, utils.ErroAPI{Erro: err.Error()})
		return
	}

	url := fmt.Sprintf("%s/publicacao", config.ApiUrl)
	response, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPost, url, bytes.NewBuffer(publicacao))
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, utils.ErroAPI{Erro: err.Error()})
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		utils.TratarStatusCodeErro(w, response)
		return
	}

	utils.JSON(w, response.StatusCode, nil)
}

//TODO: Desenvolver metodo de salvar dados do usuario que curte a publicação na API

//Chama a API para curtir uma publicação
func CurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoID, err := strconv.ParseInt(parametros["publicacaoId"], 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, utils.ErroAPI{Erro: err.Error()})
		return
	}

	url := fmt.Sprintf("%s/publicacoes/%d/curtir", config.ApiUrl, publicacaoID)
	resp, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPost, url, nil)
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

//Chama a API para descurtir uma publicação
func DescurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoID, err := strconv.ParseInt(parametros["publicacaoId"], 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, utils.ErroAPI{Erro: err.Error()})
		return
	}

	url := fmt.Sprintf("%s/publicacoes/%d/descurtir", config.ApiUrl, publicacaoID)
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

//Chama a API para Atualizar uma publicação
func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoID, err := strconv.ParseInt(parametros["publicacaoId"], 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, utils.ErroAPI{Erro: err.Error()})
		return
	}

	r.ParseForm()
	publicacao, err := json.Marshal(map[string]string{
		"titulo":   r.FormValue("titulo"),
		"conteudo": r.FormValue("conteudo"),
	})

	if err != nil {
		utils.JSON(w, http.StatusBadRequest, utils.ErroAPI{Erro: err.Error()})
		return
	}

	url := fmt.Sprintf("%s/publicacoes/%d", config.ApiUrl, publicacaoID)
	resp, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodPut, url, bytes.NewBuffer(publicacao))
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
