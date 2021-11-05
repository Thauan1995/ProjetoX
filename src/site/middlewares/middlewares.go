package middlewares

import (
	"net/http"
	"site/autenticacao"
	"site/utils"
	"site/utils/log"
)

// Autenticar verifica se o usuario fazendo a requisição está autenticado
func Autenticar(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		if err := autenticacao.ValidarToken(r); err != nil {
			log.Warningf(c, "Erro ao validar Token: %v", err)
			utils.RespondWithError(w, http.StatusUnauthorized, 0, "Erro ao validar Token")
			return
		}
		proximaFuncao(w, r)
	}
}
