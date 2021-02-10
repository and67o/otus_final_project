package logger

import (
	"errors"
	"fmt"
	"strings"

	"github.com/and67o/otus_project/internal/configuration"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

const (
	levelDebug = "DEBUG"
	levelError = "ERROR"
	levelInfo  = "INFO"
	levelWarn  = "WARN"

	outPut = "stderr"
)

var errNotFound = errors.New("not found log level")

func New(configuration configuration.LoggerConf) (*Logger, error) {
	logger := new(Logger)

	config := getConfig(configuration.IsProd)

	config.DisableStacktrace = configuration.TraceOn

	config.OutputPaths = []string{outPut}
	config.OutputPaths = []string{configuration.File}

	lvl, err := setLevel(configuration.Level)
	if err != nil {
		return nil, fmt.Errorf("level error: %w", err)
	}
	config.Level = *lvl

	l, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("config build error: %w", err)
	}

	logger.logger = l

	return logger, nil
}

func setLevel(level string) (*zap.AtomicLevel, error) {
	var lvl zap.AtomicLevel
	var l zapcore.Level

	switch strings.ToUpper(level) {
	case levelDebug:
		l = zapcore.DebugLevel
	case levelError:
		l = zapcore.ErrorLevel
	case levelInfo:
		l = zapcore.InfoLevel
	case levelWarn:
		l = zapcore.WarnLevel
	default:
		return nil, errNotFound
	}

	lvl = zap.NewAtomicLevelAt(l)

	return &lvl, nil
}

func getConfig(isProd bool) zap.Config {
	var config zap.Config

	if isProd {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	return config
}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l Logger) Fatal(msg string) {
	l.logger.Fatal(msg)
}

func (l Logger) Warn(msg string) {
	l.logger.Warn(msg)
}
