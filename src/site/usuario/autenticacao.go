package usuario

import (
	"context"
	"fmt"
	"net/http"
	"site/config"
	"site/utils"
	"site/utils/log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const SessionName = "APPTCSS"

type Claims struct {
	IDUsuario  int64
	Autorizado bool
	jwt.StandardClaims
}

func session(r *http.Request) *Claims {
	var claims = &Claims{}
	var timeNow = utils.GetTimeNow()

	c := r.Context()

	auth := r.Header.Get("Authorization")
	if auth == "" {
		log.Warningf(c, "Autenticação não informada %s", auth)
		return claims
	}

	auth = strings.Replace(auth, "Bearer", "", -1)

	dadosAuth := strings.Split(auth, "$")
	if len(dadosAuth) != 3 {
		log.Warningf(c, "Autenticação inválida: Mal-formada %s", auth)
		return claims
	}

	if sessao := dadosAuth[0]; sessao != SessionName {
		log.Warningf(c, "Autenticação inválida: Sessão incorreta %s / %s", sessao, SessionName)
		return claims
	}

	autenticacao := dadosAuth[1]
	expiracao, err := time.Parse("2006-01-02T15:04:05-0700", dadosAuth[2])
	if err != nil {
		log.Warningf(c, "Autenticação inválida: Expiração inválida %s", err)
		return claims
	}

	if timeNow.After(expiracao) {
		log.Warningf(c, "Autenticação expirou")
		return claims
	}

	chaveSecreta, err := config.GetConfig(c, config.ChaveAutenticacaoAcesso)
	if err != nil {
		log.Warningf(c, "Autenticação inválida: Chave de Autenticação não encontrada %v", err)
		return claims
	}

	tkn, err := jwt.ParseWithClaims(autenticacao, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(chaveSecreta.Value), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Warningf(c, "Autenticação inválida: Assinatura inválida %v", err)
			return claims
		}

		log.Warningf(c, "Autenticação inválida: Token mal-formado %v", err)
		return claims
	}
	if !tkn.Valid {
		log.Warningf(c, "Autenticação Inválida: Token inválido")
		return claims
	}

	claims.Autorizado = tkn.Valid

	return claims
}

func Autenticar(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()

		if r.Header.Get("X-Appengine-Cron") == "true" {
			log.Debugf(c, "Usuario AppEngine-Cron Autenticado")
			f(w, r)
			return
		}

		if r.Header.Get("X-Appengine-Queuename") != "" {
			log.Debugf(c, "Task Queue Autenticado = %v", r.Header.Get("X-Appengine-Queuename"))
			f(w, r)
			return
		}

		if r.FormValue("token") != "" {
			apiToken, err := config.GetConfig(c, config.APIToken)
			if err == nil && r.FormValue("token") == apiToken.Value {
				log.Debugf(c, "Autenticado via Token")
				f(w, r)
				return
			}
		}

		claims := session(r)
		if claims == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, 0, "Autenticação Invalida")
			return
		}

		if !claims.Autorizado {
			log.Warningf(c, "Autenticação Inválida: Token inválido")
			utils.RespondWithError(w, http.StatusUnauthorized, 0, "Autenticação Inválida: Token inválido")
			return
		}

		novaAutorizacao, code := CriarSessao(c, claims.IDUsuario)
		if code != 0 {
			log.Warningf(c, "Autenticação inválida: Atualização da autorização falhou")
			utils.RespondWithError(w, http.StatusUnauthorized, 0, "Autenticação invalida")
			return
		}

		w.Header().Set("Authorization", novaAutorizacao)
		f(w, r)
	}
}

func CriarSessao(c context.Context, idUsuario int64) (string, int) {
	confTempoAcesso := config.GetDefault(c, config.TempoExpiracaoAcesso, "40m")

	timeNow := utils.GetTimeNow()

	tempoAcesso, err := time.ParseDuration(confTempoAcesso.Value)
	if err != nil {
		log.Warningf(c, "Tempo de expiração da autenticação inválida. err: %v", err)
		tempoAcesso, _ = time.ParseDuration("40m")
	}

	expiracoAcesso := timeNow.Add(tempoAcesso)

	credenciais := &Claims{
		IDUsuario: idUsuario,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiracoAcesso.Unix(),
		},
	}

	chaveSecreta, err := config.GetConfig(c, config.ChaveAutenticacaoAcesso)
	if err != nil {
		return "", ErrChaveAutenticacao
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, credenciais)

	tokenString, err := token.SignedString([]byte(chaveSecreta.Value))
	if err != nil {
		return "", ErrAssinarChave
	}

	return fmt.Sprintf("Bearer %s$%s$%s", SessionName, tokenString, expiracoAcesso.Format("2006-01-02T15:04:05-0700")), 0

}

func RemoverAutenticacao(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Authorization", "")
}

func Current(r *http.Request) *Usuario {
	c := r.Context()
	claims := session(r)
	if claims.IDUsuario == 0 || !claims.Autorizado {
		log.Warningf(c, "Erro ao buscar claims")
		log.Warningf(c, "IDUsuario: '%d' - Autorizado: '%v'", claims.IDUsuario, claims.Autorizado)
		return nil
	}

	u := GetUsuario(c, claims.IDUsuario)
	if u == nil {
		log.Warningf(c, "Usuario não encotrado no ID: '%d'", claims.IDUsuario)
		return nil
	}

	if u.Inativo {
		log.Warningf(c, "Usuario invativo")
		return nil
	}
	return u
}
