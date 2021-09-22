package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"site/rest"

	"github.com/gorilla/mux"
)

func init() {
	chave := make([]byte, 64)

	if _, err := rand.Read(chave); err != nil {
		log.Fatal("Erro ao criar chave Secret", err)
	}

	stringBase64 := base64.StdEncoding.EncodeToString(chave)
	fmt.Println(stringBase64)
}

func main() {
	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()

	//Usuario
	r.HandleFunc("/usuario", rest.UsuarioHandler)
	r.HandleFunc("/usuario/login", rest.LoginHandler)

	http.Handle("/", router)

	var port = os.Getenv("PORT")
	if port == "" {
		port = "5000"
		log.Printf("Padronizando para porta %s", port)
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
