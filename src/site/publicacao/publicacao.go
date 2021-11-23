package publicacao

import (
	"context"
	"fmt"
	"site/usuario"
	"site/utils"
	"site/utils/consts"
	"site/utils/log"
	"strings"

	"cloud.google.com/go/datastore"
)

const (
	KindPublicacoes = "Publicacoes"
)

type Publicacao struct {
	ID          int64 `datastore:"-"`
	Titulo      string
	Conteudo    string
	AutorID     int64
	AutorNick   string
	Curtidas    int64
	DataCriacao utils.JsonSpecialDateTime
}

//Cria uma publicação
func PutPublicacao(c context.Context, publicacao *Publicacao) error {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return err
	}
	defer datastoreClient.Close()

	key := datastore.IDKey(KindPublicacoes, publicacao.ID, nil)
	key, err = datastoreClient.Put(c, key, publicacao)
	if err != nil {
		log.Warningf(c, "Erro ao atualizar publicação")
		return err
	}
	publicacao.ID = key.ID
	return nil
}

func CriarPublic(c context.Context, usuarioID int64, publicacao *Publicacao) error {
	usuarioBanco := usuario.GetUsuario(c, usuarioID)

	publicacao.AutorID = usuarioBanco.ID
	publicacao.AutorNick = usuarioBanco.Nick

	if publicacao.ID == 0 {
		publicacao.DataCriacao = utils.GetSpecialTimeNow()
	}

	if publicacao.Titulo == "" {
		return fmt.Errorf("O titulo não pode estar em branco")
	}

	if publicacao.Conteudo == "" {
		return fmt.Errorf("O conteudo não pode estar em branco")
	}

	publicacao.Titulo = strings.TrimSpace(publicacao.Titulo)
	publicacao.Conteudo = strings.TrimSpace(publicacao.Conteudo)

	return PutPublicacao(c, publicacao)
}

func GetPublicacao(c context.Context, id int64) *Publicacao {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao concectar-se com o Datastore: %v", err)
		return nil
	}
	defer datastoreClient.Close()

	key := datastore.IDKey(KindPublicacoes, id, nil)

	var publicacao Publicacao
	if err = datastoreClient.Get(c, key, &publicacao); err != nil {
		log.Warningf(c, "Falha ao buscar publicação: %v", err)
		return nil
	}
	publicacao.ID = id
	return &publicacao
}
