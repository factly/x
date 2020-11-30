package loggerx

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// GormLogger custom logger for gorm queries
type GormLogger struct {
	logger.Writer
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		requestID := middleware.GetReqID(ctx)
		if requestID == "" {
			requestID = "-"
		}
		l.Printf(l.infoStr+msg, append([]interface{}{requestID, utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		requestID := middleware.GetReqID(ctx)
		if requestID == "" {
			requestID = "-"
		}
		l.Printf(l.warnStr+msg, append([]interface{}{requestID, utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		requestID := middleware.GetReqID(ctx)
		if requestID == "" {
			requestID = "-"
		}
		l.Printf(l.errStr+msg, append([]interface{}{requestID, utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		requestID := middleware.GetReqID(ctx)
		if requestID == "" {
			requestID = "-"
		}
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", requestID, sql)
			} else {
				l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, requestID, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", requestID, sql)
			} else {
				l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, requestID, sql)
			}
		case l.LogLevel >= logger.Info:
			sql, rows := fc()
			if rows == -1 {
				l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", requestID, sql)
			} else {
				l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, requestID, sql)
			}
		}
	}
}

type traceRecorder struct {
	logger.Interface
	BeginAt      time.Time
	SQL          string
	RowsAffected int64
	Err          error
}

func (l traceRecorder) New() *traceRecorder {
	return &traceRecorder{Interface: l.Interface, BeginAt: time.Now()}
}

func (l *traceRecorder) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.BeginAt = begin
	l.SQL, l.RowsAffected = fc()
	l.Err = err
}

// NewGormLogger get new gorm logger
func NewGormLogger(config logger.Config) GormLogger {
	return GormLogger{
		Writer:       log.New(os.Stdout, "\r\n", log.LstdFlags),
		infoStr:      "[request_id:%v] %s\n[info] ",
		warnStr:      "[request_id:%v] %s\n[warn] ",
		errStr:       "[request_id:%v] %s\n[error] ",
		traceStr:     "%s\n[%.3fms] [rows:%v] [request_id:%v] %s",
		traceWarnStr: "%s %s\n[%.3fms] [rows:%v] [request_id:%v] %s",
		traceErrStr:  "%s %s\n[%.3fms] [rows:%v] [request_id:%v] %s",
		Config:       config,
	}
}
