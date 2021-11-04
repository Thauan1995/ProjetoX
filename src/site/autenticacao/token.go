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

	secretKey, err := config.GetConfig(c, config.ChaveAutenticacaoAcesso)
	if err != nil {
		log.Warningf(c, "Erro ao buscar Chave de Autenticação Acesso")
		return "", fmt.Errorf("Erro ao buscar Chave de Autenticação Acesso")
	}

	permissoes := jwt.MapClaims{}
	permissoes["authorized"] = true
	permissoes["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissoes["usuarioId"] = usuarioID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes)
	return token.SignedString([]byte(secretKey.Value))
}

// ValidarToken verifica se o token passado na requisição é valido
func ValidarToken(r *http.Request) error {
	c := r.Context()
	tokenString := extrairToken(r)
	log.Infof(c, "tokenString %v", tokenString)
	token, err := jwt.Parse(tokenString, retornaChaveVerificacao)
	if err != nil {
		log.Warningf(c, "Erro ao fazer o Parse do token: %v", err)
		return err
	}
	log.Infof(c, "token: %v", token)
	return nil
}

func extrairToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

func retornaChaveVerificacao(token *jwt.Token) (interface{}, error) {
	var c context.Context
	secretKey, err := config.GetConfig(c, config.ChaveAutenticacaoAcesso)
	if err != nil {
		return "", fmt.Errorf("Erro ao buscar Chave de Autenticação Acesso")
	}
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Método de assinatura inesperado! %v", token.Header["alg"])
	}
	return secretKey, nil
}
