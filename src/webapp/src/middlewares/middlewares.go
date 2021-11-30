package middlewares

import (
	"log"
	"net/http"
	"webapp/src/cookies"
)

// Escreve informações da requisição no terminal
func Logger(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\n %s %s %s", r.Method, r.RequestURI, r.Host)
		proximaFuncao(w, r)
	}
}

// Verifica se existe cookies de dados de autenticação no browser do usuario
func Autenticar(proximaFuncao http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := cookies.Ler(r); err != nil {
			http.Redirect(w, r, "/web/login", 302)
			return
		}
		proximaFuncao(w, r)
	}
}
