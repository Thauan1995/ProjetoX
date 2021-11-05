package autenticacao

import (
	"context"
	"fmt"
	"net/http"
	"site/config"
	"site/utils/log"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Criar token retorna um token assinado com as permissões do usuario
func CriarToken(c context.Context, usuarioID int64) (string, error) {

	permissoes := jwt.MapClaims{}
	permissoes["authorized"] = true
	permissoes["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissoes["usuarioId"] = usuarioID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes)
	return token.SignedString([]byte(config.SecretKey))
}

// ValidarToken verifica se o token passado na requisição é valido
func ValidarToken(r *http.Request) error {
	c := r.Context()
	tokenString := extrairToken(r)

	token, err := jwt.Parse(tokenString, retornaChaveVerificacao)
	if err != nil {
		log.Warningf(c, "Erro ao fazer o Parse do token: %v", err)
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	log.Warningf(c, "Token inválido")
	return fmt.Errorf("Token inválido")
}

func extrairToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

func retornaChaveVerificacao(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Método de assinatura inesperado! %v", token.Header["alg"])
	}
	return config.SecretKey, nil
}
