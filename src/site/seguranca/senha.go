package seguranca

import (
	"context"
	"fmt"
	"site/usuario"
	"site/utils/log"
)

//Formato para alteração de senha
type Senha struct {
	Nova  string `json:"nova"`
	Atual string `json:"atual"`
}

func BuscarSenha(c context.Context, usuarioID int64) (string, error) {
	usuarioBanco := usuario.GetUsuario(c, usuarioID)
	if usuarioBanco == nil {
		log.Warningf(c, "Erro na busca do usuario no banco")
		return "", fmt.Errorf("Erro na busca do usuario no banco")
	}

	return usuarioBanco.Senha, nil
}

func AtualizarSenha(c context.Context, usuarioID int64, senha string) error {
	usuarioBanco := usuario.GetUsuario(c, usuarioID)
	if usuarioBanco == nil {
		log.Warningf(c, "Erro na busca do usuario no banco")
		return fmt.Errorf("Erro na busca do usuario no banco")
	}

	usuarioBanco.Senha = senha

	return usuario.PutUsuario(c, usuarioBanco)
}
