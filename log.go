//
// File: log.go
// Created by Dizzrt on 2023/02/17.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package etlog

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Dizzrt/etlog/color"
	"github.com/Dizzrt/etstream/kafka"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	ReporterType string
	ReporterName string
	FilePath     string
	MaxFileSize  int
	MaxBackups   int
	MaxAge       int
	Compress     bool
	KafkaEnable  bool
	KafkaConfig  kafka.KafkaConfig
}

const (
	logTimeFormat = "2006-01-02 15:04:05.000 Z07:00"
)

var (
	loggers       map[string]*zap.Logger
	defaultLogger *zap.Logger

	_globalMu sync.RWMutex
)

func init() {
	loggers = make(map[string]*zap.Logger)
}

func L() *zap.Logger {
	_globalMu.RLock()
	l := defaultLogger
	_globalMu.RUnlock()
	return l
}

func LogWithType(t string) *zap.Logger {
	return loggers[t]
}

func NewLogger(config LogConfig, logType string) (err error) {
	hook := &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxFileSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	bufferdWriteSyncer := zapcore.BufferedWriteSyncer{
		WS:   zapcore.AddSync(hook),
		Size: 4096,
	}

	cores := []zapcore.Core{
		// 输出 info level 及以上级别日志到文件
		zapcore.NewCore(logEncoder(), zapcore.AddSync(&bufferdWriteSyncer), zapcore.InfoLevel),
		// 输出 debug level 及以上级别日志到控制台
		zapcore.NewCore(logColorEncoder(), zapcore.Lock(os.Stdout), zapcore.DebugLevel),
	}

	if config.KafkaEnable {
		kafkaWriter, err := kafka.NewKafkaWriter(config.KafkaConfig, nil, nil)
		if err != nil {
			return err
		}

		kafkaWriteSyncer := zapcore.AddSync(kafkaWriter)
		// 输出 info level 及以上级别日志到kafka
		kafkacore := zapcore.NewCore(logEncoder(), kafkaWriteSyncer, zapcore.InfoLevel)
		cores = append(cores, kafkacore)
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core).WithOptions(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel), zap.Fields(
		zap.Any("log_info", map[string]interface{}{
			"ptype":    config.ReporterType, // producer type
			"pname":    config.ReporterName, // producer name
			"log_type": logType,
		})))

	if defaultLogger == nil {
		defaultLogger = logger
	}
	loggers[logType] = logger

	return nil
}

func logEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    logLevelEncoder,
			EncodeTime:     logTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   logCallerEncoder,
		})
}

func logColorEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    logLevelColorEncoder,
			EncodeTime:     logTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   logCallerEncoder,
		})
}

func logLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func logLevelColorEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("[%s]", color.LevelColorEncoder(level)))
}

func logTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(logTimeFormat) + "]")
}

func logCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}
