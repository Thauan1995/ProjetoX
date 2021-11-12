package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"site/autenticacao"
	"site/seguidores"
	"site/usuario"
	"site/utils"
	"site/utils/log"
	"strconv"
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

func SeguidorHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodPost {
		SeguirUsuario(w, r)
		return
	}

	if r.Method == http.MethodPut {
		UnSeguirUsuarios(w, r)
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

	var usu usuario.Usuario

	usuarioIDNoToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		log.Warningf(c, "Erro ao extrair id do usuario da requisição %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao extrair id do usuario da requisição")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warningf(c, "Erro ao receber body de usuario")
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao receber body de usuario")
		return
	}

	if err = json.Unmarshal(body, &usu); err != nil {
		log.Warningf(c, "Erro ao realizar unmarshal do usuario %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao realizar unmarshal do usuario")
		return
	}

	if usu.ID != usuarioIDNoToken {
		log.Warningf(c, "Usuario não tem autenticação para fazer essa ação")
		utils.RespondWithError(w, http.StatusForbidden, 0, "Usuario não tem autenticação para fazer essa ação")
		return
	}

	err = usuario.DeletarUsuario(c, usu)
	if err != nil {
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
	var usu usuario.Usuario

	seguidorID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		log.Warningf(c, "Erro ao extrair token do usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao extrair token do usuario")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warningf(c, "Erro ao receber body de usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao receber body de usuario")
		return
	}

	if err = json.Unmarshal(body, &usu); err != nil {
		log.Warningf(c, "Falha ao fazer unmarshal do usuario a ser seguido: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao fazer unmarshal do usuario a ser seguido")
		return
	}

	if seguidorID == usu.ID {
		log.Warningf(c, "Não é possivel seguir você mesmo")
		utils.RespondWithError(w, http.StatusForbidden, 0, "Não é possivel seguir você mesmo")
		return
	}

	if err = seguidores.Seguir(c, usu.ID, seguidorID); err != nil {
		log.Warningf(c, "Erro seguir usuario %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao seguir usuario")
		return
	}

	log.Debugf(c, "Usuario seguido com sucesso")
	utils.RespondWithJSON(w, http.StatusOK, "Usuario seguido com sucesso")
	return
}

func UnSeguirUsuarios(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	var usu usuario.Usuario

	seguidorID, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		log.Warningf(c, "Erro ao extrair token do usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao extrair token do usuario")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warningf(c, "Erro ao receber body de usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao receber body de usuario")
		return
	}

	if err = json.Unmarshal(body, &usu); err != nil {
		log.Warningf(c, "Falha ao fazer unmarshal do usuario a ser seguido: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Falha ao fazer unmarshal do usuario a ser seguido")
		return
	}

	if seguidorID == usu.ID {
		log.Warningf(c, "Não é possivel parar de seguir você mesmo")
		utils.RespondWithError(w, http.StatusForbidden, 0, "Não é possivel parar de seguir você mesmo")
		return
	}

	if err = seguidores.PararDeSeguir(c, usu.ID, seguidorID); err != nil {
		log.Warningf(c, "Erro ao parar de seguir usuario %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao parar de seguir usuario")
		return
	}

	log.Debugf(c, "Sucesso em parar de seguir usuario")
	utils.RespondWithJSON(w, http.StatusOK, "Sucesso em parar de seguit usuario")
	return
}
