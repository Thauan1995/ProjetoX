package estabelecimento

import (
	"context"
	"fmt"
	"site/endereco"
	"site/utils"
	"site/utils/consts"
	"site/utils/log"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/badoux/checkmail"
)

const (
	KindEstabelecimento = "Estabelecimento"
)

type Estabelecimento struct {
	ID           int64 `datastore:"-"`
	CNPJ         string
	Endereco     EnderecoEstabelecimento
	IE           string
	Nome         string
	Email        string
	Telefone     string
	Setor        int64 // SELECT 1 - BAR 2 - RESTAURANTE 3 - LANCHONETE 4 - OUTROS
	HorarioFunc  string
	DiasFunc     int64 // 1 = Domingo 2 = Segunda 3 = Terça 4 = Quarta 5 = Quinta 6 = Sexta 7 = Sábado
	DataCadastro time.Time
}

type EnderecoEstabelecimento struct {
	CEP         string
	Numero      string
	Logradouro  string
	Bairro      string
	Municipio   string
	UF          string
	Pais        string
	Complemento string
}

func GetEstabelecimento(c context.Context, id int64) *Estabelecimento {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Erro ao conectar-se com o Datastore: %v", err)
		return nil
	}
	defer datastoreClient.Close()

	key := datastore.IDKey(KindEstabelecimento, id, nil)
	var estabelecimento Estabelecimento
	err = datastoreClient.Get(c, key, &estabelecimento)
	if err != nil {
		log.Warningf(c, "Erro ao buscar Empresa: %v", err)
		return nil
	}
	estabelecimento.ID = id
	return &estabelecimento
}
func GetMultiEstabelecimento(c context.Context, keys []*datastore.Key) ([]Estabelecimento, error) {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Erro ao conectar-se com o Datastore: %v", err)
		return []Estabelecimento{}, err
	}
	defer datastoreClient.Close()

	estabelecimentos := make([]Estabelecimento, len(keys))
	if err := datastoreClient.GetMulti(c, keys, estabelecimentos); err != nil {
		if errs, ok := err.(datastore.MultiError); ok {
			for _, e := range errs {
				if e == datastore.ErrNoSuchEntity {
					return []Estabelecimento{}, nil
				}
			}
		}
		log.Warningf(c, "Erro ao buscar Multi Estabelecimentos: %v", err)
		return []Estabelecimento{}, err
	}
	for i := range keys {
		estabelecimentos[i].ID = keys[i].ID
	}
	return estabelecimentos, nil
}

func PutEstabelecimento(c context.Context, estabelecimento *Estabelecimento) error {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Erro ao conectar-se com o Datastore: %v", err)
		return err
	}

	defer datastoreClient.Close()

	key := datastore.IDKey(KindEstabelecimento, estabelecimento.ID, nil)
	key, err = datastoreClient.Put(c, key, estabelecimento)
	if err != nil {
		log.Warningf(c, "Erro ao inserir Estabelecimento: %v", err)
		return err
	}

	estabelecimento.ID = key.ID
	return nil

}

func PutMultiEstabelecimentos(c context.Context, estabelecimentos []Estabelecimento) error {
	if len(estabelecimentos) == 0 {
		return nil
	}

	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Erro ao conectar-se com o Datastore: %v", err)
		return err
	}

	defer datastoreClient.Close()

	keys := make([]*datastore.Key, 0, len(estabelecimentos))

	for i := range estabelecimentos {
		keys = append(keys, datastore.IDKey(KindEstabelecimento, estabelecimentos[i].ID, nil))
		if err != nil {

			log.Warningf(c, "Erro ao inserir Multi Estabelecimentos: %v", err)
		}
	}
	return nil
}

