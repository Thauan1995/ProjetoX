package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"site/seguranca"
	"site/usuario"
	"site/utils"
	"site/utils/log"
)

// Login é responsavel por autenticar um usuario na API
func LoginHandler(w http.ResponseWriter, r *http.Request) {

}

func AutenticarUsuario(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	corpoRequisicao, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warningf(c, "Erro ao receber body para autenticar usuario %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao receber body para autenticar usuario")
		return
	}

	var usuarioLogin usuario.Usuario
	err = json.Unmarshal(corpoRequisicao, &usuarioLogin)
	if err != nil {
		log.Warningf(c, "Erro ao fazer unmarshal do corpo da requisição de usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao fazer unmarshal do corpo da requisição de usuario")
		return
	}
	usuarioBanco := usuario.GetUsuario(c, usuarioLogin.ID)

	err = seguranca.VerifcarSenha(usuarioBanco.Senha, usuarioLogin.Senha)
	if err != nil {
		log.Warningf(c, "Senha inserida no login não compativel com a cadastrada no banco: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Senha inserida no login não compativel com a cadastrada no banco")
		return
	}
	// parado na Aula 86
	// voltar para aula 85 e verificar se há necessidade ne implementar metodos de valida() e prepara()
}
