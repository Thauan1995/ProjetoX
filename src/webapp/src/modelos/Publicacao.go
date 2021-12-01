package modelos

import (
	"webapp/src/utils"
)

//Representa uma publicação feita por um usuario
type Publicacao struct {
	ID          int64
	Titulo      string
	Conteudo    string
	AutorID     int64
	AutorNick   string
	Curtidas    int64
	DataCriacao utils.JsonSpecialDateTime
}