func InserirEstabelecimento(c context.Context, estabelecimento *Estabelecimento) error {
	log.Debugf(c, "Inserindo Estabelecimento: %#v", estabelecimento)

	validEmail := checkmail.ValidateFormat(estabelecimento.Email)

	if validEmail != nil {

		return fmt.Errorf("Email inválido")
	}

	var cnpjOk bool

	estabelecimento.CNPJ, cnpjOk = utils.NormalizeCPFCNPJ(estabelecimento.CNPJ)

	if !cnpjOk {
		return fmt.Errorf("CNPJ Inválido")

	}

	if estabelecimento.IE == "" {
		return fmt.Errorf("IE Inválido")

	}

	if estabelecimento.Nome == "" {
		return fmt.Errorf("Nome deve ser informado")
	}

	estabelecimento.Telefone = utils.OnlyNumbers(estabelecimento.Telefone)
	if len(estabelecimento.Telefone) < 10 || len(estabelecimento.Telefone) > 11 {
		return fmt.Errorf("Telefone informado é inválido")

	}

	if estabelecimento.ID == 0 {

		enderEstab := endereco.Endereco{
			CEP:         estabelecimento.Endereco.CEP,
			Numero:      estabelecimento.Endereco.Numero,
			Logradouro:  estabelecimento.Endereco.Logradouro,
			Municipio:   estabelecimento.Endereco.Municipio,
			Bairro:      estabelecimento.Endereco.Bairro,
			UF:          estabelecimento.Endereco.UF,
			Pais:        estabelecimento.Endereco.Pais,
			Complemento: estabelecimento.Endereco.Complemento,
		}

		err := endereco.BuscaEnderecoPorCEP(c, &enderEstab)

		if err != nil {
			return err
		}

		estabelecimento.Endereco.CEP = enderEstab.CEP
		estabelecimento.Endereco.Numero = enderEstab.Numero
		estabelecimento.Endereco.Logradouro = enderEstab.Logradouro
		estabelecimento.Endereco.Municipio = enderEstab.Municipio
		estabelecimento.Endereco.Bairro = enderEstab.Bairro
		estabelecimento.Endereco.UF = enderEstab.UF
		estabelecimento.Endereco.Pais = enderEstab.Pais
		estabelecimento.Endereco.Complemento = enderEstab.Complemento

	}

	if estabelecimento.Endereco.CEP == "" {

		return fmt.Errorf("CEP Inválido")

	}

	if estabelecimento.Endereco.Bairro == "" {

		return fmt.Errorf("Bairro inválido")

	}

	if estabelecimento.Endereco.Logradouro == "" {

		return fmt.Errorf("Logradouro inválido")
	}

	if estabelecimento.Endereco.UF == "" {

		return fmt.Errorf("UF inválido")

	}

	if estabelecimento.Endereco.Pais == "" {

		return fmt.Errorf("País inválido")

	}

	if estabelecimento.Endereco.Numero == "" {

		return fmt.Errorf("Número inválido")

	}

	if estabelecimento.Endereco.Municipio == "" {

		return fmt.Errorf("Municipio inválido")

	}

	return PutEstabelecimento(c, estabelecimento)

}

func FiltrarEstabelecimento(c context.Context, estabelecimento Estabelecimento) ([]Estabelecimento, error) {
	datastoreClient, err := datastore.NewClient(c, consts.IDProjeto)
	if err != nil {
		log.Warningf(c, "Erro ao conectar-se com Datastore: %v", err)
		return nil, err
	}

	defer datastoreClient.Close()

	j := datastore.NewQuery(KindEstabelecimento)

	if estabelecimento.Nome != "" {

		j = j.Filter("Nome =", estabelecimento.Nome)

	}

	if estabelecimento.CNPJ != "" {

		j = j.Filter("CNPJ =", estabelecimento.CNPJ)
	}

	if estabelecimento.IE != "" {

		j = j.Filter("IE =", estabelecimento.IE)
	}

	if estabelecimento.ID != 0 {

		key := datastore.IDKey(KindEstabelecimento, estabelecimento.ID, nil)

		j = j.Filter("__key__=", key)
	}

	j = j.KeysOnly()
	keys, err := datastoreClient.GetAll(c, j, nil)
	if err != nil {
		log.Warningf(c, "Erro ao buscar estabelecimento: %v", err)
		return nil, err
	}

	return GetMultiEstabelecimento(c, keys)

}

func (estabelecimento *Estabelecimento) Validar(etapa string) error {

	if estabelecimento.Nome == "" {

		return fmt.Errorf("O campo nome é obrigatório: %v", estabelecimento.Nome)

	}

	if estabelecimento.IE == "" {

		return fmt.Errorf("O campo IE é obrigatório: %v", estabelecimento.IE)
	}

	if estabelecimento.CNPJ == "" {

		return fmt.Errorf("O campo CNPJ é obrigatório: %v", estabelecimento.CNPJ)
	}

	return nil
}
