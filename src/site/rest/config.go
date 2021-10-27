package rest

import (
	"net/http"
	"site/config"
	"site/utils"
	"site/utils/log"
	"strings"
)

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	c := r.Context()

	switch r.Method {

	case http.MethodPost:
		nome := strings.TrimSpace(r.FormValue("nome"))
		valor := strings.TrimSpace(r.FormValue("valor"))
		if nome == "" || valor == "" {
			log.Warningf(c, "nome e valor são obrigatórios.")
			utils.RespondWithError(w, http.StatusBadRequest, 0, "nome e valor são obrigatorios")
			return
		}

		var configuracao config.Config
		configuracao.Name = nome
		configuracao.Value = valor

		err := config.PutConfig(c, &configuracao)
		if err != nil {
			log.Warningf(c, "Não foi possivel salvar o valor da configuração: %v", err)
			utils.RespondWithError(w, http.StatusBadRequest, 0, "Não foi possivel salvar o valor da configuração")
			return
		}

		log.Infof(c, "Configuração salva: %#v", configuracao)
		utils.RespondWithJSON(w, http.StatusOK, configuracao)
		return

	case http.MethodGet:
		nome := strings.TrimSpace(r.FormValue("nome"))

		if nome == "" {
			configs, err := config.ListConfigs(c)
			if err != nil {
				log.Warningf(c, "Erro ao encontrar lista %v", err)
				utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao encontrar lista")
				return
			}
			log.Infof(c, "Lista de config: %v", configs)
			utils.RespondWithJSON(w, http.StatusOK, configs)
			return
		}

		config, err := config.GetConfig(c, nome)
		if err != nil {
			log.Warningf(c, "Erro ao buscar config %v", err)
			utils.RespondWithError(w, http.StatusBadRequest, 0, "Erro ao buscar config")
			return
		}

		if config == nil {
			log.Warningf(c, "Nenhuma config com esse nome %v", nome)
			utils.RespondWithError(w, http.StatusBadRequest, 0, "Nenhuma config com esse nome")
			return
		}
		log.Infof(c, "Config encontrada %#v", config)
		utils.RespondWithJSON(w, http.StatusOK, config)
		return

	default:
		log.Warningf(c, "Metodo não permitido")
		utils.RespondWithError(w, http.StatusMethodNotAllowed, 0, "Metodo não permitido")
		return
	}
}
