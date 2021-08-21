package wrap

import (
	"context"
	"net/http"
	"site/utils"
	"site/utils/log"
	"site/utils/wrap/keys"
)

func Wrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()

		logger, err := log.New(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, 0, "Erro ao conectar-se com o Logger")
			return
		}
		defer logger.Close()

		c = context.WithValue(c, keys.LoggerKey, logger)
		h.ServeHTTP(w, r.WithContext(c))
	})
}
