package endereco

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"site/utils/log"
)

type retPostmon struct {
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Cidade      string `json:"cidade"`
	Logradouro  string `json:"logradouro"`
	Cep         string `json:"cep"`
	Estado      string `json:"estado"`
}

type Endereco struct {
	CEP         string
	Numero      string
	Logradouro  string
	Bairro      string
	Municipio   string
	UF          string
	Pais        string
	Complemento string
}

func BuscarEndereco(c context.Context, endereco *Endereco) error {
	log.Debugf(c, "Buscando endereco por cep e numero %#v", endereco)

	if endereco.CEP == "" && endereco.Logradouro == "" && endereco.Bairro == "" && endereco.Municipio == "" {
		return fmt.Errorf("Dados insuficientes para busca do endereco %v", endereco)
	}

	return nil
}

func BuscaEnderecoPorCEP(c context.Context, endereco *Endereco) error {
	log.Debugf(c, "Buscando endereco por cep e numero %#v", endereco)
	urlCep := "http://api.postmon.com.br/v1/cep/%s"
	urlCep = fmt.Sprintf(urlCep, endereco.CEP)
	resp, err := http.Get(urlCep)
	log.Debugf(c, "Realizando requisicao para url %v", urlCep)
	if err != nil {
		log.Warningf(c, "Erro ao Realizar requisicao  %v", err)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warningf(c, "Erro ao Realizar ReadAll  %v", err)
		return err
	}
	log.Infof(c, "retorno da busca do cep realizada com sucesso. %#v", string(body))
	if len(body) == 0 {
		log.Warningf(c, "CEP não encontrado.")
		return fmt.Errorf("CEP não encontrado.")
	}
	var retorno retPostmon
	err = json.Unmarshal(body, &retorno)
	if err != nil {
		log.Warningf(c, "Erro ao Realizar Unmarshal  %v", err)
		return err
	}

	endereco.CEP = retorno.Cep
	endereco.Bairro = retorno.Bairro
	endereco.Logradouro = retorno.Logradouro
	endereco.UF = retorno.Estado

	log.Infof(c, "retornando endereco atualizado: %#v", retorno)
	return nil
}
