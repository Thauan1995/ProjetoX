package modelos

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"webapp/src/config"
	"webapp/src/requisicoes"
)

//Representa uma pessoa utilizando a rede social
type Usuario struct {
	ID          int64
	Nome        string
	Email       string
	Nick        string
	CriadoEm    time.Time
	Seguidores  []Usuario
	Seguindo    []Usuario
	Publicacoes []Publicacao
}

// Faz 4 requisições na API para maontar o perfil do usuario
func BuscarUsuarioCompleto(usuarioID int64, r *http.Request) (Usuario, error) {
	canalUsuario := make(chan Usuario)
	canalSeguidores := make(chan []Usuario)
	canalSeguindo := make(chan []Usuario)
	canalPublicacoes := make(chan []Publicacao)

	go BuscarDadosUsuario(canalUsuario, usuarioID, r)
	go BuscarSeguidores(canalSeguidores, usuarioID, r)
	go BuscarSeguindo(canalSeguindo, usuarioID, r)
	go BuscarPublicacoes(canalPublicacoes, usuarioID, r)

	var (
		usuario     Usuario
		seguidores  []Usuario
		seguindo    []Usuario
		publicacoes []Publicacao
	)

	for i := 0; i < 4; i++ {
		select {
		case usuarioCarregado := <-canalUsuario:
			if usuarioCarregado.ID == 0 {
				return Usuario{}, errors.New("Erro ao buscar o usuario")
			}
			usuario = usuarioCarregado

		case seguidoresCarregado := <-canalSeguidores:
			if seguidoresCarregado == nil {
				return Usuario{}, errors.New("Erro ao buscar os seguidores")
			}
			seguidores = seguidoresCarregado

		case seguindoCarregado := <-canalSeguindo:
			if seguindoCarregado == nil {
				return Usuario{}, errors.New("Erro ao buscar usuarios seguindo outros usuarios")
			}
			seguindo = seguindoCarregado

		case publicacoesCarregada := <-canalPublicacoes:
			if publicacoesCarregada == nil {
				return Usuario{}, errors.New("Erro ao buscar publcações do usuario")
			}
			publicacoes = publicacoesCarregada
		}
	}

	usuario.Seguidores = seguidores
	usuario.Seguindo = seguindo
	usuario.Publicacoes = publicacoes
	return usuario, nil

}

//Chama API para buscar os dados do usuario
func BuscarDadosUsuario(canal chan<- Usuario, usuarioID int64, r *http.Request) {
	url := fmt.Sprintf("%s/usuarios/buscar?ID=%d", config.ApiUrl, usuarioID)
	resp, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if err != nil {
		canal <- Usuario{}
		return
	}
	defer resp.Body.Close()

	var usuario Usuario
	if err = json.NewDecoder(resp.Body).Decode(&usuario); err != nil {
		canal <- Usuario{}
		return
	}

	canal <- usuario
}

//Chama API para buscar os seguidores do usuario
func BuscarSeguidores(canal chan<- []Usuario, usuarioID int64, r *http.Request) {
	url := fmt.Sprintf("%s/usuario/seguidores/%d", config.ApiUrl, usuarioID)
	resp, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if err != nil {
		canal <- nil
		return
	}
	defer resp.Body.Close()

	var seguidores []Usuario
	if err = json.NewDecoder(resp.Body).Decode(&seguidores); err != nil {
		canal <- nil
		return
	}

	canal <- seguidores
}

//Chama API para buscar usuarios seguidos por outros usuarios
func BuscarSeguindo(canal chan<- []Usuario, usuarioID int64, r *http.Request) {
	url := fmt.Sprintf("%s/usuario/seguidos/%d", config.ApiUrl, usuarioID)
	resp, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if err != nil {
		canal <- nil
		return
	}
	defer resp.Body.Close()

	var seguindo []Usuario
	if err = json.NewDecoder(resp.Body).Decode(&seguindo); err != nil {
		canal <- nil
		return
	}

	canal <- seguindo
}

//Chama API para buscar as publicações do usuario
func BuscarPublicacoes(canal chan<- []Publicacao, usuarioID int64, r *http.Request) {
	url := fmt.Sprintf("%s/usuario/%d/publicacoes", config.ApiUrl, usuarioID)
	resp, err := requisicoes.FazerRequisicaoComAutenticacao(r, http.MethodGet, url, nil)
	if err != nil {
		canal <- nil
		return
	}
	defer resp.Body.Close()

	var publicacoes []Publicacao
	if err = json.NewDecoder(resp.Body).Decode(&publicacoes); err != nil {
		canal <- nil
		return
	}

	canal <- publicacoes
}
