package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"webapp/src/config"
	"webapp/src/cookies"
	"webapp/src/modelos"
	"webapp/src/requisicoes"
	"webapp/src/utils"

	"github.com/gorilla/mux"
)

func LoginHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		CarregarTelaLogin(w, r)
		return
	}

	if r.Method == http.MethodPost {
		FazerLogin(w, r)
		return
	}

}

func CadastroHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		CarregarTelaCadastroUsuario(w, r)
		return
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		CarregarHome(w, r)
		return
	}
}

func PaginaEditPublicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CarregarPagEditPublic(w, r)
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

//Renderiza a pagina princial com as publicações
func CarregarHome(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/publicacoes", config.ApiUrl)
	resp, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, utils.ErroAPI{Erro: err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		utils.TratarStatusCodeErro(w, resp)
		return
	}

	var publicacoes []modelos.Publicacao
	if err = json.NewDecoder(resp.Body).Decode(&publicacoes); err != nil {
		utils.JSON(w, http.StatusUnprocessableEntity, utils.ErroAPI{Erro: err.Error()})
		return
	}

	cookie, _ := cookies.Ler(r)
	usuarioID, _ := strconv.ParseInt(cookie["id"], 10, 64)

	utils.ExecutarTemplate(w, "home.html", struct {
		Publicacoes []modelos.Publicacao
		UsuarioID   int64
	}{
		Publicacoes: publicacoes,
		UsuarioID:   usuarioID,
	})
}

//Renderiza a pagina para edição de uma publicação
func CarregarPagEditPublic(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicID, err := strconv.ParseInt(parametros["publicacaoId"], 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, utils.ErroAPI{Erro: err.Error()})
		return
	}

	url := fmt.Sprintf("%s/publicacao/%d", config.ApiUrl, publicID)
	resp, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, utils.ErroAPI{Erro: err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		utils.TratarStatusCodeErro(w, resp)
		return
	}

	var public modelos.Publicacao
	if err = json.NewDecoder(resp.Body).Decode(&public); err != nil {
		utils.JSON(w, http.StatusUnprocessableEntity, utils.ErroAPI{Erro: err.Error()})
		return
	}

	utils.ExecutarTemplate(w, "atualizar-publicacao.html", public)
}
