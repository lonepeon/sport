package logger

import (
	"errors"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ErrGeneric = errors.New("something wrong happened")
)

type Closer func() error

type Logger struct {
	logger *zap.SugaredLogger
	fields []Field
}

type wFlusher struct {
	io.Writer
}

func (w wFlusher) Sync() error {
	type Flusher interface {
		Flush() error
	}

	flusher, ok := w.Writer.(Flusher)
	if ok {
		flusher.Flush()
	}

	return nil
}

func NewLogger(w io.Writer) (*Logger, Closer) {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		wFlusher{w},
		zapcore.InfoLevel,
	)

	logger := zap.New(core).Sugar()
	closer := func() error { return logger.Sync() }

	return &Logger{logger: logger}, closer
}

func (l *Logger) Error(msg string) {
	l.logger.With(l.toZapFields()...).Error(msg)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.logger.With(l.toZapFields()...).Errorf(msg, args...)
}

func (l *Logger) Info(msg string) {
	l.logger.With(l.toZapFields()...).Infof(msg)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.logger.With(l.toZapFields()...).Infof(msg, args...)
}

func (l *Logger) WithFields(f ...Field) *Logger {
	return &Logger{
		logger: l.logger.With(),
		fields: append(l.fields, f...),
	}
}

func (l *Logger) toZapFields() []interface{} {
	fields := make([]interface{}, 0, len(l.fields))
	for _, field := range l.fields {
		fields = append(fields, field.ToZapField())
	}

	return fields
}
