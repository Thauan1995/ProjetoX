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

func CarregarPagUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CarregarPaginaUsuarios(w, r)
		return
	}
}

func CarregarPerfilUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CarregarPerfilUsuario(w, r)
		return
	}
}

func CarregarPerfilUsuarioLogadoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CarregarPerfilUsuarioLogado(w, r)
		return
	}
}

func PagEdicaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CarregarPagEdicao(w, r)
		return
	}
	if r.Method == http.MethodPut {
		EditarUsuario(w, r)
		return
	}
}

func PagAttSenhaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		CarregarPagAttSenha(w, r)
		return
	}
	if r.Method == http.MethodPut {
		AtualizarSenha(w, r)
		return
	}
}

//Renderiza a tela de login
func CarregarTelaLogin(w http.ResponseWriter, r *http.Request) {
	cookie, _ := cookies.Ler(r)

	if cookie["token"] != "" {
		http.Redirect(w, r, "/web/home", 302)
		return
	}

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

//Renderiza a pagina de usuarios que atendem o filtro passado
func CarregarPaginaUsuarios(w http.ResponseWriter, r *http.Request) {
	nomeOuNick := r.URL.Query().Get("usuario")
	url := fmt.Sprintf("%s/usuario/buscar?Nick=%s", config.ApiUrl, nomeOuNick)

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

	var usuarios []modelos.Usuario
	if err = json.NewDecoder(resp.Body).Decode(&usuarios); err != nil {
		utils.JSON(w, http.StatusUnprocessableEntity, utils.ErroAPI{Erro: err.Error()})
		return
	}

	utils.ExecutarTemplate(w, "usuarios.html", usuarios)
}

//Renderiza a pagina do perfil do usuario
func CarregarPerfilUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, err := strconv.ParseInt(parametros["idusuario"], 10, 64)
	if err != nil {
		utils.JSON(w, http.StatusBadRequest, utils.ErroAPI{Erro: err.Error()})
		return
	}

	cookie, _ := cookies.Ler(r)
	usuarioLogadoID, _ := strconv.ParseInt(cookie["id"], 10, 64)

	if usuarioID == usuarioLogadoID {
		http.Redirect(w, r, "/web/perfil", 302)
		return
	}

	usuario, err := modelos.BuscarUsuarioCompleto(usuarioID, r)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, utils.ErroAPI{Erro: err.Error()})
		return
	}

	utils.ExecutarTemplate(w, "usuario.html", struct {
		Usuario         modelos.Usuario
		UsuarioLogadoID int64
	}{
		Usuario:         usuario,
		UsuarioLogadoID: usuarioLogadoID,
	})
}

//Renderiza a pagina do perfil do usuario logado
func CarregarPerfilUsuarioLogado(w http.ResponseWriter, r *http.Request) {
	cookie, _ := cookies.Ler(r)
	usuarioID, _ := strconv.ParseInt(cookie["id"], 10, 64)

	usuario, err := modelos.BuscarUsuarioCompleto(usuarioID, r)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, utils.ErroAPI{Erro: err.Error()})
		return
	}

	utils.ExecutarTemplate(w, "perfil.html", usuario)
}

//Renderiza a pagina para edição do usuario
func CarregarPagEdicao(w http.ResponseWriter, r *http.Request) {
	cookie, _ := cookies.Ler(r)
	usuarioID, _ := strconv.ParseInt(cookie["id"], 10, 64)

	canal := make(chan modelos.Usuario)
	go modelos.BuscarDadosUsuario(canal, usuarioID, r)
	usuario := <-canal

	if usuario.ID == 0 {
		utils.JSON(w, http.StatusInternalServerError, utils.ErroAPI{Erro: "Erro ao buscar usuario"})
		return
	}

	utils.ExecutarTemplate(w, "editar-usuario.html", usuario)
}

//Renderiza a pagina para atualização de senha
func CarregarPagAttSenha(w http.ResponseWriter, r *http.Request) {
	utils.ExecutarTemplate(w, "atualizar-senha.html", nil)
}
