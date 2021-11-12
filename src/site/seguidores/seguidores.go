package seguidores

import (
	"context"
	"fmt"
	"site/utils"
	"site/utils/consts"
	"site/utils/log"

	"cloud.google.com/go/datastore"
)

const (
	KindSeguidores = "Seguidores"
)

type Seguidor struct {
	IDSeguidor  int64
	IDUsuario   []int64
	DataCriacao utils.JsonSpecialDateTime
}

func GetSeguidorByIDSeguidor(c context.Context, idSeguidor int64) *Seguidor {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return nil
	}
	defer datastoreClient.Close()

	key := datastore.IDKey(KindSeguidores, idSeguidor, nil)
	var seguidor Seguidor

	if err := datastoreClient.Get(c, key, &seguidor); err != nil {
		log.Warningf(c, "Falha na busca do seguidor : %v", err)
		return nil
	}

	seguidor.IDSeguidor = idSeguidor
	return &seguidor
}

func PutSeguidor(c context.Context, seguidor *Seguidor) error {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return err
	}
	defer datastoreClient.Close()

	key := datastore.IDKey(KindSeguidores, seguidor.IDSeguidor, nil)
	key, err = datastoreClient.Put(c, key, seguidor)
	if err != nil {
		log.Warningf(c, "Erro ao atualizar seguidor")
		return err
	}
	seguidor.IDSeguidor = key.ID
	return nil
}

func FiltrarSeguidores(c context.Context, filtro Seguidor) []*datastore.Key {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Erro ao conectar-se com o Datastore: %v", err)
		return nil
	}
	defer datastoreClient.Close()

	q := datastore.NewQuery(KindSeguidores)

	if filtro.IDSeguidor != 0 {
		key := datastore.IDKey(KindSeguidores, filtro.IDSeguidor, nil)
		q = q.Filter("__key__ =", key)
	}

	q = q.Order("IDSeguidor").KeysOnly()
	keys, err := datastoreClient.GetAll(c, q, nil)
	if err != nil {
		log.Warningf(c, "Erro ao buscar Seguidor: %v", err)
		return nil
	}
	return keys
}

func InserirSeguidor(c context.Context, seguidor *Seguidor) error {
	log.Debugf(c, "Inserindo Seguidor no banco: %v", seguidor)

	if seguidor.IDUsuario == nil {
		return fmt.Errorf("IDUsuario não informado: %v", seguidor.IDUsuario)
	}

	if seguidor.IDSeguidor == 0 {
		return fmt.Errorf("IDSeguidor não informado: %v", seguidor.IDSeguidor)
	}

	if seguidor.DataCriacao.IsZero() {
		seguidor.DataCriacao = utils.GetSpecialTimeNow()
	}

	keysSeguidor := FiltrarSeguidores(c, *seguidor)
	if len(keysSeguidor) > 0 && seguidor.IDSeguidor != keysSeguidor[0].ID {
		log.Debugf(c, "Seguidor ja existe %#v", seguidor)
		return nil
	}

	return PutSeguidor(c, seguidor)
}

func Seguir(c context.Context, usuarioID, seguidorID int64) error {
	var seguidor Seguidor
	seguidorBanco := GetSeguidorByIDSeguidor(c, seguidorID)

	if seguidorBanco == nil {
		seguidor.IDUsuario = append(seguidor.IDUsuario, usuarioID)

	} else {
		for _, v := range seguidorBanco.IDUsuario {
			if v == usuarioID {
				log.Warningf(c, "Usuario ja está sendo seguido")
				return fmt.Errorf("Usuario já está sendo seguido")
			}
		}
		seguidorBanco.IDUsuario = append(seguidorBanco.IDUsuario, usuarioID)
		seguidor.IDUsuario = append(seguidor.IDUsuario, seguidorBanco.IDUsuario...)
	}
	seguidor.IDSeguidor = seguidorID

	if err := InserirSeguidor(c, &seguidor); err != nil {
		log.Warningf(c, "Erro na inserção do seguidor no banco: %v", err)
		return fmt.Errorf("Erro na inserção do seguidor no banco")
	}
	return nil
}

func PararDeSeguir(c context.Context, usuarioID, seguidorID int64) error {
	seguidorBanco := GetSeguidorByIDSeguidor(c, seguidorID)

	var usuariosAtt []int64
	for _, v := range seguidorBanco.IDUsuario {
		if v != usuarioID {
			usuariosAtt = append(usuariosAtt, v)
		}
	}

	seguidorBanco.IDSeguidor = seguidorID
	seguidorBanco.IDUsuario = append(seguidorBanco.IDUsuario, usuariosAtt...)

	if err := InserirSeguidor(c, seguidorBanco); err != nil {
		log.Warningf(c, "Erro na inserção do seguidor atualizado: %v", err)
		return fmt.Errorf("Erro na inserção do seguidor atualizado")
	}
	return nil
}
