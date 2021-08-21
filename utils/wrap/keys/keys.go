package keys

type ctxKeyType string

var (
	LoggerKey    = ctxKeyType("logger-key")
	DatastoreKey = ctxKeyType("datastore-key")
)
