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
	c := r.Context()

	if r.Method == http.MethodPost {
		AutenticarUsuario(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
}

func AutenticarUsuario(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	corpoRequisicao, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Warningf(c, "Erro ao receber body para autenticar usuario %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao receber body para autenticar usuario")
		return
	}

	log.Infof(c, "Corpo : %v", string(corpoRequisicao))
	var usuarioLogin usuario.Login
	log.Infof(c, "estrutura: %v", usuarioLogin)

	err = json.Unmarshal(corpoRequisicao, &usuarioLogin)
	if err != nil {
		log.Warningf(c, "Erro ao fazer unmarshal do corpo da requisição de usuario: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao fazer unmarshal do corpo da requisição de usuario")
		return
	}
	log.Infof(c, "fez o Unmarshal")

	usuarioBanco := usuario.GetUsuarioByEmail(c, usuarioLogin.Email)

	err = seguranca.VerifcarSenha(usuarioBanco.Senha, usuarioLogin.Senha)
	if err != nil {
		log.Warningf(c, "Senha inserida no login não compativel com a cadastrada no banco: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, 0, "Senha inserida no login não compativel com a cadastrada no banco")
		return
	}

	w.Write([]byte("Você está logado! Parabens!"))

}
