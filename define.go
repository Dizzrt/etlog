//
// File: define.go
// Created by Dizzrt on 2023/02/27.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package etlog

import (
	"time"

	"github.com/Dizzrt/etstream/kafka"
	"go.uber.org/zap/zapcore"
)

const logTimeFormat = "2006-01-02 15:04:05.000 Z07:00"

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

type Log struct {
	Time      time.Time
	Level     zapcore.Level
	Caller    string
	Message   string
	ExtraData map[string]interface{}
	Traceback string
	RawData   string
}
