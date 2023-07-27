package logrous

import (
	"os"
	"time"

	"github.com/nolleh/caption_json_formatter"
	"github.com/sirupsen/logrus"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/models"
)

type logrusLogger struct {
	level      string
	encoding   string
	logger     *logrus.Logger
	logOptions *config2.LogOptions
}

// For mapping config logger
var loggerLevelMap = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
}

// NewLogrusLogger creates a new logrus logger
func NewLogrusLogger(cfg *config2.LogOptions, env environemnt.Environment) logger.Logger {
	logrusLogger := &logrusLogger{level: cfg.LogLevel, logOptions: cfg}
	logrusLogger.initLogger(env)

	return logrusLogger
}

// InitLogger Init logger
func (l *logrusLogger) initLogger(env environemnt.Environment) {
	logLevel := l.GetLoggerLevel()

	// Create a new instance of the logger. You can have any number of instances.
	logrusLogger := logrus.New()

	logrusLogger.SetLevel(logLevel)

	// Output to stdout instead of the defaultLogger stderr
	// Can be any io.Writer, see below for File example
	logrusLogger.SetOutput(os.Stdout)

	if env.IsDevelopment() {
		logrusLogger.SetReportCaller(false)
		logrusLogger.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			ForceColors:   true,
			FullTimestamp: true,
		})
	} else {
		logrusLogger.SetReportCaller(false)
		// https://github.com/nolleh/caption_json_formatter
		logrusLogger.SetFormatter(&caption_json_formatter.Formatter{PrettyPrint: true})
	}

	l.logger = logrusLogger
}

func (l *logrusLogger) GetLoggerLevel() logrus.Level {
	level, exist := loggerLevelMap[l.level]
	if !exist {
		return logrus.DebugLevel
	}

	return level
}

func (l *logrusLogger) LogType() models.LogType {
	return models.Logrus
}

func (l *logrusLogger) Configure(cfg func(internalLog interface{})) {
	cfg(l.logger)
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *logrusLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *logrusLogger) Debugw(msg string, fields logger.Fields) {
	entry := l.mapToFields(fields)
	entry.Debug(msg)
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *logrusLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *logrusLogger) Infow(msg string, fields logger.Fields) {
	entry := l.mapToFields(fields)
	entry.Info(msg)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *logrusLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *logrusLogger) WarnMsg(msg string, err error) {
	l.logger.Warn(msg, logrus.WithField("error", err.Error()))
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *logrusLogger) Errorw(msg string, fields logger.Fields) {
	entry := l.mapToFields(fields)
	entry.Error(msg)
}

func (l *logrusLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *logrusLogger) Err(msg string, err error) {
	l.logger.Error(msg, logrus.WithField("error", err.Error()))
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *logrusLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *logrusLogger) Printf(template string, args ...interface{}) {
	l.logger.Printf(template, args...)
}

func (l *logrusLogger) WithName(name string) {
	l.logger.WithField(constants.NAME, name)
}

func (l *logrusLogger) GrpcMiddlewareAccessLogger(
	method string,
	time time.Duration,
	metaData map[string][]string,
	err error,
) {
	l.Info(
		constants.GRPC,
		logrus.WithField(constants.METHOD, method),
		logrus.WithField(constants.TIME, time),
		logrus.WithField(constants.METADATA, metaData),
		logrus.WithError(err),
	)
}

func (l *logrusLogger) GrpcClientInterceptorLogger(
	method string,
	req interface{},
	reply interface{},
	time time.Duration,
	metaData map[string][]string,
	err error,
) {
	l.Info(
		constants.GRPC,
		logrus.WithField(constants.METHOD, method),
		logrus.WithField(constants.REQUEST, req),
		logrus.WithField(constants.REPLY, reply),
		logrus.WithField(constants.TIME, time),
		logrus.WithField(constants.METADATA, metaData),
		logrus.WithError(err),
	)
}

func (l *logrusLogger) mapToFields(fields map[string]interface{}) *logrus.Entry {
	return l.logger.WithFields(logrus.Fields{"ss": 1})
}
