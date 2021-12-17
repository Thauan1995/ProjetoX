package modelos

import "time"

//Representa uma pessoa utilizando a rede social
type Usuario struct {
	ID         int64
	Nome       string
	Email      string
	Nick       string
	CriadoEm   time.Time
	Seguidores []Usuario
	Seguindo   []Usuario
	Pulicacoes []Publicacao
}
