package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"site/autenticacao"
	"site/seguidores"
	"site/seguranca"
	"site/usuario"
	"site/utils"
	"site/utils/log"
	"strconv"

	"github.com/gorilla/mux"
)

func BuscaUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodGet {
		BuscaUsuario(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
	return
}

func RegistraUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodPost {
		InsereUsuario(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
	return
}

func AtualizaUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodPut {
		AtualizaUsuario(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
	return
}

func AtualizaSenhaHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodPut {
		AtualizaSenha(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
	return
}

func DeletaUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodDelete {
		DeletaUsuario(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
	return
}

func SeguirHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodPut {
		SeguirUsuario(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
	return
}

func UnFollowHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodPut {
		UnFollowUsuarios(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
	return
}

func BuscaUsuariosSeguidosHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodGet {
		BuscaUsuariosSeguidos(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
	return
}

func BuscaSeguidoresHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodGet {
		BuscaSeguidores(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
	return
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

	var usuarios usuario.Usuario

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

	checkBanco := usuario.GetUsuarioByEmail(c, usuarios)
	log.Debugf(c, "checkBanco: %v", checkBanco)

	if checkBanco {
		log.Warningf(c, "Email ou nick ja existe")
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Email ou nick ja existe")
		return
	}
	err = usuario.InserirUsuario(c, &usuarios)
	if err != nil {
		log.Warningf(c, "Falha ao inserir Usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao inserir Usuario")
		return
	}

	log.Debugf(c, "Usuario inserido com sucesso")
	utils.RespondWithJSON(w, http.StatusOK, "Usuario Inserido com sucesso")
}

func AtualizaUsuario(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	var usu usuario.Usuario

	usuarioIDNoToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		log.Warningf(c, "Erro ao extrair token do usuario da requisição %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao extrair token do usuario da requisição")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warningf(c, "Erro ao receber body de usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao receber body de usuario")
		return
	}

	err = json.Unmarshal(body, &usu)
	if err != nil {
		log.Warningf(c, "Falha ao fazer unmarshal de usuario %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao fazer unmarshal de usuario")
		return
	}

	if usu.ID != usuarioIDNoToken {
		log.Warningf(c, "Usuario não tem autorizaçao para fazer essa ação")
		utils.RespondWithError(w, http.StatusForbidden, 0, "Usuario não tem autorização para fazer essa ação")
		return
	}

	err = usuario.AtualizarUsuario(c, usu)
	if err != nil {
		log.Warningf(c, "Erro ao atualizar usuario %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao atualizar usuario")
		return
	}

	log.Debugf(c, "Usuario atualizado com sucesso")
	utils.RespondWithJSON(w, http.StatusOK, "Usuario atualizado com sucesso")
}

func DeletaUsuario(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	params := mux.Vars(r)
	idUsu, err := strconv.ParseInt(params["idusuario"], 10, 64)
	if err != nil {
		log.Warningf(c, "Falha ao converter id do usuário: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao converter id do usuário")
	}

	usuarioID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		log.Warningf(c, "Erro ao extrair id do usuário da requisição %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao extrair id do usuário da requisição")
		return
	}

	usu := usuario.GetUsuario(c, idUsu)

	if usu.ID != usuarioID {
		log.Warningf(c, "Usuario %v não tem autenticação para fazer essa ação", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Usuario não tem autenticação para fazer essa ação")
		return
	}

	if err = usuario.DeletarUsuario(c, *usu); err != nil {
		log.Warningf(c, "Falha ao deletar usuario")
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao deletar usuario")
		return
	}

	log.Warningf(c, "Usuario deletado")
	utils.RespondWithJSON(w, http.StatusOK, "Usuario deletado")
}

//Permite que um usuario siga outro
func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	params := mux.Vars(r)
	idSeguidor, err := strconv.ParseInt(params["idusuario"], 10, 64)
	if err != nil {
		log.Warningf(c, "Falha ao converter id do usuário: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao converter id do usuário")
	}
	seguidorID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		log.Warningf(c, "Erro ao extrair token do usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao extrair token do usuario")
		return
	}

	seg := usuario.GetUsuario(c, idSeguidor)

	if seg.ID == seguidorID {
		log.Warningf(c, "Não é possivel seguir você mesmo")
		utils.RespondWithError(w, http.StatusForbidden, 0, "Não é possivel seguir você mesmo")
		return
	}

	if err = seguidores.Seguir(c, seg.ID, seguidorID); err != nil {
		log.Warningf(c, "Erro seguir usuario %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao seguir usuario")
		return
	}

	log.Debugf(c, "Usuario seguido com sucesso")
	utils.RespondWithJSON(w, http.StatusOK, "Usuario seguido com sucesso")
	return
}

//Permite que um usuario pare de seguir outro
func UnFollowUsuarios(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	params := mux.Vars(r)
	idSeguidor, err := strconv.ParseInt(params["idusuario"], 10, 64)
	if err != nil {
		log.Warningf(c, "Falha ao converter id do usuário: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao converter id do usuário")
	}
	seguidorID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		log.Warningf(c, "Erro ao extrair token do usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao extrair token do usuario")
		return
	}

	seg := usuario.GetUsuario(c, idSeguidor)

	if seg.ID == seguidorID {
		log.Warningf(c, "Não é possivel parar de seguir você mesmo")
		utils.RespondWithError(w, http.StatusForbidden, 0, "Não é possivel parar de seguir você mesmo")
		return
	}

	if err = seguidores.PararDeSeguir(c, seg.ID, seguidorID); err != nil {
		log.Warningf(c, "Erro ao parar de seguir usuario %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao parar de seguir usuario")
		return
	}

	log.Debugf(c, "Sucesso em deixar de seguir usuario")
	utils.RespondWithJSON(w, http.StatusOK, "Sucesso em deixar de seguir usuario")
	return
}

// Traz todas as pessoas que o usuario esta seguindo
func BuscaUsuariosSeguidos(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	params := mux.Vars(r)
	idSeguidor, err := strconv.ParseInt(params["idusuario"], 10, 64)
	if err != nil {
		log.Warningf(c, "Falha ao converter id do usuário: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao converter id do usuário")
	}

	usuarios, err := seguidores.BuscarUsuariosSeguidos(c, idSeguidor)
	if err != nil {
		log.Warningf(c, "Erro ao efetuar a busca de usuarios %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao efetuar a busca de usuarios")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, usuarios)
	return

}

// Traz todas as pessoas que seguem o usuario
func BuscaSeguidores(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	params := mux.Vars(r)
	idUsu, err := strconv.ParseInt(params["idusuario"], 10, 64)
	if err != nil {
		log.Warningf(c, "Falha ao converter id do usuário: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao converter id do usuário")
	}

	filtro := seguidores.Seguidor{}

	seguidors, err := seguidores.FiltrarSeguidores(c, filtro)
	if err != nil {
		log.Warningf(c, "Erro ao filtrar seguidores %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao filtrar seguidor")
		return
	}

	usuarios, err := seguidores.BuscarSeguidores(c, seguidors, idUsu)
	if err != nil {
		log.Warningf(c, "Erro ao efetuar busca de seguidores %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao efetuar busca de seguidores")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, usuarios)
	return
}

// Permite alterar a senha do usuario
func AtualizaSenha(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	usuarioID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		log.Warningf(c, "Erro ao extrair ID do usuario da requisição %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao extrair ID do usuario da requisição")
		return
	}

	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		log.Warningf(c, "Falha ao converter id %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao converter id")
		return
	}

	if usuarioID != id {
		log.Warningf(c, "Não é possivel atualizar a senha de um usuario que não seja o seu")
		utils.RespondWithError(w, http.StatusForbidden, 0, "Não é possivel atualizar a senha de um usuario que não seja o seu")
		return
	}

	corpoRequisicao, err := ioutil.ReadAll(r.Body)

	var senha seguranca.Senha
	if err = json.Unmarshal(corpoRequisicao, &senha); err != nil {
		log.Warningf(c, "Falha ao realizar unmarshal da requisição para alterar senha %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao realizar unmarshal da requisição para alterar senha")
		return
	}

	senhaBanco, err := seguranca.BuscarSenha(c, usuarioID)
	if err != nil {
		log.Warningf(c, "Falha ao buscar senha do usuario no banco %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao buscar senha do usuario no banco")
		return
	}

	if err = seguranca.VerifcarSenha(senhaBanco, senha.Atual); err != nil {
		log.Warningf(c, "A senha atual não condiz com a que está no banco %v", err)
		utils.RespondWithError(w, http.StatusUnauthorized, 0, "A senha atual não condiz com a que está no banco")
		return
	}

	senhaNovaHash, err := seguranca.Hash(senha.Nova)
	if err != nil {
		log.Warningf(c, "Falha ao criptografar nova senha %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ai criptografar nova senha")
		return
	}

	if err = seguranca.AtualizarSenha(c, usuarioID, string(senhaNovaHash)); err != nil {
		log.Warningf(c, "Erro ao atualizar senha %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao atualizar senha")
		return
	}

	log.Debugf(c, "Senha atualizada com sucesso")
	utils.RespondWithJSON(w, http.StatusOK, "Senha atualizada com sucesso")
	return
}
