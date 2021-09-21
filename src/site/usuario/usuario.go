package usuario

import (
	"context"
	"fmt"
	"site/utils"
	"site/utils/consts"
	"site/utils/log"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

const (
	KindUsuario = "Usuario"
)

type Usuario struct {
	ID          int64 `datastore:"-"`
	Nome        string
	Nick        string
	Email       string
	Senha       string
	DataCriacao time.Time
}

func GetUsuario(c context.Context, id int64) *Usuario {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return nil
	}
	defer datastoreClient.Close()

	key := datastore.IDKey(KindUsuario, id, nil)

	var usuario Usuario
	if err = datastoreClient.Get(c, key, &usuario); err != nil {
		log.Warningf(c, "Falha ao buscar Usuario: %v", err)
		return nil
	}
	usuario.ID = id
	return &usuario
}

func GetMultUsuario(c context.Context, keys []*datastore.Key) ([]Usuario, error) {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return []Usuario{}, err
	}
	defer datastoreClient.Close()

	usuario := make([]Usuario, len(keys))
	if err := datastoreClient.GetMulti(c, keys, usuario); err != nil {
		if errs, ok := err.(datastore.MultiError); ok {
			for _, e := range errs {
				if e == datastore.ErrNoSuchEntity {
					return []Usuario{}, nil
				}
			}
		}
		log.Warningf(c, "Erro ao buscar Multi Usuarios: %v", err)
		return []Usuario{}, err
	}
	for i := range keys {
		usuario[i].ID = keys[i].ID
	}
	return usuario, nil
}

func PutUsuario(c context.Context, usuario *Usuario) error {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se ao Datastore: %v", err)
		return err
	}
	defer datastoreClient.Close()

	key := datastore.IDKey(KindUsuario, usuario.ID, nil)
	key, err = datastoreClient.Put(c, key, usuario)
	if err != nil {
		log.Warningf(c, "Erro ao atualizar usuario: %v", err)
		return err
	}
	usuario.ID = key.ID
	return nil
}

func PutMultUsuario(c context.Context, usuario []Usuario) error {
	if len(usuario) == 0 {
		return nil
	}
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return err
	}
	defer datastoreClient.Close()

	keys := make([]*datastore.Key, 0, len(usuario))
	for i := range usuario {
		keys = append(keys, datastore.IDKey(KindUsuario, usuario[i].ID, nil))
	}
	keys, err = datastoreClient.PutMulti(c, keys, usuario)
	if err != nil {
		log.Warningf(c, "Erro ao inserir Multi Usuarios: %v", err)
		return err
	}
	return nil
}
func FiltrarUsuario(c context.Context, usuario Usuario) ([]Usuario, error) {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return nil, err
	}
	defer datastoreClient.Close()

	q := datastore.NewQuery(KindUsuario)

	if usuario.Nome != "" {
		q = q.Filter("Nome =", usuario.Nome)
	}

	if usuario.Nick != "" {
		q = q.Filter("Nick =", usuario.Nick)
	}

	if usuario.Email != "" {
		q = q.Filter("Email =", usuario.Email)
	}

	if usuario.ID != 0 {
		key := datastore.IDKey(KindUsuario, usuario.ID, nil)
		q = q.Filter("__key__ =", key)
	}

	q = q.KeysOnly()
	keys, err := datastoreClient.GetAll(c, q, nil)
	if err != nil {
		log.Warningf(c, "Erro ao buscar Usuario: %v", err)
		return nil, err
	}
	return GetMultUsuario(c, keys)
}

// validar() valida os campos do processo
func (usuario *Usuario) validar(etapa string) error {
	if usuario.Nome == "" {
		return fmt.Errorf("O campo nome é obrigatório: %v", usuario.Nome)
	}

	if usuario.Nick == "" {
		return fmt.Errorf("O campo nick é obrigatório: %v", usuario.Nick)
	}

	if usuario.Email == "" {
		return fmt.Errorf("O campo email é obrigatório: %v", usuario.Email)
	}

	if etapa == "cadastro" && usuario.Senha == "" {
		return fmt.Errorf("O campo senha é obrigatório")
	}

	return nil
}

func (usuario *Usuario) Preparar(etapa string) error {
	if err := usuario.validar(etapa); err != nil {
		return err
	}
	return nil
}
func InserirUsuario(c context.Context, usuario *Usuario) error {
	log.Debugf(c, "Inserindo Usuario no banco: %v", usuario)

	if err := usuario.Preparar("cadastro"); err != nil {
		return fmt.Errorf("Erro ao preparar campos de usuario %v", err)
	}

	validaEmail := checkmail.ValidateFormat(usuario.Email)
	if validaEmail != nil {
		return fmt.Errorf("Email inserido inválido")
	}

	cost := bcrypt.DefaultCost

	hash, err := bcrypt.GenerateFromPassword([]byte(usuario.Senha), cost)
	if err != nil {
		panic(err.Error())
	}
	usuario.Senha = string(hash)

	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nick = strings.TrimSpace(usuario.Nick)
	usuario.Email = strings.TrimSpace(usuario.Email)

	if usuario.ID == 0 {
		usuario.DataCriacao = utils.GetTimeNow()
	}

	return PutUsuario(c, usuario)
}

func GetUsuarioByEmail(c context.Context, email string) *Usuario {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Falha ao conectar-se com o Datastore: %v", err)
		return nil
	}
	defer datastoreClient.Close()

	key := datastore.NameKey(KindUsuario, email, nil)

	var usuario Usuario
	if err = datastoreClient.Get(c, key, &usuario); err != nil {
		log.Warningf(c, "Falha ao buscar Usuario: %v", err)
		return nil
	}
	usuario.Email = email
	return &usuario
}
