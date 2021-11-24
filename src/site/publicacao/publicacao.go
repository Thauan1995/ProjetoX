package publicacao

import (
	"context"
	"fmt"
	"site/seguidores"
	"site/usuario"
	"site/utils"
	"site/utils/consts"
	"site/utils/log"
	"sort"
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

func GetMultPublicacao(c context.Context, keys []*datastore.Key) ([]Publicacao, error) {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return []Publicacao{}, err
	}
	defer datastoreClient.Close()

	publicacao := make([]Publicacao, len(keys))
	if err := datastoreClient.GetMulti(c, keys, publicacao); err != nil {
		if errs, ok := err.(datastore.MultiError); ok {
			for _, e := range errs {
				if e == datastore.ErrNoSuchEntity {
					return []Publicacao{}, nil
				}
			}
		}
		log.Warningf(c, "Erro ao buscar Multi Usuarios: %v", err)
		return []Publicacao{}, err
	}
	for i := range keys {
		publicacao[i].ID = keys[i].ID
	}
	return publicacao, nil
}

func FiltrarPublicacoes(c context.Context, publicacao Publicacao) ([]Publicacao, error) {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return nil, err
	}
	defer datastoreClient.Close()

	q := datastore.NewQuery(KindPublicacoes)

	if publicacao.AutorNick != "" {
		q = q.Filter("AutorNick =", publicacao.AutorNick)
	}

	if publicacao.AutorID != 0 {
		q = q.Filter("AutorID =", publicacao.AutorID)
	}

	if publicacao.ID != 0 {
		key := datastore.IDKey(KindPublicacoes, publicacao.ID, nil)
		q = q.Filter("__key__ =", key)
	}

	q = q.KeysOnly()
	keys, err := datastoreClient.GetAll(c, q, nil)
	if err != nil {
		log.Warningf(c, "Erro ao buscar Publicação")
		return nil, err
	}
	return GetMultPublicacao(c, keys)
}

func Buscar(c context.Context, usuarioID int64) ([]Publicacao, error) {
	var publicacao Publicacao

	publicacao.AutorID = usuarioID

	publics, err := FiltrarPublicacoes(c, publicacao)
	if err != nil {
		log.Warningf(c, "Erro ao filtrar publicações pelo usuarioID: %v", err)
		return nil, err
	}

	seguidos, err := seguidores.BuscarUsuariosSeguidos(c, usuarioID)
	if err != nil {
		log.Warningf(c, "Erro ao buscar usuarios seguidos para perara busca de publicações %v", err)
		return nil, err
	}

	for _, v := range seguidos {
		publicacao.AutorID = v.ID
		publicSeguidos, err := FiltrarPublicacoes(c, publicacao)
		if err != nil {
			log.Warningf(c, "Erro filtrar publicações pelo usuarioID dos seguidos: %v", err)
			return nil, err
		}
		publics = append(publics, publicSeguidos...)

	}

	//ordenando publicações pela data mais recente
	sort.Slice(publics, func(i, j int) bool {
		return publics[i].DataCriacao.After(publics[j].DataCriacao.Time)
	})

	return publics, nil
	// TODO: Corrigir erro de bad request ocorrido quando um usuario que não segue ngm efetua a busca de publicações
}

func Atualizar(c context.Context, publicacao Publicacao) error {
	publicBanco := GetPublicacao(c, publicacao.ID)

	publicBanco.Titulo = publicacao.Titulo
	publicBanco.Conteudo = publicacao.Conteudo

	publicacao.AutorID = publicBanco.AutorID
	publicacao.AutorNick = publicBanco.AutorNick

	publicacao.DataCriacao = utils.GetSpecialTimeNow()

	return PutPublicacao(c, &publicacao)
}

func Deletar(c context.Context, publicacao Publicacao) error {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o datastore: %v", err)
		return err
	}
	defer datastoreClient.Close()

	key := datastore.IDKey(KindPublicacoes, publicacao.ID, nil)
	if err = datastoreClient.Delete(c, key); err != nil {
		log.Warningf(c, "Falha ao deletar publicação: %v", err)
		return err
	}
	return nil
}

func BuscarPorUsuario(c context.Context, usuarioID int64) ([]Publicacao, error) {
	var publicacao Publicacao
	publicacao.AutorID = usuarioID

	publics, err := FiltrarPublicacoes(c, publicacao)
	if err != nil {
		log.Warningf(c, "Erro ao filtrar publicações %v", err)
		return nil, err
	}

	return publics, nil
}

func Curtir(c context.Context, publicacaoID int64) error {
	public := GetPublicacao(c, publicacaoID)

	public.Curtidas++

	if err := Atualizar(c, *public); err != nil {
		log.Warningf(c, "Erro ao atualizar curtida da publicação no banco: %v", err)
		return err
	}
	return nil
}

func Descurtir(c context.Context, publicacaoID int64) error {
	public := GetPublicacao(c, publicacaoID)

	if public.Curtidas > 0 {
		public.Curtidas--
	} else {
		log.Warningf(c, "Não tem como descurtir uma publicação que não foi curtido")
		return fmt.Errorf("Não tem como descurtir uma publicação que não foi curtida")
	}

	if err := Atualizar(c, *public); err != nil {
		log.Warningf(c, "Erro ao atualizar descurtida da publicação no banco: %v", err)
		return err
	}
	return nil
}
