package seguidores

import (
	"context"
	"fmt"
	"site/usuario"
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
func GetMultSeguidor(c context.Context, keys []*datastore.Key) ([]Seguidor, error) {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return []Seguidor{}, err
	}
	defer datastoreClient.Close()

	seguidores := make([]Seguidor, len(keys))
	if err := datastoreClient.GetMulti(c, keys, seguidores); err != nil {
		if errs, ok := err.(datastore.MultiError); ok {
			for _, e := range errs {
				if e == datastore.ErrNoSuchEntity {
					return []Seguidor{}, nil
				}
			}
		}
		log.Warningf(c, "Erro ao buscar Multi Seguidors: %v", err)
		return []Seguidor{}, err
	}
	for i := range keys {
		seguidores[i].IDSeguidor = keys[i].ID
	}
	return seguidores, nil
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

func FiltrarSeguidores(c context.Context, filtro Seguidor) ([]Seguidor, error) {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Erro ao conectar-se com o Datastore: %v", err)
		return nil, err
	}
	defer datastoreClient.Close()

	q := datastore.NewQuery(KindSeguidores)

	if filtro.IDSeguidor != 0 {
		key := datastore.IDKey(KindSeguidores, filtro.IDSeguidor, nil)
		q = q.Filter("__key__ =", key)
	}

	q = q.KeysOnly()
	keys, err := datastoreClient.GetAll(c, q, nil)
	if err != nil {
		log.Warningf(c, "Erro ao buscar Seguidor: %v", err)
		return nil, err
	}
	return GetMultSeguidor(c, keys)
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

	keysSeguidor, err := FiltrarSeguidores(c, *seguidor)
	if err != nil {
		log.Warningf(c, "Erro ao filtrar seguidores %v", err)
		return err
	}
	if len(keysSeguidor) > 0 && seguidor.IDSeguidor != keysSeguidor[0].IDSeguidor {
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

	if len(usuariosAtt) < 1 {
		usuariosAtt = []int64{0}
	}

	seguidorBanco.IDSeguidor = seguidorID
	seguidorBanco.IDUsuario = usuariosAtt

	if err := InserirSeguidor(c, seguidorBanco); err != nil {
		log.Warningf(c, "Erro na inserção do seguidor atualizado: %v", err)
		return fmt.Errorf("Erro na inserção do seguidor atualizado")
	}
	return nil
}

func BuscarUsuariosSeguidos(c context.Context, seguidorID int64) ([]usuario.Usuario, error) {
	var usuarios []usuario.Usuario
	seguidorBanco := GetSeguidorByIDSeguidor(c, seguidorID)

	for _, v := range seguidorBanco.IDUsuario {
		if v == 0 {
			continue
		}
		usu := usuario.GetUsuario(c, v)
		usuarios = append(usuarios, usuario.Usuario{
			ID:          usu.ID,
			Nome:        usu.Nome,
			Nick:        usu.Nick,
			Email:       usu.Email,
			DataCriacao: usu.DataCriacao,
		})

	}
	return usuarios, nil
}

func BuscarSeguidores(c context.Context, seguidores []Seguidor, usuarioID int64) ([]usuario.Usuario, error) {
	var usuarios []usuario.Usuario

	for _, v := range seguidores {
		for _, x := range v.IDUsuario {
			if x == usuarioID {
				usu := usuario.GetUsuario(c, v.IDSeguidor)
				usuarios = append(usuarios, usuario.Usuario{
					ID:          usu.ID,
					Nome:        usu.Nome,
					Nick:        usu.Nick,
					Email:       usu.Email,
					DataCriacao: usu.DataCriacao,
				})
			}
		}
	}
	return usuarios, nil
}
