package log

import (
	"fmt"

	"os"
	"time"
)

var (
	logger          *Logger
	transport       *Complex
	lineFormatter   *LineFormatter
	stdoutTransport *Transporter
	isDebug         = false
	transportsCache []EntryTransporter
)

func init() {

	lineFormatter = &LineFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}
	stdoutTransport = NewTransport(os.Stderr, InfoLevel)
	stdoutTransport.SetFormatter(lineFormatter)
	transport = NewComplex()
	logger = NewLogger(transport, false, "")
	Reset()
	RegisterExitHandler(func() {
		Close()
	})
}
func InitDebug(d bool) {
	isDebug = d
	if isDebug {
		stdoutTransport.SetLevel(DebugLevel)
	} else {
		transport.setLevel(InfoLevel)
	}
	Reset(transportsCache...)
}

func Reset(transports ...EntryTransporter) {
	transportsCache = transports
	if isDebug || len(transports) == 0 {
		transportsTmp := append(transportsCache, stdoutTransport)
		transport.Reset(transportsTmp...)
	} else {
		transport.Reset(transports...)
	}

}

//Close 关闭
func Close() {
	if transport != nil {
		transport.Close()
	}
}
func SetPrefix(prefix string) {
	logger.SetPrefix(prefix)
}

//WithFields 写域
func WithFields(fields Fields) Builder {

	return logger.WithFields(fields)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {

	logger.Debug(args...)
}

// Debug logs a message at level Debug on the standard logger.
func DebugF(format string, args ...interface{}) {

	logger.Debugf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {

	logger.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {

	logger.Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {

	logger.Error(args...)
}
func panicOut(args ...interface{}) string {
	defer func() {
		if e := recover(); e != nil {
			Close()
		}
	}()
	s, _ := encode(PanicLevel, args...)
	_, _ = os.Stderr.Write(s)

	return string(s)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	panic(panicOut(args...))
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	s, e := encode(FatalLevel, args...)
	if e != nil {
		return
	}
	_, _ = os.Stderr.Write(s)
	logger.Fatal(args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {

	logger.Infof(format, args...)
}

// Warnf logs a message at level Info on the standard logger.
func Warnf(format string, args ...interface{}) {

	logger.Warnf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {

	logger.Errorf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {

	s, e := encode(FatalLevel, fmt.Sprintf(format, args...))
	if e != nil {
		return
	}
	_, _ = os.Stdout.Write(s)
	logger.Fatalf(format, args...)

}

func encode(level Level, args ...interface{}) ([]byte, error) {
	entry := &Entry{}
	entry.Message = fmt.Sprintln(args...)
	entry.Level = level
	entry.Time = time.Now()
	s, e := lineFormatter.Format(entry)
	if e != nil {
		return s, e
	}
	return s, e
}
