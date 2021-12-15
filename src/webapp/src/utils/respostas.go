package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// Representa a resposta de erro da API
type ErroAPI struct {
	Erro string `json:"err"`
}

//Retorna em formato JSON para a requisição
func JSON(w http.ResponseWriter, statusCode int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if statusCode != http.StatusNoContent {
		if err := json.NewEncoder(w).Encode(dados); err != nil {
			log.Fatal(err)
		}
	}
}

//Trata as requisições com status code 400 ou superior
func TratarStatusCodeErro(w http.ResponseWriter, r *http.Response) {
	var err ErroAPI
	json.NewDecoder(r.Body).Decode(&err)
	JSON(w, r.StatusCode, err)
}
