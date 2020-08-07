package loggerx

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger
var req *http.Request

func Init() func(next http.Handler) http.Handler {
	logger = logrus.New()
	logger.Formatter = &logrus.TextFormatter{}

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Println(err)
	}

	logger.Out = file

	return NewStructuredLogger(logger)
}

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

type StructuredLogger struct {
	Logger *logrus.Logger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {

	req = r

	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Infoln("request started")

	return entry
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Infoln("request complete")
}

func Error(err error) {
	logFields := logrus.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	if reqID := middleware.GetReqID(req.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if req.TLS != nil {
		scheme = "https"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = req.Proto
	logFields["http_method"] = req.Method

	logFields["remote_addr"] = req.RemoteAddr
	logFields["user_agent"] = req.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, req.Host, req.RequestURI)

	if pc, file, line, ok := runtime.Caller(1); ok {
		funcName := runtime.FuncForPC(pc).Name()

		pwd, _ := os.Getwd()
		relPath := file[len(pwd):]

		logFields["source"] = fmt.Sprintf("%s:%s:%v", relPath, path.Base(funcName), line)
	}

	logger.WithFields(logFields).Error(err)
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func GetLogEntry(r *http.Request) logrus.FieldLogger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
}

func LogEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithField(key, value)
	}
}

func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}
