package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"site/autenticacao"
	"site/seguranca"
	"site/usuario"
	"site/utils"
	"site/utils/log"
)

// Login é responsavel por autenticar um usuario na API
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	if r.Method == http.MethodPost {
		AcessarUsuario(w, r)
		return
	}

	log.Warningf(c, "Método não permitido")
	utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Método não permitido")
}

func AcessarUsuario(w http.ResponseWriter, r *http.Request) {
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

	usuarioBanco, err := usuario.FiltrarUsuario(c, usuarioLogin)

	for _, usu := range usuarioBanco {
		err = seguranca.VerifcarSenha(usu.Senha, usuarioLogin.Senha)
		if err != nil {
			log.Warningf(c, "Senha inserida no login não compativel com a cadastrada no banco: %v", err)
			utils.RespondWithError(w, http.StatusBadRequest, 0, "Senha inserida no login não compativel com a cadastrada no banco")
			return
		}
		token, _ := autenticacao.CriarToken(usu.ID)
		log.Infof(c, "Token gerado para autenticação: %v", token)
	}

	w.Write([]byte("Você está logado! Parabens!"))

}
