//
// File: stash.go
// Created by Dizzrt on 2023/02/27.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package etlog

import (
	"encoding/json"
	"time"

	"github.com/dlclark/regexp2"
	"go.uber.org/zap/zapcore"
)

var (
	logTimeReg      *regexp2.Regexp
	logLevelReg     *regexp2.Regexp
	logCallerReg    *regexp2.Regexp
	logMsgReg       *regexp2.Regexp
	logExtraReg     *regexp2.Regexp
	logTracebackReg *regexp2.Regexp
)

func init() {
	logTimeReg = regexp2.MustCompile(`(?<=^\[)\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}\.\d{3}\s\+\d{2}:\d{2}(?=\])`, 0)
	logLevelReg = regexp2.MustCompile(`(?<=\[)[A-Z]*(?=\])`, 0)
	logCallerReg = regexp2.MustCompile(`(?<=(\[[A-Z]+\]\t)\[).*(?=\])`, 0)
	logMsgReg = regexp2.MustCompile(`(?<=\[.*\]\t\[.*\]\t\[.*\]\t).+(?=\t{)`, 0)
	logExtraReg = regexp2.MustCompile(`{.*}`, 0)
	logTracebackReg = regexp2.MustCompile(`(?<=}\n)[\s\S]+`, 0)
}

func Stash(l string) (Log, error) {
	var err error
	ll := Log{
		RawData: l,
	}

	for i := 0; i < 1; i += 1 {
		// time
		match, err := logTimeReg.FindStringMatch(l)
		if err != nil {
			break
		}

		t, err := time.Parse(logTimeFormat, match.String())
		if err != nil {
			break
		}
		ll.Time = t

		// level
		match, err = logLevelReg.FindStringMatch(l)
		if err != nil {
			break
		}

		ll.Level, err = zapcore.ParseLevel(match.String())
		if err != nil {
			break
		}

		// caller
		match, err = logCallerReg.FindStringMatch(l)
		if err != nil {
			break
		}
		ll.Caller = match.String()

		// msg
		match, err = logMsgReg.FindStringMatch(l)
		if err != nil {
			break
		}
		ll.Message = match.String()

		// extra
		match, err = logExtraReg.FindStringMatch(l)
		if err != nil {
			break
		}

		err = json.Unmarshal([]byte(match.String()), &ll.ExtraData)
		if err != nil {
			break
		}

		// traceback
		match, err = logTracebackReg.FindStringMatch(l)
		if err != nil {
			break
		}
		ll.Traceback = match.String()
	}

	return ll, err
}
