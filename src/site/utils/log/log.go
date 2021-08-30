package log

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"site/utils/consts"
	"site/utils/wrap/keys"
	"strings"
	"time"

	"cloud.google.com/go/logging"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

const (
	// DefaultLogID is the default log ID of the underlying Stackdriver Logging logger. Request
	// logs are logged under the ID "request_log", so use "app_log" for consistency. To use a
	// different ID create your logger with NewWithID.
	DefaultLogID = "app_log"

	traceContextHeaderName = "X-Cloud-Trace-Context"
)

func traceID(projectID, trace string) string {
	return fmt.Sprintf("projects/%s/traces/%s", projectID, trace)
}

type envVarError struct {
	varName string
}

func (e *envVarError) Error() string {
	return fmt.Sprintf("log: %s env var is not set, falling back to standard library log", e.varName)
}

// A Logger logs messages to Stackdriver Logging (though in certain cases it may fall back to the
// standard library's "log" package; see New). Logs will be correlated with requests in Stackdriver.
type Logger struct {
	client *logging.Client
	logger *logging.Logger
	monRes *monitoredres.MonitoredResource
	trace  string
}

// NewWithID creates a new Logger. The Logger is initialized using environment variables that are
// present on App Engine:
//
//   • GOOGLE_CLOUD_PROJECT
//   • GAE_SERVICE
//   • GAE_VERSION
//
// The given log ID will be passed through to the underlying Stackdriver Logging logger.
//
// Additionally, options (of type LoggerOption, from cloud.google.com/go/logging) will be passed
// through to the underlying Stackdriver Logging logger. Note that the option CommonResource will
// have no effect because the MonitoredResource is set when each log entry is made, thus overriding
// any value set with CommonResource. This is intended: much of the value of this package is in
// setting up the MonitoredResource so that log entries correlate with requests.
//
// The Logger will be valid in all cases, even when the error is non-nil. In the case of a non-nil
// error the Logger will fall back to the standard library's "log" package. There are three cases
// in which the error will be non-nil:
//
//   1. Any of the aforementioned environment variables are not set.
//   2. The given http.Request does not have the X-Cloud-Trace-Context header.
//   3. Initialization of the underlying Stackdriver Logging client produced an error.
func NewWithID(r *http.Request, logID string, options ...logging.LoggerOption) (*Logger, error) {
	traceContext := r.Header.Get(traceContextHeaderName)
	if traceContext == "" {
		return &Logger{}, fmt.Errorf("log: %s header is not set, falling back to standard library log", traceContextHeaderName)
	}

	client, err := logging.NewClient(r.Context(), fmt.Sprintf("projects/%s", consts.IDProjeto))
	if err != nil {
		return &Logger{}, err
	}

	monRes := &monitoredres.MonitoredResource{
		Labels: map[string]string{
			"module_id":  consts.IDServico,
			"project_id": consts.IDProjeto,
			"version_id": consts.IDVersao,
		},
		Type: "gae_app",
	}

	return &Logger{
		client: client,
		logger: client.Logger(logID, options...),
		monRes: monRes,
		trace:  traceID(consts.IDProjeto, strings.Split(traceContext, "/")[0]),
	}, nil
}

// New is identical to NewWithID with the exception that it uses the default log ID.
func New(r *http.Request, options ...logging.LoggerOption) (*Logger, error) {
	return NewWithID(r, DefaultLogID, options...)
}

// Close closes the Logger, ensuring all logs are flushed and closing the underlying
// Stackdriver Logging client.
func (lg *Logger) Close() error {
	if lg.client != nil {
		return lg.client.Close()
	}

	return nil
}

// Logf logs with the given severity. Remaining arguments are handled in the manner of fmt.Printf.
func (lg *Logger) logf(severity logging.Severity, format string, v ...interface{}) {
	if lg.logger == nil {
		log.Printf(format, v...)
		return
	}

	const lim = 240 << 10
	payload := fmt.Sprintf(format, v...)

	if len(payload) > lim {
		suffix := fmt.Sprintf("...(length %d)", len(payload))
		payload = fmt.Sprintf("%s%s", payload[:lim-len(suffix)], suffix)
	}

	lg.logger.Log(logging.Entry{
		Timestamp: time.Now(),
		Severity:  severity,
		Payload:   payload,
		Trace:     lg.trace,
		Resource:  lg.monRes,
	})
}

func Logf(ctx context.Context, severity logging.Severity, format string, v ...interface{}) {
	cv := ctx.Value(keys.LoggerKey)
	if cv == nil {
		log.Printf(format, v...)
		return
	}

	logger := cv.(*Logger)
	logger.logf(severity, format, v...)
}

// Debugf calls Logf with debug severity.
func Debugf(ctx context.Context, format string, v ...interface{}) {
	Logf(ctx, logging.Debug, format, v...)
}

// Infof calls Logf with info severity.
func Infof(ctx context.Context, format string, v ...interface{}) {
	Logf(ctx, logging.Info, format, v...)
}

// Noticef calls Logf with notice severity.
func Noticef(ctx context.Context, format string, v ...interface{}) {
	Logf(ctx, logging.Notice, format, v...)
}

// Warningf calls Logf with warning severity.
func Warningf(ctx context.Context, format string, v ...interface{}) {
	Logf(ctx, logging.Warning, format, v...)
}

// Errorf calls Logf with error severity.
func Errorf(ctx context.Context, format string, v ...interface{}) {
	Logf(ctx, logging.Error, format, v...)
}

// Criticalf calls Logf with critical severity.
func Criticalf(ctx context.Context, format string, v ...interface{}) {
	Logf(ctx, logging.Critical, format, v...)
}

// Alertf calls Logf with alert severity.
func Alertf(ctx context.Context, format string, v ...interface{}) {
	Logf(ctx, logging.Alert, format, v...)
}

// Emergencyf calls Logf with emergency severity.
func Emergencyf(ctx context.Context, format string, v ...interface{}) {
	Logf(ctx, logging.Emergency, format, v...)
}
