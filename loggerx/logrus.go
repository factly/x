package loggerx

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"

	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

var logrusLogger *logrus.Logger
var req *http.Request

func Init() func(next http.Handler) http.Handler {
	logrusLogger = logrus.New()
	logrusLogger.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}
	return NewStructuredLogger(logrusLogger)
}

func NewStructuredLogger(logrusLogger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logrusLogger})
}

type StructuredLogger struct {
	Logger *logrus.Logger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {

	req = r

	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := GetDefaultFields()

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
	logFields := GetDefaultFields()

	if pc, file, line, ok := runtime.Caller(1); ok {
		funcName := runtime.FuncForPC(pc).Name()

		pwd, _ := os.Getwd()
		relPath := file[len(pwd):]

		logFields["source"] = fmt.Sprintf("%s:%s:%v", relPath, path.Base(funcName), line)
	}

	logrusLogger.WithFields(logFields).Error(err)
}

func ErrorWithoutRequest(err error) {
	logrusFields := logrus.Fields{}
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		pwd, _ := os.Getwd()
		relPath := file[len(pwd):]
		logrusFields["source"] = fmt.Sprintf("%s:%s:%v", relPath, path.Base(funcName), line)
	}
	logrusLogger.WithFields(logrusFields).Error(err)
}

func Info(info string) {
	logrusFields := logrus.Fields{}
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		pwd, _ := os.Getwd()
		relPath := file[len(pwd):]
		logrusFields["source"] = fmt.Sprintf("%s:%s:%v", relPath, path.Base(funcName), line)
	}
	logrusLogger.WithFields(logrusFields).Info(info)
}

func Warning(warning string) {
	logrusFields := logrus.Fields{}
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		pwd, _ := os.Getwd()
		relPath := file[len(pwd):]
		logrusFields["source"] = fmt.Sprintf("%s:%s:%v", relPath, path.Base(funcName), line)
	}
	logrusLogger.WithFields(logrusFields).Warning(warning)
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

func GetDefaultFields() logrus.Fields {
	logFields := logrus.Fields{}

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

	return logFields
}
