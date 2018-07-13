package utility

import (
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/richardsang2008/BotManager/model"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"time"
)

type Log struct {
}

var log *logrus.Logger

func (l *Log) New(filename string, logLevel model.LogLevel) *logrus.Logger {
	if log != nil {
		return log
	}
	path := filename
	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	logrus.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.DebugLevel: writer,
			logrus.WarnLevel:  writer,
			logrus.FatalLevel: writer,
			logrus.PanicLevel: writer,
		},
		//&logrus.JSONFormatter{},
		&logrus.TextFormatter{},
	))
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  path,
		logrus.ErrorLevel: path,
		logrus.DebugLevel: path,
		logrus.WarnLevel:  path,
		logrus.FatalLevel: path,
		logrus.PanicLevel: path,
	}
	log = logrus.New()
	switch logLevel {
	case model.LogLevelDEBUG:
		log.SetLevel(logrus.DebugLevel)
		log.Hooks.Add(lfshook.NewHook(
			pathMap,
			&logrus.TextFormatter{},
		))
	case model.LogLevelINFO:
		log.SetLevel(logrus.InfoLevel)
		log.Hooks.Add(lfshook.NewHook(
			pathMap,
			&logrus.JSONFormatter{},
		))
	case model.LogLevelERROR:
		log.SetLevel(logrus.ErrorLevel)
		log.Hooks.Add(lfshook.NewHook(
			pathMap,
			&logrus.JSONFormatter{},
		))
	case model.LogLevelWARNING:
		log.SetLevel(logrus.WarnLevel)
		log.Hooks.Add(lfshook.NewHook(
			pathMap,
			&logrus.JSONFormatter{},
		))
	case model.LogLevelPANIC:
		log.SetLevel(logrus.FatalLevel)
		log.Hooks.Add(lfshook.NewHook(
			pathMap,
			&logrus.JSONFormatter{},
		))
	default:
		log.SetLevel(logrus.PanicLevel)
		log.Hooks.Add(lfshook.NewHook(
			pathMap,
			&logrus.JSONFormatter{},
		))
	}
	return log
}
func (l *Log) GetLogger() *logrus.Logger {
	return log
}
func (l *Log) Debug(args ...interface{}) {
	log.Debug(args)
}
func (l *Log) Panic(args ...interface{}) {
	log.Panic(args)
}
func (l *Log) Info(args ...interface{}) {
	log.Info(args)
}
func (l *Log) Error(args ...interface{}) {
	log.Error(args)
}
func (l *Log) Warning(args ...interface{}) {
	log.Warning(args)
}
