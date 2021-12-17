package rest

import (
	"net/http"
	"webapp/src/cookies"
)

func FazerLogout(w http.ResponseWriter, r *http.Request) {
	cookies.Deletar(w)
	http.Redirect(w, r, "/web/login", 302)
}
