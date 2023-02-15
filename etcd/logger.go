package etcd

import (
	"github.com/eolinker/eosc/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Build struct {
	logger  *Logger
	builder log.Builder
}

func (b *Build) Enabled(level zapcore.Level) bool {
	return b.logger.Enabled(level)
}

func (b *Build) With(fields []zapcore.Field) zapcore.Core {

	return &Build{
		logger:  b.logger,
		builder: b.builder.WithFields(toFields(fields)),
	}
}

func (b *Build) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return b.logger.Check(ent, ce)
}

func (b *Build) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	if b.logger.Enabled(ent.Level) {
		b.builder.WithFields(toFields(fields)).Logln(toLevel(ent.Level), ent.Message)
	}
	return nil
}

func (b *Build) Sync() error {
	return nil
}

type Logger struct {
	logger  *log.Logger
	encoder zapcore.Encoder
}

func NewLogger() *Logger {

	return &Logger{
		logger:  log.GetLogger(),
		encoder: zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
	}
}

func (l *Logger) Enabled(level zapcore.Level) bool {

	return l.logger.IsLevelEnabled(toLevel(level))

}

func (l *Logger) With(fields []zapcore.Field) zapcore.Core {

	return &Build{
		logger:  l,
		builder: l.logger.WithFields(toFields(fields)),
	}
}

func (l *Logger) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(ent, l)
}

func (l *Logger) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	level := toLevel(ent.Level)
	if l.logger.IsLevelEnabled(level) {
		buf, err := l.encoder.EncodeEntry(ent, fields)
		if err != nil {
			return err
		}
		l.logger.Logln(level, buf.String())
	}
	return nil
}

func (l *Logger) Sync() error {
	return nil
}

var (
	levels = map[zapcore.Level]log.Level{
		zapcore.DebugLevel: log.DebugLevel,
		// InfoLevel is the default logging priority.
		zapcore.InfoLevel: log.InfoLevel,
		// WarnLevel logs are more important than Info, but don't need individual
		// human review.
		zapcore.WarnLevel: log.WarnLevel,
		// ErrorLevel logs are high-priority. If an application is running smoothly,
		// it shouldn't generate any error-level logs.
		zapcore.ErrorLevel: log.ErrorLevel,
		// DPanicLevel logs are particularly important errors. In development the
		// logger panics after writing the message.
		zapcore.DPanicLevel: log.TraceLevel,
		// PanicLevel logs a message, then panics.
		zapcore.PanicLevel: log.PanicLevel,
		// FatalLevel logs a message, then calls os.Exit(1).
		zapcore.FatalLevel: log.FatalLevel,
	}
)

func toLevel(level zapcore.Level) log.Level {
	return levels[level]
}
func toFields(fields []zapcore.Field) log.Fields {

	encoder := zapcore.NewMapObjectEncoder()

	for _, f := range fields {

		f.AddTo(encoder)
	}
	return encoder.Fields
}
