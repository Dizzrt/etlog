//
// File: color.go
// Created by Dizzrt on 2023/02/23.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package color

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

type Color uint8

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

var (
	_levelToColor = map[zapcore.Level]Color{
		zapcore.DebugLevel:  Magenta,
		zapcore.InfoLevel:   Blue,
		zapcore.WarnLevel:   Yellow,
		zapcore.ErrorLevel:  Red,
		zapcore.DPanicLevel: Red,
		zapcore.PanicLevel:  Red,
		zapcore.FatalLevel:  Red,
	}
	_unknownLevelColor = Red

	_levelToCapitalColorString = make(map[zapcore.Level]string, len(_levelToColor))
)

func init() {
	for level, color := range _levelToColor {
		_levelToCapitalColorString[level] = color.Add(level.String())
	}
}

func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

func LevelColorEncoder(l zapcore.Level) string {
	s, ok := _levelToCapitalColorString[l]
	if !ok {
		s = _unknownLevelColor.Add(l.String())
	}
	return s
}
